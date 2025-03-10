package config

import (
	"github.com/spf13/pflag"
	"github.com/warmans/rsk-search/pkg/flag"
)

type SearchServiceConfig struct {
	Env                   string
	BlugeIndexPath        string
	Scheme                string
	Hostname              string
	RewardsDisabled       bool
	FilesBasePath         string
	AudioUriPattern       string
	MediaBasePath         string
	VideoPartialsBasePath string
	ArchiveBasePath       string
}

func (c *SearchServiceConfig) RegisterFlags(fs *pflag.FlagSet, prefix string) {
	flag.StringVarEnv(fs, &c.Env, prefix, "env", "unknown", "set an environment label (used mostly in logging)")
	flag.StringVarEnv(fs, &c.Scheme, prefix, "scheme", "http://", "scheme to use for absolute links")
	flag.StringVarEnv(fs, &c.Hostname, prefix, "hostname", "localhost", "hostname to use for absolute links")
	flag.StringVarEnv(fs, &c.BlugeIndexPath, prefix, "bluge-index-path", "./var/gen/rsk.bluge", "location of bluge search index")
	flag.StringVarEnv(fs, &c.FilesBasePath, prefix, "files-base-path", "./var", "location of static data files")
	flag.StringVarEnv(fs, &c.MediaBasePath, prefix, "media-base-path", "/audio", "location of media files")
	flag.StringVarEnv(fs, &c.VideoPartialsBasePath, prefix, "video-partials-base-path", "./var/video-partials", "partial video files used to generate gifs")
	flag.StringVarEnv(fs, &c.ArchiveBasePath, prefix, "archive-base-path", "./var/archive", "archived files dir")
	flag.BoolVarEnv(fs, &c.RewardsDisabled, prefix, "rewards-disabled", false, "Disable claiming rewards (but sill calculate them)")
	flag.StringVarEnv(fs, &c.AudioUriPattern, prefix, "audio-uri-pattern", "/dl/media/episode/%s.mp3", "episode ID e.g. xfm-S1E01 will be interpolated into this string")
}
