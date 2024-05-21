package http

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/warmans/rsk-search/pkg/data"
	extract_audio "github.com/warmans/rsk-search/pkg/extract-audio"
	"github.com/warmans/rsk-search/pkg/mediacache"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/quota"
	"github.com/warmans/rsk-search/pkg/store/rw"
	"github.com/warmans/rsk-search/pkg/util"
	"github.com/warmans/rsk-search/service/config"
	"github.com/warmans/rsk-search/service/metrics"
	"go.uber.org/zap"
	"io"
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
	mediaCache *mediacache.Cache,
) *DownloadService {
	return &DownloadService{
		logger:        logger.With(zap.String("component", "downloads-http-server")),
		serviceConfig: serviceConfig,
		rwStoreConn:   rwStoreConn,
		httpMetrics:   httpMetrics,
		episodeCache:  episodeCache,
		mediaCache:    mediaCache,
	}
}

type DownloadService struct {
	logger        *zap.Logger
	serviceConfig config.SearchServiceConfig
	rwStoreConn   *rw.Conn
	httpMetrics   *metrics.HTTPMetrics
	episodeCache  *data.EpisodeCache
	mediaCache    *mediacache.Cache
}

func (c *DownloadService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/dl/archive/episodes-json.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadJSONArchive)))
	router.Path("/dl/archive/episodes-plaintext.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadPlaintextArchive)))
	router.Path("/dl/episode/{episode}.json").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodeJSON)))
	router.Path("/dl/episode/{episode}.txt").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodePlaintext)))

	router.Path("/dl/media/{media_type}/{id}.mp3").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadMP3)))
	router.Path("/dl/media/file/{name}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadFile)))

	//video stuff
	router.Path("/dl/media/gif/{id}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadGif)))
	router.Path("/dl/media/sprite/{episode_id}.jpg").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadVideoSprite)))
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
	if req.URL.Query().Get("pos") != "" || req.URL.Query().Get("ts") != "" {

		if mediaType != MediaTypeEpisode {
			http.Error(resp, "Section exports only supported for episodes", http.StatusNotImplemented)
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

		var startTimestamp, endTimestamp time.Duration
		if pos := req.URL.Query().Get("pos"); pos != "" {
			startTimestamp, endTimestamp, _, err = ep.GetTimestampRange(pos)
			if err != nil {
				http.Error(resp, "invalid position specification", http.StatusBadRequest)
				return
			}
		} else {
			if ts := req.URL.Query().Get("ts"); ts != "" {
				startTimestamp, endTimestamp, err = parseTsParam(ts)
				if err != nil {
					http.Error(resp, fmt.Sprintf("invalid timestamp specification: %s", ts), http.StatusBadRequest)
				}
			} else {
				http.Error(resp, "either position range (pos) or timestamp range (ts) must be specified", http.StatusBadRequest)
				return
			}
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
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s-%s.mp3", episode.ID(), startTimestamp.String(), endTimestamp.String()))
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

func (c *DownloadService) DownloadGif(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileID, ok := vars["id"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}
	ep, err := c.episodeCache.GetEpisode(fileID)
	if errors.Is(err, data.ErrNotFound) || ep == nil {
		http.Error(resp, fmt.Sprintf("unknown episode ID: %s", fileID), http.StatusNotFound)
		return
	}
	if ep.MediaType != models.MediaTypeVideo {
		http.Error(resp, fmt.Sprintf("episode is not a video: %s", fileID), http.StatusBadRequest)
		return
	}
	if ep.MediaFileName == "" {
		http.Error(resp, fmt.Sprintf("no media found for given file: %s", fileID), http.StatusNotFound)
		return
	}
	if err := c.checkQuotas(req.Context()); err != nil {
		if errors.Is(err, DownloadsOverQuota) {
			http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
		} else {
			http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
		}
		return
	}

	pos := req.URL.Query().Get("pos")
	if pos == "" {
		http.Error(resp, "position not given", http.StatusBadRequest)
		return
	}

	// disable caching for custom text
	noCache := false

	startTimestamp, endTimestamp, dialog, err := ep.GetTimestampRange(pos)
	if err != nil {
		c.logger.Error("invalid position", zap.Error(err))
		http.Error(resp, "invalid position specification", http.StatusBadRequest)
		return
	}
	customText := strings.TrimSpace(req.URL.Query().Get("custom_text"))
	if customText != "" {
		if len(customText) > 200 {
			http.Error(resp, "custom_text cannot be more than 200 characters", http.StatusBadRequest)
			return
		}
		c.logger.Info("custom text used", zap.String("custom_text", customText))
		dialog = []string{customText}
		noCache = true
	}
	clipDuration := endTimestamp - startTimestamp
	if clipDuration > time.Second*15 {
		http.Error(resp, fmt.Sprintf("gifs cannot be more than 15 seconds. Given range was %s", clipDuration), http.StatusBadRequest)
		return
	}
	cacheKey := fmt.Sprintf("%s-%s.gif", fileID, pos)
	startTime := time.Now()

	cacheHit, err := c.mediaCache.Get(cacheKey, resp, noCache, func(writer io.Writer) error {
		text := []string{}
		for k, v := range strings.Split(strings.Replace(strings.Join(dialog, " "), "'", "", -1), " ") {
			if k%12 == 0 {
				text = append(text, "\n", v)
				continue
			}
			text = append(text, " ", v)
		}
		countingWriter := NewCountingWriter(writer)
		err = ffmpeg_go.
			Input(path.Join(c.serviceConfig.MediaBasePath, "video", ep.MediaFileName),
				ffmpeg_go.KwArgs{
					"ss": fmt.Sprintf("%0.2f", startTimestamp.Seconds()),
					"to": fmt.Sprintf("%0.2f", endTimestamp.Seconds()),
				}).
			Output("pipe:",
				ffmpeg_go.KwArgs{
					"format": "gif",
					"filter_complex": fmt.Sprintf(
						"[0:v]drawtext=text='%s':fontcolor=white:fontsize=16:box=1:boxcolor=black@0.5:boxborderw=5:x=(w-text_w)/2:y=(h-(text_h+10))",
						util.FfmpegSanitizeDrawtext(strings.TrimSpace(strings.Join(text, ""))),
					),
				},
			).WithOutput(countingWriter, os.Stderr).Run()
		if err != nil {
			c.logger.Error("ffmpeg failed", zap.Error(err))
			return err
		}
		if err := c.incrementQuotas(req.Context(), "gif", fileID, int64(countingWriter.BytesWritten())); err != nil {
			c.logger.Error("Failed to increment quotas", zap.Error(err))
		}
		return nil
	})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		c.logger.Error("cache fetch failed", zap.Error(err))
		return
	}
	if cacheHit {
		c.logger.Info("Cache hit", zap.Duration("time taken", time.Since(startTime)), zap.String("cache_key", cacheKey))
	} else {
		c.logger.Info("Cache miss", zap.Duration("time taken", time.Since(startTime)), zap.String("cache_key", cacheKey))
	}
}

// DownloadVideoSprite returns a sprite containing a single frame for each position stitched together into a
// single image.
func (c *DownloadService) DownloadVideoSprite(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	episodeID, ok := vars["episode_id"]
	if !ok {
		http.Error(resp, "No dialog identifier given", http.StatusBadRequest)
		return
	}
	ep, err := c.episodeCache.GetEpisode(episodeID)
	if errors.Is(err, data.ErrNotFound) || ep == nil {
		http.Error(resp, fmt.Sprintf("unknown episode ID: %s", episodeID), http.StatusNotFound)
		return
	}
	if ep.MediaType != models.MediaTypeVideo {
		http.Error(resp, fmt.Sprintf("episode is not a video: %s", episodeID), http.StatusBadRequest)
		return
	}
	filePath := path.Join(c.serviceConfig.MediaBasePath, "image", "sprite", episodeID+".jpg")
	fileStat, err := os.Stat(filePath)
	if err != nil {
		c.logger.Error("Failed to find media file", zap.String("path", filePath))
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}

	if err = c.incrementQuotas(req.Context(), "sprite", path.Base(filePath), fileStat.Size()); err != nil {
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

func (c *DownloadService) checkQuotas(ctx context.Context) error {
	return c.rwStoreConn.WithStore(func(s *rw.Store) error {
		_, currentMib, err := s.GetMediaStatsForCurrentMonth(ctx)
		if err != nil {
			if err == context.Canceled || strings.HasSuffix(err.Error(), "driver: bad connection") {
				return nil
			}
			return errors.Wrap(err, "failed to get current usage")
		}
		if currentMib > quota.BandwidthQuotaInMiB {
			return DownloadsOverQuota
		}
		return nil
	})
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

func parseTsParam(ts string) (time.Duration, time.Duration, error) {
	parts := strings.Split(ts, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid timestamp range: %s", ts)
	}
	startTs, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start timestamp %s: %w", ts, err)
	}
	endTs, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end timestamp %s: %w", ts, err)
	}
	return time.Duration(startTs) * time.Millisecond, time.Duration(endTs) * time.Millisecond, nil
}

// CountingWriter via https://github.com/jeanfric/goembed/blob/master/countingwriter/countingwriter.go
type CountingWriter struct {
	writer       io.Writer
	bytesWritten int
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{
		writer:       w,
		bytesWritten: 0,
	}
}

func (w *CountingWriter) Write(b []byte) (int, error) {
	n, err := w.writer.Write(b)
	w.bytesWritten += n
	return n, err
}

// BytesWritten returns the number of bytes that were written to the wrapped writer.
func (w *CountingWriter) BytesWritten() int {
	return w.bytesWritten
}
