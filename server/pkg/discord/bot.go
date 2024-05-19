package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/searchterms"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var punctuation = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
var spaces = regexp.MustCompile(`[\s]{2,}`)

type CustomID struct {
	EpisodeID       string          `json:"e,omitempty"`
	Position        int             `json:"p,omitempty"`
	ContentModifier ContentModifier `json:"t,omitempty"`
}

type ContentModifier uint8

const (
	ContentModifierNone ContentModifier = iota
	ContentModifierTextOnly
	ContentModifierAudioOnly
)

func NewBot(
	logger *zap.Logger,
	session *discordgo.Session,
	guildID string,
	webUrl string,
	transcriptApiClient api.TranscriptServiceClient,
	searchApiClient api.SearchServiceClient,
) *Bot {
	bot := &Bot{
		logger:              logger,
		session:             session,
		guildID:             guildID,
		webUrl:              webUrl,
		transcriptApiClient: transcriptApiClient,
		searchApiClient:     searchApiClient,
		commands: []*discordgo.ApplicationCommand{
			{
				Name:        "scrimp",
				Description: "Search with confirmation",
				Type:        discordgo.ChatApplicationCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "query",
						Description:  "enter a partial quote",
						Type:         discordgo.ApplicationCommandOptionString,
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
	}
	bot.commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"scrimp": bot.scrimptonQueryBegin,
	}
	bot.buttonHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, suffix string){
		"scrimp_confirm": bot.scrimptonQueryComplete,
	}

	return bot
}

type Bot struct {
	logger              *zap.Logger
	session             *discordgo.Session
	guildID             string
	webUrl              string
	transcriptApiClient api.TranscriptServiceClient
	searchApiClient     api.SearchServiceClient
	commands            []*discordgo.ApplicationCommand
	commandHandlers     map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	buttonHandlers      map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, customIdPayload string)
	createdCommands     []*discordgo.ApplicationCommand
}

func (b *Bot) Start() error {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			// exact match
			if h, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionApplicationCommandAutocomplete:
			// exact match
			if h, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			// prefix match buttons to allow additional data in the customID
			for k, h := range b.buttonHandlers {
				actionPrefix := fmt.Sprintf("%s:", k)
				if strings.HasPrefix(i.MessageComponentData().CustomID, actionPrefix) {
					h(s, i, strings.TrimPrefix(i.MessageComponentData().CustomID, actionPrefix))
				} else {
					b.respondError(s, i, fmt.Errorf("invalid customID format: %s", k))
					return
				}
			}
		}
	})
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}
	var err error
	b.createdCommands, err = b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, b.guildID, b.commands)
	if err != nil {
		return fmt.Errorf("cannot register commands: %w", err)
	}
	return nil
}

func (b *Bot) Close() error {
	// cleanup commands
	for _, cmd := range b.createdCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.guildID, cmd.ID)
		if err != nil {
			return fmt.Errorf("cannot delete %s command: %w", cmd.Name, err)
		}
	}
	return b.session.Close()
}

