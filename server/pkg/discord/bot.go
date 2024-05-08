package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	buttonHandlers      map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, customIdSuffix string)
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
				if strings.HasPrefix(i.MessageComponentData().CustomID, fmt.Sprintf("%s:", k)) {
					h(s, i, strings.TrimPrefix(i.MessageComponentData().CustomID, fmt.Sprintf("%s:", k)))
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

		interactionResponse, err, cleanup := b.queryInteractionResponse(ids[0], pos, false)
		if err != nil {
			b.respondError(s, i, err)
			return
		}
		defer cleanup()

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
						CustomID: fmt.Sprintf("scrimp_confirm:%s", selection),
					},
					discordgo.Button{
						// Label is what the user will see on the button.
						Label: "Post without audio",
						// Style provides coloring of the button. There are not so many styles tho.
						Style: discordgo.SecondaryButton,
						// Disabled allows bot to disable some buttons for users.
						Disabled: false,
						// CustomID is a thing telling Discord which data to send when this button will be pressed.
						CustomID: fmt.Sprintf("scrimp_confirm:%s:textonly", selection),
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

		prefixTerm := strings.TrimSpace(data.Options[0].StringValue())
		exactMatch := false

		// looks like a quoted string =
		if strings.HasPrefix(prefixTerm, `"`) {
			exactMatch = true
		}

		res, err := b.searchApiClient.PredictSearchTerm(
			context.Background(),
			&api.PredictSearchTermRequest{
				Prefix:         data.Options[0].StringValue(),
				MaxPredictions: 25,
				Exact:          exactMatch,
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

func (b *Bot) scrimptonQueryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, commandSuffix string) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	if commandSuffix == "" {
		return
	}
	state := strings.Split(commandSuffix, ":")
	if len(state) < 2 {
		b.logger.Error("unexpected selection format", zap.String("id_suffix", commandSuffix))
		b.respondError(s, i, fmt.Errorf("invalid selection"))
		return
	}
	pos, err := strconv.Atoi(state[1])
	if err != nil {
		b.logger.Error("failed to parse position", zap.String("pos", state[1]), zap.Error(err))
		b.respondError(s, i, fmt.Errorf("internal error, unknown selection"))
		return
	}
	var omitAudio bool
	if len(state) > 2 && state[2] == "textonly" {
		omitAudio = true
	}
	interactionResponse, err, cleanup := b.queryInteractionResponse(state[0], pos, omitAudio)
	if err != nil {
		b.respondError(s, i, err)
		return
	}
	defer cleanup()

	err = s.InteractionRespond(i.Interaction, interactionResponse)
	if err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
		b.respondError(s, i, fmt.Errorf("unknown command type"))
	}
}

func (b *Bot) queryInteractionResponse(episodeId string, pos int, omitAudio bool) (*discordgo.InteractionResponse, error, func()) {
	dialog, err := b.transcriptApiClient.GetTranscriptDialog(context.Background(), &api.GetTranscriptDialogRequest{
		Epid:            episodeId,
		Pos:             int32(pos),
		NumContextLines: 2,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch selected line"), func() {}
	}

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

	var files []*discordgo.File
	cancelFunc := func() {}
	if !omitAudio {
		audioFileURL := fmt.Sprintf("%s%s?pos=%d", b.webUrl, dialog.TranscriptMeta.AudioUri, matchedDialogRow.Pos)
		resp, err := http.Get(audioFileURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch selected line"), func() {}
		}
		files = append(files, &discordgo.File{
			Name:        fmt.Sprintf("%s-%d.mp3", dialog.TranscriptMeta.Id, matchedDialogRow.Pos),
			ContentType: "audio/mpeg",
			Reader:      resp.Body,
		})
		cancelFunc = func() {
			resp.Body.Close()
		}
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"%s\n\n %s",
				dialogFormatted.String(),
				fmt.Sprintf(
					"`%s` @ `%s` | [%s](%s)",
					dialog.TranscriptMeta.Id,
					(time.Second*time.Duration(matchedDialogRow.OffsetSec)).String(),
					strings.TrimPrefix(b.webUrl, "https://"),
					b.webUrl,
				),
			),
			Files: files,
		},
	}, nil, cancelFunc
}

func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
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
