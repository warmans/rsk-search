package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/searchterms"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var punctuation = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
var spaces = regexp.MustCompile(`[\s]{2,}`)
var metaWhitespace = regexp.MustCompile(`[\n\r\t]+`)

var rendersInProgress = map[string]string{}
var renderMutex = sync.RWMutex{}
var errRenderInProgress = errors.New("render in progress")
var errDuplicateInteraction = errors.New("interaction already processing")

func lockRenderer(username string, interactionIdentifier string) (func(), error) {
	renderMutex.Lock()
	defer renderMutex.Unlock()
	if oldInteractionID, found := rendersInProgress[username]; found {
		if interactionIdentifier == oldInteractionID {
			return func() {}, errDuplicateInteraction
		}
		return func() {}, errRenderInProgress
	}
	rendersInProgress[username] = interactionIdentifier
	return func() {
		renderMutex.Lock()
		delete(rendersInProgress, username)
		renderMutex.Unlock()
	}, nil
}

type CustomID struct {
	EpisodeID       string          `json:"e,omitempty"`
	Position        int             `json:"p,omitempty"`
	ContentModifier ContentModifier `json:"t,omitempty"`
}

func (c CustomID) String() string {
	return fmt.Sprintf("%s-%d-%d", c.EpisodeID, c.Position, c.ContentModifier)
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
		"scrimp_custom":  bot.editModal,
	}
	bot.modalHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, suffix string){
		"scrimp_confirm": bot.scrimptonQueryCompleteCustom,
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
	modalHandlers       map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, customIdPayload string)
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
		case discordgo.InteractionModalSubmit:
			// prefix match buttons to allow additional data in the customID
			for k, h := range b.modalHandlers {
				actionPrefix := fmt.Sprintf("%s:", k)
				if strings.HasPrefix(i.ModalSubmitData().CustomID, actionPrefix) {
					h(s, i, strings.TrimPrefix(i.ModalSubmitData().CustomID, actionPrefix))
					return
				}
			}
			b.respondError(s, i, fmt.Errorf("unknown customID format: %s", i.MessageComponentData().CustomID))
			return
		case discordgo.InteractionMessageComponent:
			// prefix match buttons to allow additional data in the customID
			for k, h := range b.buttonHandlers {
				actionPrefix := fmt.Sprintf("%s:", k)
				if strings.HasPrefix(i.MessageComponentData().CustomID, actionPrefix) {
					h(s, i, strings.TrimPrefix(i.MessageComponentData().CustomID, actionPrefix))
					return
				}
			}
			b.respondError(s, i, fmt.Errorf("unknown customID format: %s", i.MessageComponentData().CustomID))
			return
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
		customID, err := parseCustomID(selection)
		if err != nil {
			b.respondError(s, i, err)
			return
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

		switch dialog.TranscriptMeta.MediaType {
		case api.MediaType_VIDEO:
			if err := b.beginVideoResponse(s, i, dialog, *customID, username); err != nil {
				b.logger.Error("Failed to begin video response", zap.Error(err))
			}
			return
		case api.MediaType_AUDIO:
			if err := b.beginAudioResponse(s, i, dialog, *customID, username); err != nil {
				b.logger.Error("Failed to begin video response", zap.Error(err))
			}
			return
		}
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
		data := i.ApplicationCommandData()

		rawTerms := strings.TrimSpace(data.Options[0].StringValue())

		terms, err := searchterms.Parse(rawTerms)
		if err != nil {
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

func (b *Bot) beginVideoResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	dialog *api.TranscriptDialog,
	customID CustomID,
	username string,
) error {
	// send a placeholder
	interactionResponse, cleanup, err := b.buildVideoResponse(dialog, customID, username, true, nil)
	defer cleanup()
	if err != nil {
		if errors.Is(err, errDuplicateInteraction) {
			fmt.Println("Duplicated interaction")
			return nil
		}
		b.respondError(s, i, err)
		return err
	}
	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral
	if err = s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
		return err
	}

	// update with the gif
	go func() {
		interactionResponse, cleanup, err = b.buildVideoResponse(dialog, customID, username, false, nil)
		defer cleanup()
		if err != nil {
			if errors.Is(err, errDuplicateInteraction) {
				return
			}
			b.logger.Error("interaction failed", zap.Error(err))
			_, err := s.InteractionResponseEdit(
				i.Interaction,
				&discordgo.WebhookEdit{
					Content: util.ToPtr(fmt.Sprintf("Failed (%s)...", err.Error())),
				},
			)
			if err != nil {
				b.logger.Error("edit failed", zap.Error(err))
			}
			return
		}

		noModifier, err := encodeCustomID("scrimp_confirm", customID.EpisodeID, customID.Position, ContentModifierNone)
		if err != nil {
			b.logger.Error("edit failed", zap.Error(fmt.Errorf("failed to marshal customID: %w", err)))
			return
		}
		customise, err := encodeCustomID("scrimp_custom", customID.EpisodeID, customID.Position, ContentModifierNone)
		if err != nil {
			b.logger.Error("edit failed", zap.Error(fmt.Errorf("failed to marshal customID: %w", err)))
			return
		}
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: util.ToPtr("Complete!" + interactionResponse.Data.Content),
			Components: util.ToPtr([]discordgo.MessageComponent{
				// ActionRow is a container of all buttons within the same row.
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							// Label is what the user will see on the button.
							Label: "Post",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: discordgo.PrimaryButton,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: noModifier,
						},
						discordgo.Button{
							// Label is what the user will see on the button.
							Label: "Post Custom",
							// Style provides coloring of the button. There are not so many styles tho.
							Style: discordgo.SecondaryButton,
							// Disabled allows bot to disable some buttons for users.
							Disabled: false,
							// CustomID is a thing telling Discord which data to send when this button will be pressed.
							CustomID: customise,
						},
					},
				},
			}),
			Files: interactionResponse.Data.Files,
		})
		if err != nil {
			b.logger.Error("edit failed", zap.Error(err))
			return
		}
	}()
	return nil
}