func (b *Bot) scrimptonQueryBegin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		selection := i.ApplicationCommandData().Options[0].StringValue()
		if selection == "" {
			return
		}
		ids := strings.Split(selection, ":")
		if len(ids) < 2 {
			b.logger.Error("unexpected selection format", zap.String("selection", selection))
			b.respondError(s, i, fmt.Errorf("invalid selection"))
			return
		}
		pos, err := strconv.Atoi(ids[1])
		if err != nil {
			b.logger.Error("failed to parse position", zap.String("pos", ids[1]), zap.Error(err))
			b.respondError(s, i, fmt.Errorf("internal error, unknown selection"))
			return
		}
		username := "unknown"
		if i.Member != nil {
			username = i.Member.DisplayName()
		}
		dialog, err := b.transcriptApiClient.GetTranscriptDialog(context.Background(), &api.GetTranscriptDialogRequest{
			Epid:            ids[0],
			Pos:             int32(pos),
			NumContextLines: 2,
		})
		if err != nil {
			b.respondError(s, i, fmt.Errorf("failed to fetch selected line"))
			return
		}

		// videos will respond with a gif instead of audio
		if dialog.TranscriptMeta.MediaType == api.MediaType_VIDEO {
			interactionResponse, err, cleanup := b.videoFileResponse(dialog, ids[0], pos, username)
			if err != nil {
				b.respondError(s, i, err)
				return
			}
			defer cleanup()
			interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral

			noModifier, err := encodeCustomID("scrimp_confirm", ids[0], pos, ContentModifierNone)
			if err != nil {
				b.logger.Error("failed to marshal customID", zap.Error(err))
				b.respondError(s, i, fmt.Errorf("failed to create interaction ID"))
			}
			interactionResponse.Data.Components = []discordgo.MessageComponent{
				// ActionRow is a container of all buttons within the same row.
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							// Label is what the user will see on the button.
							Label: "Post",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: discordgo.SuccessButton,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: noModifier,
						},
					},
				},
			}
			if err = s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
				b.logger.Error("failed to respond", zap.Error(err))
			}
			return
		}

		// normal audio response...

		interactionResponse, err, cleanup := b.audioFileResponse(dialog, ids[0], pos, ContentModifierNone, username)
		if err != nil {
			b.respondError(s, i, err)
			return
		}
		defer cleanup()

		withAudioID, err := encodeCustomID("scrimp_confirm", ids[0], pos, ContentModifierNone)
		if err != nil {
			b.logger.Error("failed to marshal customID", zap.Error(err))
			b.respondError(s, i, fmt.Errorf("failed to create interaction ID"))
		}
		withoutAudioID, err := encodeCustomID("scrimp_confirm", ids[0], pos, ContentModifierTextOnly)
		if err != nil {
			b.logger.Error("failed to marshal customID", zap.Error(err))
			b.respondError(s, i, fmt.Errorf("failed to create interaction ID"))
		}
		withoutTextID, err := encodeCustomID("scrimp_confirm", ids[0], pos, ContentModifierAudioOnly)
		if err != nil {
			b.logger.Error("failed to marshal customID", zap.Error(err))
			b.respondError(s, i, fmt.Errorf("failed to create interaction ID"))
		}
		interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral
		interactionResponse.Data.Components = []discordgo.MessageComponent{
			// ActionRow is a container of all buttons within the same row.
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Post with audio",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SuccessButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: withAudioID,
					},
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Post without audio",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SecondaryButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: withoutAudioID,
					},
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Post without text",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SecondaryButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: withoutTextID,
					},
				},
			},
		}
		err = s.InteractionRespond(i.Interaction, interactionResponse)
		if err != nil {
			b.logger.Error("failed to respond", zap.Error(err))
		}
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()

		rawTerms := strings.TrimSpace(data.Options[0].StringValue())

		terms, err := searchterms.Parse(rawTerms)
		if err != nil {
			b.respondError(s, i, fmt.Errorf("invalid search terms: %w", err))
			return
		}
		if len(terms) == 0 {
			if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{},
				},
			}); err != nil {
				b.logger.Error("Failed to respond with autocomplete options", zap.Error(err))
			}
			return
		}

		filterString, err := filter.Print(searchterms.TermsToFilter(terms))
		if err != nil {
			b.respondError(s, i, fmt.Errorf("failed to create filter: %w", err))
			return
		}
		res, err := b.searchApiClient.PredictSearchTerm(
			context.Background(),
			&api.PredictSearchTermRequest{
				Query:          filterString,
				MaxPredictions: 25,
			},
		)
		if err != nil {
			b.logger.Error("Failed to fetch autocomplete options", zap.Error(err))
			return
		}

		choices := []*discordgo.ApplicationCommandOptionChoice{}
		for _, v := range res.Predictions {
			if v.Actor == "" {
				continue
			}
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  util.TrimToN(fmt.Sprintf("%s: %s", v.Actor, v.Line), 100),
				Value: fmt.Sprintf("%s:%d", v.Epid, v.Pos),
			})
		}
		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		}); err != nil {
			b.logger.Error("Failed to respond with autocomplete options", zap.Error(err))
		}
		return
	}
	b.respondError(s, i, fmt.Errorf("unknown command type"))
}

