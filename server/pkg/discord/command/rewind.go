package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/discord/common"
	"github.com/warmans/rsk-search/pkg/meta"
	"go.uber.org/zap"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var tagRegex = regexp.MustCompile(`^\\tag\s+(.+)\s+([0-9hms\.]+)$`)

func NewRewindCommand(
	logger *zap.Logger,
	rewindStateDir string,
	transcriptApiClient api.TranscriptServiceClient) (*RewindCommand, error) {

	b := &RewindCommand{
		logger:              logger,
		rewindStateDir:      rewindStateDir,
		transcriptApiClient: transcriptApiClient,
		rewindThreadCache:   &sync.Map{},
	}
	if err := b.initRewindCache(); err != nil {
		return nil, err
	}

	return b, nil
}

type RewindCommand struct {
	logger              *zap.Logger
	rewindStateDir      string
	transcriptApiClient api.TranscriptServiceClient
	rewindStateLock     sync.RWMutex
	rewindThreadCache   *sync.Map
}

func (r *RewindCommand) Kind() discordgo.ApplicationCommandOptionType {
	return discordgo.ApplicationCommandOptionString
}

func (r *RewindCommand) Name() string {
	return "scrimp-rewind"
}

func (r *RewindCommand) Description() string {
	return "Start a rewind thread for the given episode"
}

func (r *RewindCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:         "episode",
			Description:  "The episode ID in the format xfm-S01E01",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
	}
}

func (r *RewindCommand) ButtonHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"rewind-start":   r.rewindStart,
		"episode-rating": r.rateEpisode,
	}
}

func (r *RewindCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{}
}

func (r *RewindCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		r.Name(): r.rewindBegin,
	}
}

func (r *RewindCommand) AutoCompleteHandler() discord.InteractionHandler {
	return r.handleAutocomplete
}

func (r *RewindCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{
		r.handleThreadMessage,
	}
}

func (r *RewindCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	data := i.ApplicationCommandData()

	prefix := strings.TrimSpace(data.Options[0].StringValue())

	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, epID := range meta.EpisodeList() {
		if strings.HasPrefix(strings.ToLower(epID), strings.ToLower(strings.TrimSpace(prefix))) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  epID,
				Value: epID,
			})
		}
		if len(choices) > 24 {
			break
		}
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}); err != nil {
		r.logger.Error("Failed to respond with autocomplete options", zap.Error(err))
	}
	return nil
}

func (r *RewindCommand) rewindBegin(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	selection := i.ApplicationCommandData().Options[0].StringValue()
	if selection == "" {
		return nil
	}

	return r.confirmRewindStart(s, i, selection)
}

func (r *RewindCommand) rewindStart(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
	}
	epid := args[0]

	content, err := r.getEpisodeSummary(epid)
	if err != nil {
		return err
	}
	initialMessage, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "1️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:episode-rating:%s:1", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "2️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:episode-rating:%s:2", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "3️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:episode-rating:%s:3", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "4️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:episode-rating:%s:4", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "5️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:episode-rating:%s:5", r.Name(), epid),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	thread, err := s.MessageThreadStartComplex(initialMessage.ChannelID, initialMessage.ID, &discordgo.ThreadStart{
		Name: fmt.Sprintf("%s REWIND", epid),
		Type: discordgo.ChannelTypeGuildPublicThread,
	})
	if err != nil {
		return err
	}

	if err := r.createRewindState(RewindState{
		OriginalMessageID:      initialMessage.ID,
		OriginalMessageChannel: initialMessage.ChannelID,
		AnswerThreadID:         thread.ID,
		EpisodeID:              epid,
	}); err != nil {
		return err
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Done!",
		},
	}); err != nil {
		return err
	}

	return nil
}

func (r *RewindCommand) rateEpisode(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	idAndRatingParts := strings.Split(args[0], ":")
	rating, err := strconv.ParseFloat(idAndRatingParts[1], 32)
	if err != nil {
		return fmt.Errorf("failed to parse rating: %w", err)
	}

	if _, err := r.transcriptApiClient.BulkSetTranscriptRatingScore(context.Background(), &api.BulkSetTranscriptRatingScoreRequest{
		Epid:        idAndRatingParts[0],
		OauthSource: "discord",
		Scores: map[string]float32{
			i.Member.User.Username: float32(rating),
		},
	}); err != nil {
		return err
	}

	if _, err := s.ChannelMessageSend(i.Message.ID, fmt.Sprintf("%s rated the episode %s/5", i.Member.DisplayName(), idAndRatingParts[1])); err != nil {
		return err
	}

	common.RespondConfirm(r.logger, s, i, "Rating accepted")
	return nil
}

func (r *RewindCommand) confirmRewindStart(s *discordgo.Session, i *discordgo.InteractionCreate, epid string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Confirm Rewind",
			Content: fmt.Sprintf("Are you sure you want to start a rewind thread for %s?", epid),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Confirm",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("%s:rewind-start:%s", r.Name(), epid),
						},
					},
				},
			},
		},
	})
}