func (b *Bot) editModal(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {

	customID, err := decodeCustomIDPayload(customIDPayload)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to decode customID: %w", err))
		return
	}

	noModifier, err := encodeCustomID("scrimp_confirm", customID.EpisodeID, customID.Position, ContentModifierNone)
	if err != nil {
		b.logger.Error("edit failed", zap.Error(fmt.Errorf("failed to marshal customID: %w", err)))
		return
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: noModifier,
			Title:    "Edit Gif",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "custom_text",
							Label:     "Gif Text",
							Style:     discordgo.TextInputParagraph,
							Required:  true,
							MaxLength: 128,
						},
					},
				},
			},
		},
	})
	if err != nil {
		b.respondError(s, i, err)
	}
}

func (b *Bot) beginAudioResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	dialog *api.TranscriptDialog,
	customID CustomID,
	username string,
) error {

	interactionResponse, err, cleanup := b.audioFileResponse(dialog, customID, username)
	if err != nil {
		b.respondError(s, i, err)
		return err
	}
	defer cleanup()

	withAudioID, err := encodeCustomID("scrimp_confirm", customID.EpisodeID, customID.Position, ContentModifierNone)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to create interaction ID: %w", err))
		return err
	}
	withoutAudioID, err := encodeCustomID("scrimp_confirm", customID.EpisodeID, customID.Position, ContentModifierTextOnly)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to create interaction ID: %w", err))
		return err
	}
	withoutTextID, err := encodeCustomID("scrimp_confirm", customID.EpisodeID, customID.Position, ContentModifierAudioOnly)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to create interaction ID: %w", err))
		return err
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
	return nil
}

func (b *Bot) scrimptonQueryCompleteCustom(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {
	if customIDPayload == "" {
		b.respondError(s, i, fmt.Errorf("missing customID"))
		return
	}
	customID, err := decodeCustomIDPayload(customIDPayload)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to decode customID: %w", err))
		return
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
	customText := i.Interaction.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	if err = b.completeVideoResponse(s, i, dialog, username, *customID, util.ToPtr(customText)); err != nil {
		b.respondError(s, i, err)
	}
}

func (b *Bot) scrimptonQueryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	spew.Dump(customIDPayload, i)

	if customIDPayload == "" {
		b.respondError(s, i, fmt.Errorf("missing customID"))
		return
	}
	customID, err := decodeCustomIDPayload(customIDPayload)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to decode customID: %w", err))
		return
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

	switch dialog.TranscriptMeta.MediaType {
	case api.MediaType_VIDEO:
		if err := b.completeVideoResponse(s, i, dialog, username, *customID, nil); err != nil {
			b.logger.Error("Failed to complete video response", zap.Error(err))
		}
	case api.MediaType_AUDIO:
		// respond audio
		interactionResponse, err, cleanup := b.audioFileResponse(dialog, *customID, username)
		defer cleanup()
		if err != nil {
			b.respondError(s, i, err)
			return
		}
		if err = s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
			b.respondError(s, i, err)
		}
	}
}

func (b *Bot) completeVideoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, dialog *api.TranscriptDialog, username string, customID CustomID, customText *string) error {

	interactionResponse, cleanup, err := b.buildVideoResponse(dialog, customID, username, true, nil)
	defer cleanup()
	if err != nil {
		if errors.Is(err, errDuplicateInteraction) {
			fmt.Println("Duplicated interaction")
			return nil
		}
		if errors.Is(err, errRenderInProgress) {
			b.respondError(s, i, errors.New("you already have a render in progress"))
		}
		b.respondError(s, i, err)
		return err
	}
	if err = s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		return fmt.Errorf("failed to respond: %w", err)
	}
	go func() {
		interactionResponse, cleanup, err = b.buildVideoResponse(dialog, customID, username, false, customText)
		defer cleanup()
		if err != nil {
			if errors.Is(err, errDuplicateInteraction) {
				return
			}
			b.logger.Error("interaction failed", zap.Error(err))
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: util.ToPtr("Failed....")})
			if err != nil {
				b.logger.Error("edit failed", zap.Error(err))
			}
			return
		}
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: util.ToPtr(interactionResponse.Data.Content),
			Files:   interactionResponse.Data.Files,
		})
		if err != nil {
			b.logger.Error("edit failed", zap.Error(err))
		}
	}()

	return nil
}