func (b *Bot) scrimptonQueryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	if customIDPayload == "" {
		return
	}

	customID, err := decodeCustomIDPayload(customIDPayload)
	if err != nil {
		b.logger.Error("failed to decode customID", zap.Error(err))
		b.respondError(s, i, fmt.Errorf("failed to decode customID"))
	}

	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}
	dialog, err := b.transcriptApiClient.GetTranscriptDialog(context.Background(), &api.GetTranscriptDialogRequest{
		Epid:            customID.EpisodeID,
		Pos:             int32(customID.Position),
		NumContextLines: 2,
	})
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to fetch selected line"))
		return
	}

	// respond gif
	if dialog.TranscriptMeta.MediaType == api.MediaType_VIDEO {
		interactionResponse, err, cleanup := b.videoFileResponse(dialog, customID.EpisodeID, customID.Position, username)
		if err != nil {
			b.respondError(s, i, err)
			return
		}
		defer cleanup()

		err = s.InteractionRespond(i.Interaction, interactionResponse)
		if err != nil {
			b.logger.Error("failed to respond", zap.Error(err))
		}
		return
	}

	// respond audio
	interactionResponse, err, cleanup := b.audioFileResponse(dialog, customID.EpisodeID, customID.Position, customID.ContentModifier, username)
	if err != nil {
		b.respondError(s, i, err)
		return
	}
	defer cleanup()

	if err = s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
		b.respondError(s, i, fmt.Errorf("unknown command type"))
	}
}

func (b *Bot) audioFileResponse(dialog *api.TranscriptDialog, episodeId string, pos int, contentModifier ContentModifier, username string) (*discordgo.InteractionResponse, error, func()) {

	dialogFormatted := strings.Builder{}
	var matchedDialogRow *api.Dialog
	for k, d := range dialog.Dialog {
		switch d.Type {
		case api.Dialog_CHAT:
			if d.Actor == "" {
				dialogFormatted.WriteString(fmt.Sprintf("\n> *%s*", d.Content))
			} else {
				if d.IsMatchedRow {
					dialogFormatted.WriteString(fmt.Sprintf("\n> **%s: %s**", d.Actor, d.Content))
				} else {
					dialogFormatted.WriteString(fmt.Sprintf("\n> %s: %s", d.Actor, d.Content))
				}
			}
		case api.Dialog_NONE:
			dialogFormatted.WriteString(fmt.Sprintf("\n> *%s*", d.Content))
		case api.Dialog_SONG:
			dialogFormatted.WriteString(fmt.Sprintf("\n> SONG: %s", d.Content))
		}
		if d.IsMatchedRow {
			matchedDialogRow = dialog.Dialog[k]
		}
	}
	if matchedDialogRow == nil {
		return nil, fmt.Errorf("no line was matched"), func() {}
	}

	var files []*discordgo.File
	cancelFunc := func() {}

	if contentModifier != ContentModifierTextOnly {
		audioFileURL := fmt.Sprintf("%s%s?pos=%d", b.webUrl, dialog.TranscriptMeta.AudioUri, matchedDialogRow.Pos)
		resp, err := http.Get(audioFileURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch selected line"), func() {}
		}
		if resp.StatusCode != http.StatusOK {
			b.logger.Error("failed to fetch audio", zap.Error(err), zap.String("url", audioFileURL), zap.Int("status_code", resp.StatusCode))
			return nil, fmt.Errorf("failed to fetch audio: %s", resp.Status), func() {}
		}
		files = append(files, &discordgo.File{
			Name:        createFileName(dialog, matchedDialogRow, "mp3"),
			ContentType: "audio/mpeg",
			Reader:      resp.Body,
		})
		cancelFunc = func() {
			resp.Body.Close()
		}
	}
	var content string
	if contentModifier != ContentModifierAudioOnly {
		content = fmt.Sprintf(
			"%s\n\n %s",
			dialogFormatted.String(),
			fmt.Sprintf(
				"`%s` @ `%s` | [%s](%s) | Posted by %s",
				dialog.TranscriptMeta.Id,
				(time.Millisecond*time.Duration(matchedDialogRow.OffsetMs)).String(),
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d", b.webUrl, episodeId, pos),
				username,
			),
		)
	} else {
		content = fmt.Sprintf("Posted by %s", username)
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Files:   files,
		},
	}, nil, cancelFunc
}

