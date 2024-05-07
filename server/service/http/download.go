package http

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/pkg/data"
	extract_audio "github.com/warmans/rsk-search/pkg/extract-audio"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/quota"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/service/config"
	"github.com/warmans/rsk-search/service/metrics"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var DownloadsOverQuota = errors.New("download quota exceeded")

const (
	MediaTypeEpisode = "episode"
	MediaTypeChunk   = "chunk"
	MediaTypeClip    = "clip"
)

func NewDownloadService(
	logger *zap.Logger,
	serviceConfig config.SearchServiceConfig,
	rwStoreConn *rw.Conn,
	httpMetrics *metrics.HTTPMetrics,
	episodeCache *data.EpisodeCache,
) *DownloadService {
	return &DownloadService{
		logger:        logger.With(zap.String("component", "downloads-http-server")),
		serviceConfig: serviceConfig,
		rwStoreConn:   rwStoreConn,
		httpMetrics:   httpMetrics,
		episodeCache:  episodeCache,
	}
}

type DownloadService struct {
	logger        *zap.Logger
	serviceConfig config.SearchServiceConfig
	rwStoreConn   *rw.Conn
	httpMetrics   *metrics.HTTPMetrics
	episodeCache  *data.EpisodeCache
}

func (c *DownloadService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/dl/archive/episodes-json.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadJSONArchive)))
	router.Path("/dl/archive/episodes-plaintext.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadPlaintextArchive)))
	router.Path("/dl/episode/{episode}.json").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodeJSON)))
	router.Path("/dl/episode/{episode}.txt").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodePlaintext)))

	router.Path("/dl/media/{media_type}/{id}.mp3").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadMP3)))
	router.Path("/dl/media/file/{name}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadFile)))
}

func (c *DownloadService) DownloadJSONArchive(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Disposition", "attachment; filename=episodes-json.zip")
	resp.Header().Set("Content-Type", "application/zip")
	http.ServeFile(resp, req, path.Join(c.serviceConfig.FilesBasePath, "gen", "episodes-json.zip"))
}

func (c *DownloadService) DownloadPlaintextArchive(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Disposition", "attachment; filename=episodes-plaintext.zip")
	resp.Header().Set("Content-Type", "application/zip")
	http.ServeFile(resp, req, path.Join(c.serviceConfig.FilesBasePath, "gen", "episodes-plaintext.zip"))
}

