package http

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	ffmpeg_go "github.com/warmans/ffmpeg-go/v2"
	"github.com/warmans/rsk-search/pkg/archive"
	"github.com/warmans/rsk-search/pkg/data"
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
	"math/rand/v2"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type partialFileOptions struct {
	stripID3   bool
	wantFormat string
}

type partialFileOption func(opts *partialFileOptions)

func withoutID3Metadata(enabled bool) partialFileOption {
	return func(opts *partialFileOptions) {
		opts.stripID3 = enabled
	}
}

func withOutputFormat(format string) partialFileOption {
	return func(opts *partialFileOptions) {
		opts.wantFormat = format
	}
}

func resolvePartialFileOptions(opts []partialFileOption) *partialFileOptions {
	options := &partialFileOptions{}
	for _, v := range opts {
		v(options)
	}
	return options
}

var DownloadsOverQuota = errors.New("download quota exceeded")

func NewDownloadService(
	logger *zap.Logger,
	serviceConfig config.SearchServiceConfig,
	rwStoreConn *rw.Conn,
	httpMetrics *metrics.HTTPMetrics,
	episodeCache *data.EpisodeCache,
	mediaCache *mediacache.Cache,
	archiveStore *archive.Store,
) (*DownloadService, error) {

	partials := map[string][]string{}
	files, err := os.ReadDir(serviceConfig.VideoPartialsBasePath)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".mp4") {
			continue
		}
		actorName := strings.Split(f.Name(), "-")[0]
		if _, ok := partials[actorName]; !ok {
			partials[actorName] = []string{f.Name()}
		} else {
			partials[actorName] = append(partials[actorName], f.Name())
		}
	}

	return &DownloadService{
		logger:        logger.With(zap.String("component", "downloads-http-server")),
		serviceConfig: serviceConfig,
		rwStoreConn:   rwStoreConn,
		httpMetrics:   httpMetrics,
		episodeCache:  episodeCache,
		mediaCache:    mediaCache,
		videoPartials: partials,
		archiveStore:  archiveStore,
	}, nil
}

type DownloadService struct {
	logger        *zap.Logger
	serviceConfig config.SearchServiceConfig
	rwStoreConn   *rw.Conn
	httpMetrics   *metrics.HTTPMetrics
	episodeCache  *data.EpisodeCache
	mediaCache    *mediacache.Cache
	videoPartials map[string][]string
	archiveStore  *archive.Store
}

