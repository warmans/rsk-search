package command

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/discord/common"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func NewRateCommand(
	logger *zap.Logger,
	transcriptApiClient api.TranscriptServiceClient) *RateCommand {
	return &RateCommand{
		logger:              logger,
		transcriptApiClient: transcriptApiClient,
	}
}

type RateCommand struct {
	logger              *zap.Logger
	transcriptApiClient api.TranscriptServiceClient
}

func (r *RateCommand) Kind() discordgo.ApplicationCommandOptionType {
	return discordgo.ApplicationCommandOptionString
}

func (r *RateCommand) Name() string {
	return "scrimp-rate"
}

func (r *RateCommand) Description() string {
	return "Create a rating post in the current thread/channel"
}

func (r *RateCommand) Options() []*discordgo.ApplicationCommandOption {
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

func (r *RateCommand) ButtonHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"create-rating-msg": r.handleCreateRatingMsg,
		"submit-rating":     r.handleRating,
		"open-rate-modal":   r.handleOpenRateModal,
	}
}

func (r *RateCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"submit-custom-rating": r.handleCustomRating,
	}
}

func (r *RateCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		r.Name(): r.confirmStartRating,
	}
}

func (r *RateCommand) AutoCompleteHandler() discord.InteractionHandler {
	return r.handleAutocomplete
}

func (r *RateCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{}
}

func (r *RateCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
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

func (r *RateCommand) confirmStartRating(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	selection := i.ApplicationCommandData().Options[0].StringValue()
	if selection == "" {
		return nil
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Confirm",
			Content: fmt.Sprintf("Are you sure you want to request ratings for %s in the current channel?", selection),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Confirm",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("%s:create-rating-msg:%s", r.Name(), selection),
						},
					},
				},
			},
		},
	})
}

func (r *RateCommand) handleCreateRatingMsg(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments (expected 1, got %d)", len(args))
	}
	epid := args[0]

	transcript, err := r.transcriptApiClient.GetTranscript(context.Background(), &api.GetTranscriptRequest{Epid: epid})
	if err != nil {
		return fmt.Errorf("failed to get episode: %w", err)
	}

	// get a list of people that voted for a the previous episode, but not this one and mention them.
	mentions := r.getMissingRatingMentions(transcript)

	_, err = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Content: ratingMessageContent(epid, transcript, mentions),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "1️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:submit-rating:%s:1", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "2️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:submit-rating:%s:2", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "3️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:submit-rating:%s:3", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "4️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:submit-rating:%s:4", r.Name(), epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "5️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:submit-rating:%s:5", r.Name(), epid),
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Custom Rating",
						Emoji: &discordgo.ComponentEmoji{
							Name: "⭐",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("%s:open-rate-modal:%s", r.Name(), epid),
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	common.RespondConfirm(r.logger, s, i, "Created!")

	return nil
}

func (r *RateCommand) handleOpenRateModal(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	fields := []discordgo.MessageComponent{discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.TextInput{
				CustomID:    "value",
				Label:       fmt.Sprintf("%s Rating", args[0]),
				Placeholder: "0-5",
				Style:       discordgo.TextInputShort,
				Required:    true,
				MaxLength:   3,
			},
		},
	}}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("%s:submit-custom-rating:%s", r.Name(), args[0]),
			Title:      "Custom Rating",
			Components: fields,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *RateCommand) handleRating(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

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

	if err := r.confirmSubmission(s, i, idAndRatingParts[0], rating); err != nil {
		return err
	}

	common.RespondConfirm(r.logger, s, i, "Thanks")

	return nil
}

func (r *RateCommand) handleCustomRating(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	strVal := i.Interaction.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	floatVal, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		return fmt.Errorf("invalid value %s: %w", strVal, err)
	}

	if floatVal < 0 || floatVal > 5 {
		return fmt.Errorf("rating must be between 0-5")
	}

	if _, err := r.transcriptApiClient.BulkSetTranscriptRatingScore(context.Background(), &api.BulkSetTranscriptRatingScoreRequest{
		Epid:        args[0],
		OauthSource: "discord",
		Scores: map[string]float32{
			i.Member.User.Username: float32(floatVal),
		},
	}); err != nil {
		return err
	}

	if err := r.confirmSubmission(s, i, args[0], floatVal); err != nil {
		return err
	}

	common.RespondConfirm(r.logger, s, i, "Thanks")

	return nil
}

func (r *RateCommand) confirmSubmission(s *discordgo.Session, i *discordgo.InteractionCreate, episode string, rating float64) error {

	existingRating := ""
	transcript, err := r.transcriptApiClient.GetTranscript(context.Background(), &api.GetTranscriptRequest{Epid: episode})
	if err != nil {
		r.logger.Error("failed to get transcript", zap.Error(err))
	} else {
		existingRating = fmt.Sprintf("(currently %0.2f from %d ratings)", transcript.Ratings.ScoreAvg, transcript.Ratings.NumScores)
	}

	mentions := r.getMissingRatingMentions(transcript)

	if _, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         i.Message.ID,
		Content:    util.ToPtr(ratingMessageContent(episode, transcript, mentions)),
		Components: util.ToPtr(i.Message.Components),
		Channel:    i.Message.ChannelID,
		Embed:      nil,
	}); err != nil {
		return err
	}

	if _, err := s.ChannelMessageSend(i.Message.ChannelID, fmt.Sprintf("%s rated %s %0.2f/5.00 %s", i.Member.DisplayName(), episode, rating, existingRating)); err != nil {
		return err
	}
	return nil
}

func (r *RateCommand) getMissingRatingMentions(current *api.Transcript) []string {
	previousEpid, ok := meta.PreviousEpisode(current.ShortId)
	if ok {
		prevTranscript, err := r.transcriptApiClient.GetTranscript(
			context.Background(),
			&api.GetTranscriptRequest{Epid: previousEpid},
		)
		if err == nil {
			missing := []string{}
			for author := range prevTranscript.Ratings.Scores {
				if _, ok := current.Ratings.Scores[author]; !ok {
					if strings.HasPrefix(author, "discord:") {
						missing = append(missing, fmt.Sprintf("@%s", strings.TrimPrefix(author, "discord:")))
					}
				}
			}
			return missing
		} else {
			r.logger.Error("Failed to get previous episode", zap.Error(err))
		}
	}
	return []string{}
}

func ratingMessageContent(epid string, transcript *api.Transcript, mentions []string) string {
	mentionText := ""
	if len(mentions) > 0 {
		mentionText = fmt.Sprintf("\n --- \n 👀 %s", strings.Join(mentions, " "))
	}
	return fmt.Sprintf(
		"## Rate %s\n-# %s | %s | currently %0.2f from %d ratings %s \n --- \n",
		epid,
		transcript.Name,
		transcript.ReleaseDate,
		transcript.Ratings.ScoreAvg,
		transcript.Ratings.NumScores,
		mentionText,
	)
}