func (r *RewindCommand) handleThreadMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, isValidThread := r.rewindThreadCache.Load(m.ChannelID); !isValidThread {
		return
	}
	if !tagRegex.MatchString(m.Content) {
		return
	}

	matches := tagRegex.FindStringSubmatch(m.Content)
	tag := matches[1]
	duration, err := time.ParseDuration(matches[2])
	if err != nil {
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔥"); err != nil {
			fmt.Println("failed to add reaction ", err.Error())
		}
		if err2 := r.sendThreadMessage(s, m.ChannelID, fmt.Sprintf("Failed to parse tag location %s: %s", matches[2], err.Error())); err2 != nil {
			r.logger.Error("Failed to send error message to thread", zap.Error(err2))
		}
		return
	}
	if err := r.openRewindStateForReading(m.ChannelID, func(cw *RewindState) error {
		_, err := r.transcriptApiClient.BulkSetTranscriptTags(context.Background(), &api.BulkSetTranscriptTagsRequest{
			Epid: cw.EpisodeID,
			Tags: []*api.Tag{
				{
					Name:      strings.TrimSpace(tag),
					Timestamp: duration.String(),
				},
			},
		})
		return err
	}); err != nil {
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔥"); err != nil {
			fmt.Println("failed to add reaction ", err.Error())
		}
		if err2 := r.sendThreadMessage(s, m.ChannelID, fmt.Sprintf("Failed to store tag: %s", err.Error())); err2 != nil {
			r.logger.Error("Failed to send error message to thread", zap.Error(err2))
		}
		return
	}
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
		fmt.Println("failed to add reaction ", err.Error())
		return
	}
}

func (r *RewindCommand) getEpisodeSummary(epid string) (string, error) {
	transcript, err := r.transcriptApiClient.GetTranscript(context.Background(), &api.GetTranscriptRequest{
		Epid:    fmt.Sprintf("ep-%s", epid),
		WithRaw: false,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		`**%s REWIND** | %s | 🔉 https://scrimpton.com/ep/%s

This is a rewind thread. Listen to the episode using the link above. 

**Other Commands:**
 * \tag [tag name] [timestamp in duration format e.g. 12m19s]

**Rating:**
`,
		transcript.ShortId,
		transcript.ReleaseDate,
		transcript.Id,
	), nil

}

func (r *RewindCommand) initRewindCache() error {
	entries, err := os.ReadDir(r.rewindStateDir)
	if err != nil {
		return err
	}
	for _, v := range entries {
		if !strings.HasSuffix(v.Name(), ".json") || v.IsDir() {
			continue
		}
		threadID := strings.TrimSuffix(path.Base(v.Name()), ".json")
		r.rewindThreadCache.Store(threadID, struct{}{})
	}
	return nil
}

func (r *RewindCommand) createRewindState(state RewindState) error {
	_, err := os.Stat(path.Join(r.rewindStateDir, fmt.Sprintf("%s.json", state.AnswerThreadID)))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return nil
	}
	f, err := os.Create(path.Join(r.rewindStateDir, fmt.Sprintf("%s.json", state.AnswerThreadID)))
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(state); err != nil {
		return err
	}

	r.rewindThreadCache.Store(state.AnswerThreadID, struct{}{})

	return nil

}

func (r *RewindCommand) sendThreadMessage(s *discordgo.Session, threadID string, message string) error {
	if _, err := s.ChannelMessageSend(
		threadID,
		message,
	); err != nil {
		return err
	}
	return nil
}

func (r *RewindCommand) openRewindStateForReading(channelID string, cb func(cw *RewindState) error) error {
	r.rewindStateLock.RLock()
	defer r.rewindStateLock.RUnlock()

	f, err := os.Open(path.Join(r.rewindStateDir, fmt.Sprintf("%s.json", channelID)))
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	s := RewindState{}
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return err
	}

	return cb(&s)
}

//func (r *RewindCommand) openRewindStateForWriting(channelID string, cb func(cw *RewindState) (*RewindState, error)) error {
//	r.rewindStateLock.Lock()
//	defer r.rewindStateLock.Unlock()
//
//	f, err := os.OpenFile(path.Join(r.rewindStateDir, fmt.Sprintf("%s.json", channelID)), os.O_RDWR|os.O_EXCL, 0666)
//	if err != nil {
//		return err
//	}
//	defer func(f *os.File) {
//		err := f.Close()
//		if err != nil {
//			r.logger.Error("failed to close rewind state", zap.Error(err))
//		}
//	}(f)
//
//	s := &RewindState{}
//	if err := json.NewDecoder(f).Decode(s); err != nil {
//		return err
//	}
//	s, err = cb(s)
//	if err != nil || s == nil {
//		return err
//	}
//
//	if err := f.Truncate(0); err != nil {
//		return err
//	}
//	if _, err := f.Seek(0, 0); err != nil {
//		return err
//	}
//
//	enc := json.NewEncoder(f)
//	enc.SetIndent("", "  ")
//	return enc.Encode(s)
//}

type RewindState struct {
	OriginalMessageID      string
	OriginalMessageChannel string
	AnswerThreadID         string
	EpisodeID              string
}