func (c *DownloadService) RegisterHTTP(ctx context.Context, router *mux.Router) {
	router.Path("/dl/archive/episodes-json.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadJSONArchive)))
	router.Path("/dl/archive/episodes-plaintext.zip").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadPlaintextArchive)))
	router.Path("/dl/archive/file/{name}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadArchiveFile)))

	router.Path("/dl/episode/{episode}.json").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodeJSON)))
	router.Path("/dl/episode/{episode}.txt").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodePlaintext)))

	router.Path("/dl/media/sprite/{episode_id}.jpg").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadVideoSprite)))
	router.Path("/dl/media/file/{name}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadFile)))
	router.Path("/dl/media/{episode_id}.{format}").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadEpisodeMedia)))
	router.Path("/dl/sprite/{episode_id}.jpg").Handler(handlers.RecoveryHandler()(http.HandlerFunc(c.DownloadVideoSprite)))

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
	req *http.Request,
	resp http.ResponseWriter,
	episode *models.Transcript,
	mp3Path string,
	startTimestamp time.Duration,
	endTimestamp time.Duration,
	opts ...partialFileOption,
) error {

	options := resolvePartialFileOptions(opts)

	var format, mimeType string
	var writeData func(ss time.Duration, to time.Duration, w io.Writer) error

	switch options.wantFormat {
	case "gif":
		rawDialog := episode.GetDialogAtTimestampRange(startTimestamp, endTimestamp)

		// todo concat videos to for multiple lines
		if len(rawDialog) != 1 {
			return fmt.Errorf("can only create a gif of a single line")
		}
		dialog := []string{}
		for _, v := range rawDialog {
			dialog = append(dialog, v.Content)
		}

		format = "gif"
		mimeType = "image/gif"
		writeData = func(ss time.Duration, to time.Duration, w io.Writer) error {

			videoName := c.getVideoPartialName(rawDialog[0].Actor)
			args := ffmpeg_go.KwArgs{}
			if strings.HasSuffix(videoName, ".mp4") {
				args["ss"] = fmt.Sprintf("%d", rand.IntN(5))
			}

			return ffmpeg_go.
				Output(
					[]*ffmpeg_go.Stream{
						ffmpeg_go.Input(
							fmt.Sprintf(
								"%s/%s",
								c.serviceConfig.VideoPartialsBasePath,
								videoName,
							),
							args,
						),
					},
					"pipe:",
					ffmpeg_go.KwArgs{
						"format": "gif",
						"t":      fmt.Sprintf("%0.2f", (endTimestamp - startTimestamp).Seconds()),
						"filter_complex": fmt.Sprintf(
							"scale=596:336,drawtext=text='%s':fontcolor=white:fontsize=16:box=1:boxcolor=black@0.5:boxborderw=5:x=(w-text_w)/2:y=(h-(text_h+10))",
							util.FfmpegSanitizeDrawtext(FormatGifText(56, reSplitDialog(dialog))),
						),
					},
				).
				WithOutput(w, os.Stderr).
				Run()
		}
	case "webm":
		format = "webm"
		mimeType = "video/webm"
		writeData = func(ss time.Duration, to time.Duration, w io.Writer) error {
			input := []*ffmpeg_go.Stream{
				ffmpeg_go.Input(
					mp3Path,
					ffmpeg_go.KwArgs{
						"ss": fmt.Sprintf("%0.2f", startTimestamp.Seconds()),
						"to": fmt.Sprintf("%0.2f", endTimestamp.Seconds()),
					},
				),
				ffmpeg_go.Input("assets/default-lg.jpg"),
			}

			return ffmpeg_go.
				Output(
					input,
					"pipe:",
					ffmpeg_go.KwArgs{
						"_map_0": "0:a",
						"_map_1": "1:v",
						"vf":     "scale=220:220",
						"format": "webm",
					},
				).
				WithOutput(w, os.Stderr).
				Run()
		}
	default:
		format = "mp3"
		mimeType = "audio/mpeg"
		writeData = func(ss time.Duration, to time.Duration, w io.Writer) error {

			outputArgs := ffmpeg_go.KwArgs{
				"format": "mp3",
				"acodec": "copy",
			}
			if options.stripID3 {
				outputArgs["map_metadata"] = "-1"
			}

			return ffmpeg_go.
				Input(mp3Path,
					ffmpeg_go.KwArgs{
						"ss": fmt.Sprintf("%0.2f", ss.Seconds()),
						"to": fmt.Sprintf("%0.2f", to.Seconds()),
					}).
				Output("pipe:", outputArgs).WithOutput(w, os.Stderr).Run()
		}
	}

	if req.Header.Get("Range") == "" {
		// just return the whole file
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", mimeType)
		resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s-%s.%s", episode.ID(), startTimestamp.String(), endTimestamp.String(), format))
		return writeData(startTimestamp, endTimestamp, resp)
	}

	// If they client is asking for a range, there doesn't seem to be a way to efficiently stream the data.
	// Most solutions just buffer in memory to get a ReadSeeker. A temp file seems preferable since memory is
	// more scarce than disk.

	partial, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("partial-*.%s", format))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if err := os.Remove(path.Join(partial.Name())); err != nil {
			c.logger.Error("failed to remove temporary file", zap.Error(err), zap.String("path", partial.Name()))
		}
	}()
	defer partial.Close()

	if err := writeData(startTimestamp, endTimestamp, partial); err != nil {
		partial.Close()
		return fmt.Errorf("failed to extract data: %w", err)
	}

	resp.Header().Set("Content-Type", mimeType)
	resp.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s-%s.%s", episode.ID(), startTimestamp.String(), endTimestamp.String(), format))
	http.ServeContent(resp, req, path.Base(partial.Name()), time.Now(), partial)
	return nil
}