func (c *DownloadService) DownloadMP3(resp http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	fileID, ok := vars["id"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}
	mediaType, ok := vars["media_type"]
	if !ok || (mediaType != MediaTypeEpisode && mediaType != MediaTypeChunk && mediaType != MediaTypeClip) {
		http.Error(resp, "No/unknown media type given", http.StatusBadRequest)
		return
	}

	if mediaType == MediaTypeEpisode {
		if !meta.IsValidEpisodeID(fileID) {
			http.Error(resp, "Episode not found", http.StatusNotFound)
			return
		}
	}

	filePath := path.Join(c.serviceConfig.MediaBasePath, mediaType, fmt.Sprintf("%s.mp3", fileID))

	// partial file

	if pos := req.URL.Query().Get("pos"); pos != "" {
		if mediaType != MediaTypeEpisode {
			http.Error(resp, "Offsets only supported for episodes", http.StatusNotImplemented)
			return
		}

		ep, err := c.episodeCache.GetEpisode(fileID)
		if errors.Is(err, data.ErrNotFound) || ep == nil {
			http.Error(resp, fmt.Sprintf("unknown episode ID: %s", fileID), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(resp, "failed to fetch episode metadata", http.StatusInternalServerError)
			return
		}
		startTimestamp, endTimestamp, err := ep.GetTimestampRange(pos)
		if err != nil {
			http.Error(resp, "invalid position specification", http.StatusInternalServerError)
			return
		}
		if err := c.incrementDownloadQuotas(req.Context(), mediaType, fileID, calculateDownloadQuotaUsage(ep, endTimestamp-startTimestamp)); err != nil {
			if errors.Is(err, DownloadsOverQuota) {
				http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
			} else {
				http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
			}
			return
		}
		if err = c.servePartialAudioFile(resp, ep, filePath, startTimestamp, endTimestamp); err != nil {
			c.logger.Error("Failed to serve partial audio file", zap.Error(err))
		}
		return
	}

	// whole file
	fileStat, err := os.Stat(filePath)
	if err != nil {
		c.logger.Error("Failed to find media file", zap.String("path", filePath))
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}
	if err := c.incrementDownloadQuotas(req.Context(), mediaType, fileID, fileStat.Size()); err != nil {
		if errors.Is(err, DownloadsOverQuota) {
			http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
		} else {
			http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
		}
		return
	}
	http.ServeFile(resp, req, filePath)
}

func (c *DownloadService) incrementDownloadQuotas(ctx context.Context, mediaType string, fileID string, fileBytes int64) error {

	c.logger.Debug("Incrementing download quotas", zap.String("media_type", mediaType), zap.String("file_id", fileID), zap.Int64("bytes", fileBytes))

	if err := c.incrementQuotas(ctx, mediaType, fileID, fileBytes); err != nil {
		if errors.Is(err, context.Canceled) {
			// user went away
			return nil
		}
		c.logger.Error(
			"Download failed processing quota",
			zap.Error(err),
			zap.String("episode", fileID),
			zap.Int64("num_bytes", fileBytes),
		)
		return err
	}
	c.httpMetrics.OutboundMediaBytesTotal.Set(float64(fileBytes))
	return nil
}

func (c *DownloadService) servePartialAudioFile(
	resp http.ResponseWriter,
	episode *models.Transcript,
	mp3Path string,
	startTimestamp time.Duration,
	endTimestamp time.Duration,
) error {
	resp.Header().Set("Content-Type", "audio/mpeg")
	//resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%d-%d.mp3", ep.ID(), startTimestamp, endTimestamp))
	return extract_audio.ExtractAudio(resp, mp3Path, startTimestamp, endTimestamp)
}

func (c *DownloadService) DownloadFile(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileName, ok := vars["name"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}

	fileName = path.Base(fileName)
	if fileName == "/" {
		http.Error(resp, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := path.Join(c.serviceConfig.MediaBasePath, "file", fileName)

	fileStat, err := os.Stat(filePath)
	if err != nil {
		c.logger.Error("Failed to find media file", zap.String("path", filePath))
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}

	if err = c.incrementQuotas(req.Context(), "file", fileName, fileStat.Size()); err != nil {
		c.logger.Error("Download failed processing quota", zap.Error(err))
		if errors.Is(err, DownloadsOverQuota) {
			http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
		} else {
			http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
		}
		return
	}

	http.ServeFile(resp, req, filePath)
}

func (c *DownloadService) DownloadEpisodeJSON(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	episode, ok := vars["episode"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}
	if !meta.IsValidEpisodeID(episode) {
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}
	fileName := fmt.Sprintf("%s.json", episode)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	resp.Header().Set("Content-Type", "application/json")
	http.ServeFile(resp, req, path.Join(c.serviceConfig.FilesBasePath, "data", "episodes", fileName))
}

func (c *DownloadService) DownloadEpisodePlaintext(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	episode, ok := vars["episode"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}
	if !meta.IsValidEpisodeID(episode) {
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}
	fileName := fmt.Sprintf("%s.txt", episode)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	resp.Header().Set("Content-Type", "text/plain")
	http.ServeFile(resp, req, path.Join(c.serviceConfig.FilesBasePath, "gen", "plaintext", fileName))
}

func (c *DownloadService) incrementQuotas(ctx context.Context, mediaType string, fileID string, fileBytes int64) error {
	return c.rwStoreConn.WithStore(func(s *rw.Store) error {

		fileMib := quota.BytesAsMib(fileBytes)

		_, currentMib, err := s.GetMediaStatsForCurrentMonth(ctx)
		if err != nil {
			if err == context.Canceled || strings.HasSuffix(err.Error(), "driver: bad connection") {
				return nil
			}
			return errors.Wrap(err, "failed to get current usage")
		}
		c.httpMetrics.OutboundMediaQuotaRemaining.Set(float64(quota.BandwidthQuotaInMiB - currentMib))
		if currentMib+fileMib > quota.BandwidthQuotaInMiB {
			return DownloadsOverQuota
		}
		if err := s.IncrementMediaAccessLog(ctx, mediaType, fileID, fileMib); err != nil {
			if err == context.Canceled {
				return nil
			}
			return errors.Wrap(err, "failed to increment access log bytes")
		}
		return nil
	})
}

func calculateDownloadQuotaUsage(ep *models.Transcript, duration time.Duration) int64 {
	var bitrate = 320.00 // default to a high bitrate
	if bitrateStr, ok := ep.Meta[models.MetadataTypeBitrateKbps]; ok {
		if bitrateFloat, err := strconv.ParseFloat(bitrateStr, 64); err == nil {
			bitrate = bitrateFloat
		}
	}
	// bit rate needs to be converted into bytes
	return int64((bitrate * duration.Seconds()) / 8)
}