func (b *Bot) videoFileResponse(dialog *api.TranscriptDialog, episodeId string, pos int, username string) (*discordgo.InteractionResponse, error, func()) {

	var matchedDialogRow *api.Dialog
	for k, d := range dialog.Dialog {
		if d.IsMatchedRow {
			matchedDialogRow = dialog.Dialog[k]
			break
		}
	}
	if matchedDialogRow == nil {
		return nil, fmt.Errorf("no line was matched"), func() {}
	}

	var files []*discordgo.File

	fileURL := fmt.Sprintf("%s/dl/media/gif/%s?pos=%d", b.webUrl, dialog.TranscriptMeta.Id, matchedDialogRow.Pos)

	b.logger.Info("Fetching GIF", zap.String("url", fileURL))
	resp, err := http.Get(fileURL)
	if err != nil {
		b.logger.Error("failed to fetch GIF", zap.Error(err), zap.String("url", fileURL))
		return nil, fmt.Errorf("failed to fetch GIF: %w", err), func() {}
	}
	if resp.StatusCode != http.StatusOK {
		b.logger.Error("failed to fetch GIF", zap.Error(err), zap.String("url", fileURL), zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to fetch GIF: %s", resp.Status), func() {}
	}

	files = append(files, &discordgo.File{
		Name:        createFileName(dialog, matchedDialogRow, "gif"),
		ContentType: "image/gif",
		Reader:      resp.Body,
	})
	cancelFunc := func() {
		resp.Body.Close()
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"`%s` @ `%s` | [%s](%s) | Posted by %s",
				dialog.TranscriptMeta.Id,
				(time.Millisecond * time.Duration(matchedDialogRow.OffsetMs)).String(),
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d", b.webUrl, episodeId, pos),
				username,
			),
			Files: files,
		},
	}, nil, cancelFunc
}

func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	b.logger.Error("Error response was sent", zap.Error(err))
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Failed to fetch quote due to error: %s", err.Error()),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
		return
	}
}

func encodeCustomID(action string, episodeID string, position int, contentModifier ContentModifier) (string, error) {
	encoded, err := json.Marshal(&CustomID{
		EpisodeID:       episodeID,
		Position:        position,
		ContentModifier: contentModifier,
	})
	return fmt.Sprintf("%s:%s", action, string(encoded)), err
}

func decodeCustomIDPayload(data string) (*CustomID, error) {
	decoded := &CustomID{}
	return decoded, json.Unmarshal([]byte(data), decoded)
}

func createFileName(dialog *api.TranscriptDialog, matchedDialogRow *api.Dialog, suffix string) string {
	if contentFilename := contentToFilename(matchedDialogRow.Content); contentFilename != "" {
		return fmt.Sprintf("%s.%s", contentFilename, suffix)
	}
	return fmt.Sprintf("%s-%d.%s", dialog.TranscriptMeta.Id, matchedDialogRow.Pos, suffix)
}

func contentToFilename(rawContent string) string {
	rawContent = punctuation.ReplaceAllString(rawContent, "")
	rawContent = spaces.ReplaceAllString(rawContent, " ")
	rawContent = strings.TrimSpace(rawContent)
	split := strings.Split(rawContent, " ")
	if len(split) > 9 {
		split = split[:8]
	}
	return strings.Join(split, "-")
}