func (c *DownloadService) DownloadEpisodeMedia(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	episodeID, ok := vars["episode_id"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}
	wantFormat, ok := vars["format"]
	if !ok || (wantFormat != "mp3" && wantFormat != "gif" && wantFormat != "webm") {
		http.Error(resp, "only mp3, webm and gif are supported", http.StatusBadRequest)
		return
	}

	episode, err := c.episodeCache.GetEpisode(episodeID)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			http.Error(resp, "Episode not found", http.StatusNotFound)
			return
		}
		http.Error(resp, "Failed to fetch episode metadata", http.StatusInternalServerError)
		return
	}
	if (wantFormat == "mp3" || wantFormat == "webm") && episode.Media.AudioFileName == "" {
		http.Error(resp, "this episode doesn't have audio available", http.StatusNotFound)
		return
	}

	// partial file download
	if req.URL.Query().Has("pos") || req.URL.Query().Has("ts") {

		stripID3Tags := req.URL.Query().Get("strip_tags")

		//determine the time range to export
		var startTimestamp, endTimestamp time.Duration
		if pos := req.URL.Query().Get("pos"); pos != "" {
			startTimestamp, endTimestamp, _, err = episode.GetDialogAtPosition(pos)
			if err != nil {
				http.Error(resp, "invalid position specification", http.StatusBadRequest)
				return
			}
		} else {
			if ts := req.URL.Query().Get("ts"); ts != "" {
				startTimestamp, endTimestamp, err = parseTsParam(ts)
				if err != nil {
					http.Error(resp, fmt.Sprintf("invalid timestamp specification: %s", ts), http.StatusBadRequest)
					return
				}
				if endTimestamp == 0 {
					if episode.Media.AudioDurationMs > 0 {
						endTimestamp = time.Duration(episode.Media.AudioDurationMs) * time.Millisecond
					} else {
						http.Error(resp, "end timestamp must be specified", http.StatusBadRequest)
						return
					}
				}
			} else {
				http.Error(resp, "either position range (pos) or timestamp range (ts) must be specified", http.StatusBadRequest)
				return
			}
		}
		var customText *string
		if req.URL.Query().Has("custom_text") {
			customTextParam := strings.TrimSpace(req.URL.Query().Get("custom_text"))
			if customTextParam != "" {
				if len(customTextParam) > 200 {
					http.Error(resp, "custom_text cannot be more than 200 characters", http.StatusBadRequest)
					return
				}
				customText = util.ToPtr(customTextParam)
			} else {
				customText = util.ToPtr("")
			}
		}

		c.logger.Debug(
			"Exporting partial file",
			zap.String("episode_id", episode.ShortID()),
			zap.String("format", wantFormat),
			zap.Duration("start", startTimestamp),
			zap.Duration("end", endTimestamp),
			zap.Stringp("custom_text", customText),
		)

		if episode.MediaType == models.MediaTypeVideo {
			switch wantFormat {
			// serve partial video
			case "gif":
				c.downloadGif(req.Context(), resp, episode, startTimestamp, endTimestamp, customText)
				return
				// serve partial audio
			default:
				http.Error(resp, "Unknown format requested for video media", http.StatusBadRequest)
				return
			}
		}

		// audio file
		if wantFormat != "mp3" && wantFormat != "webm" && wantFormat != "gif" {
			http.Error(resp, "Only audio can be exported without either pos or ts specified", http.StatusBadRequest)
			return
		}

		// todo: this is wrong because the Range header may prevent the full file being returned,
		// note sure if it's worth just using a counting response writer.
		if err := c.incrementDownloadQuotas(
			req.Context(),
			wantFormat,
			episode.ShortID(),
			calculateDownloadQuotaUsage(episode, endTimestamp-startTimestamp),
		); err != nil {
			if errors.Is(err, DownloadsOverQuota) {
				http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
			} else {
				http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
			}
			return
		}
		if err = c.servePartialAudioFile(
			req,
			resp,
			episode,
			path.Join(c.serviceConfig.MediaBasePath, "episode", episode.Media.AudioFileName),
			startTimestamp,
			endTimestamp,
			withoutID3Metadata(stripID3Tags == "true"),
			withOutputFormat(wantFormat),
		); err != nil {
			c.logger.Error("Failed to serve partial audio file", zap.Error(err))
			return
		}
		return

	}

	// whole audio file

	var filePath string
	if req.URL.Query().Get("remastered") == "1" && episode.Media.RemasteredAudioFileName != "" {
		filePath = path.Join(c.serviceConfig.MediaBasePath, "episode", "remaster", episode.Media.RemasteredAudioFileName)
	} else {
		filePath = path.Join(c.serviceConfig.MediaBasePath, "episode", episode.Media.AudioFileName)
	}

	c.logger.Debug(
		"Exporting full file",
		zap.String("format", wantFormat),
		zap.String("episode_id", episode.ShortID()),
		zap.String("file_path", filePath),
	)

	fileStat, err := os.Stat(filePath)
	if err != nil {
		c.logger.Error("Failed to find media file", zap.String("path", filePath))
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}
	if err := c.incrementDownloadQuotas(req.Context(), wantFormat, episode.ShortID(), fileStat.Size()); err != nil {
		if errors.Is(err, DownloadsOverQuota) {
			http.Error(resp, "Bandwidth quota exhausted", http.StatusTooManyRequests)
		} else {
			http.Error(resp, "Failed to calculate bandwidth quota", http.StatusInternalServerError)
		}
		return
	}
	http.ServeFile(resp, req, filePath)
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