func (b *Bot) audioFileResponse(dialog *api.TranscriptDialog, customID CustomID, username string) (*discordgo.InteractionResponse, error, func()) {

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

	if customID.ContentModifier != ContentModifierTextOnly {
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
	if customID.ContentModifier != ContentModifierAudioOnly {
		content = fmt.Sprintf(
			"%s\n\n %s",
			dialogFormatted.String(),
			fmt.Sprintf(
				"`%s` @ `%s` | [%s](%s) | Posted by %s",
				dialog.TranscriptMeta.Id,
				(time.Millisecond*time.Duration(matchedDialogRow.OffsetMs)).String(),
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d", b.webUrl, customID.EpisodeID, customID.Position),
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

func (b *Bot) buildVideoResponse(dialog *api.TranscriptDialog, customID CustomID, username string, placeholder bool, customText *string) (*discordgo.InteractionResponse, func(), error) {
	cleanup, err := lockRenderer(username, customID.String())
	defer cleanup()
	if err != nil {
		if errors.Is(err, errDuplicateInteraction) {
			return nil, func() {}, errDuplicateInteraction
		}
		if errors.Is(err, errRenderInProgress) {
			return nil, func() {}, errRenderInProgress
		}
		return nil, func() {}, err
	}

	var matchedDialogRow *api.Dialog
	for k, d := range dialog.Dialog {
		if d.IsMatchedRow {
			matchedDialogRow = dialog.Dialog[k]
			break
		}
	}
	if matchedDialogRow == nil {
		return nil, func() {}, fmt.Errorf("no line was matched")
	}

	var files []*discordgo.File

	customTextParam := ""
	if customText != nil {
		customTextParam = fmt.Sprintf("&custom_text=%s", url.QueryEscape(*customText))
	}
	fileURL := fmt.Sprintf("%s/dl/media/gif/%s?pos=%d%s", b.webUrl, dialog.TranscriptMeta.Id, matchedDialogRow.Pos, customTextParam)
	cancelFunc := func() {}
	bodyText := ""

	if !placeholder {
		b.logger.Info("Fetching GIF", zap.String("url", fileURL))
		resp, err := http.Get(fileURL)
		if err != nil {
			b.logger.Error("failed to fetch GIF", zap.Error(err), zap.String("url", fileURL))
			return nil, func() {}, fmt.Errorf("failed to fetch GIF: %w", err)
		}
		cancelFunc = func() {
			resp.Body.Close()
		}
		if resp.StatusCode != http.StatusOK {
			b.logger.Error("failed to fetch GIF", zap.Error(err), zap.String("url", fileURL), zap.Int("status_code", resp.StatusCode))
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, cancelFunc, fmt.Errorf("failed to fetch GIF: %s", resp.Status)
			}
			return nil, cancelFunc, fmt.Errorf("failed to fetch GIF: %s", string(body))
		}
		files = append(files, &discordgo.File{
			Name:        createFileName(dialog, matchedDialogRow, "gif"),
			ContentType: "image/gif",
			Reader:      resp.Body,
		})

	} else {
		bodyText = ":timer: Rendering gif..."
	}
	editLabel := ""
	if customText != nil {
		editLabel = " (edited)"
	}
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				"%s\n\n`%s` @ `%s`%s | [%s](%s) | Posted by %s",
				bodyText,
				dialog.TranscriptMeta.Id,
				(time.Millisecond * time.Duration(matchedDialogRow.OffsetMs)).String(),
				editLabel,
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d", b.webUrl, customID.EpisodeID, customID.Position),
				username,
			),
			Files: files,
		},
	}, cancelFunc, nil
}

func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	b.logger.Error("Error response was sent", zap.Error(err))
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Request failed with error: %s", err.Error()),
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

func parseCustomID(raw string) (*CustomID, error) {
	ids := strings.Split(raw, ":")
	if len(ids) < 2 {
		return nil, fmt.Errorf("invalid selection")
	}
	pos, err := strconv.Atoi(ids[1])
	if err != nil {
		return nil, fmt.Errorf("invalid position %s: %w", ids[1], err)
	}
	return &CustomID{
		EpisodeID: ids[0],
		Position:  pos,
	}, nil
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
	rawContent = metaWhitespace.ReplaceAllString(rawContent, " ")
	rawContent = strings.ToLower(strings.TrimSpace(rawContent))
	split := strings.Split(rawContent, " ")
	if len(split) > 9 {
		split = split[:8]
	}
	return strings.Join(split, "-")
}