func (c *DownloadService) downloadGif(ctx context.Context, resp http.ResponseWriter, episode *models.Transcript, startTimestamp time.Duration, endTimestamp time.Duration, customText *string) {

	c.logger.Debug("Exporting gif", zap.Duration("start", startTimestamp), zap.Duration("end", endTimestamp), zap.Stringp("custom_text", customText))

	clipDuration := endTimestamp - startTimestamp
	if clipDuration > time.Second*15 {
		// clip dialog to the maximum
		endTimestamp = startTimestamp + time.Second*15
	}
	cacheKey := fmt.Sprintf("%s-%s-%s.gif", episode.ShortID(), startTimestamp.String(), endTimestamp.String())
	startTime := time.Now()

	disableCaching := false
	dialog := reSplitDialog(episode.GetDialogContentAtTimestampRange(startTimestamp, endTimestamp))
	if customText != nil {
		disableCaching = true
		if *customText == "" {
			dialog = []string{}
		} else {
			dialog = strings.Split(*customText, "\n")
		}
	}

	cacheHit, err := c.mediaCache.Get(cacheKey, resp, disableCaching, func(writer io.Writer) error {

		//todo: write content type headers?

		countingWriter := NewCountingWriter(writer)
		err := ffmpeg_go.
			Input(path.Join(c.serviceConfig.MediaBasePath, "video", episode.Media.VideoFileName),
				ffmpeg_go.KwArgs{
					"ss": fmt.Sprintf("%0.2f", startTimestamp.Seconds()),
					"to": fmt.Sprintf("%0.2f", endTimestamp.Seconds()),
				}).
			Output("pipe:",
				ffmpeg_go.KwArgs{
					"format": "gif",
					"filter_complex": fmt.Sprintf(
						"[0:v]drawtext=text='%s':fontcolor=white:fontsize=16:box=1:boxcolor=black@0.5:boxborderw=5:x=(w-text_w)/2:y=(h-(text_h+10))",
						util.FfmpegSanitizeDrawtext(FormatGifText(56, dialog)),
					),
				},
			).WithOutput(countingWriter, os.Stderr).Run()
		if err != nil {
			c.logger.Error("ffmpeg failed", zap.Error(err))
			return err
		}
		if err := c.incrementQuotas(ctx, "gif", episode.ShortID(), int64(countingWriter.BytesWritten())); err != nil {
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

func (c *DownloadService) incrementQuotas(ctx context.Context, mediaType string, fileID string, fileBytes int64) error {
	return c.rwStoreConn.WithStore(func(s *rw.Store) error {

		fileMib := quota.BytesAsMib(fileBytes)

		_, currentMib, err := s.GetMediaStatsForCurrentMonth(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || strings.HasSuffix(err.Error(), "driver: bad connection") {
				return nil
			}
			return errors.Wrap(err, "failed to get current usage")
		}
		c.httpMetrics.OutboundMediaQuotaRemaining.Set(float64(quota.BandwidthQuotaInMiB - currentMib))
		if currentMib+fileMib > quota.BandwidthQuotaInMiB {
			return DownloadsOverQuota
		}
		if err := s.IncrementMediaAccessLog(ctx, mediaType, fileID, fileMib); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return errors.Wrap(err, "failed to increment access log bytes")
		}
		return nil
	})
}

func (c *DownloadService) getVideoPartialName(actor string) string {
	if names, ok := c.videoPartials[actor]; ok && len(names) > 0 {
		return names[rand.IntN(len(names))]
	}
	return "xfm.jpg"
}

func (c *DownloadService) DownloadArchiveFile(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileName, ok := vars["name"]
	if !ok {
		http.Error(resp, "No episode identifier given", http.StatusBadRequest)
		return
	}

	fileName = path.Clean(path.Base(fileName))
	if fileName == "/" {
		http.Error(resp, "Invalid filename", http.StatusBadRequest)
		return
	}

	valid, err := c.archiveStore.FileExists(fileName)
	if err != nil {
		http.Error(resp, "Failed to process request", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(resp, "Unknown file", http.StatusNotFound)
		return
	}

	filePath := path.Join(c.serviceConfig.ArchiveBasePath, fileName)

	fileStat, err := os.Stat(filePath)
	if err != nil {
		c.logger.Error("Failed to find media file", zap.String("path", filePath))
		http.Error(resp, "Episode not found", http.StatusNotFound)
		return
	}

	if err = c.incrementQuotas(req.Context(), "archive_file", fileName, fileStat.Size()); err != nil {
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
	var startTs int64
	var endTs int64
	var err error
	if len(parts) > 0 {
		startTs, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid start timestamp %s: %w", ts, err)
		}
	}
	if len(parts) > 1 {
		endTs, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end timestamp %s: %w", ts, err)
		}
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

// FormatGifText
// max length should be 56ish
func FormatGifText(maxLineLength int, lines []string) string {
	text := []string{}
	for _, line := range lines {
		currentLine := []string{}
		for _, word := range strings.Split(line, " ") {
			if lineLength(currentLine)+(len(word)+1) > maxLineLength {
				text = append(text, strings.Join(currentLine, " "))
				currentLine = []string{word}
				continue
			}
			currentLine = append(currentLine, word)
		}
		if len(currentLine) > 0 {
			text = append(text, strings.Join(currentLine, " "))
		}
	}
	return strings.TrimSpace(strings.Replace(strings.Join(text, "\n"), "'", "’", -1))
}

func lineLength(line []string) int {
	if len(line) == 0 {
		return 0
	}
	total := 0
	for _, v := range line {
		total += len(v)
	}
	// total + number of spaces that would be in the line
	return total + (len(line) - 1)
}

// ensure each line in the slice is free from linebreak
func reSplitDialog(dialog []string) []string {
	fixed := []string{}
	for _, line := range dialog {
		fixed = append(fixed, strings.Split(line, "\n")...)
	}
	return fixed
}
