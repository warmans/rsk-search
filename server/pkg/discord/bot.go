package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/archive"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/searchterms"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultContext = 0

var punctuation = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
var spaces = regexp.MustCompile(`[\s]{2,}`)
var metaWhitespace = regexp.MustCompile(`[\n\r\t]+`)

type customIDOpt func(c *CustomID)

func withModifier(mod ContentModifier) customIDOpt {
	return func(c *CustomID) {
		c.ContentModifier = mod
	}
}

func withStartLine(pos int32) customIDOpt {
	return func(c *CustomID) {
		c.StartLine = pos
	}
}
func withEndLine(pos int32) customIDOpt {
	return func(c *CustomID) {
		c.EndLine = pos
	}
}

type CustomID struct {
	EpisodeID       string          `json:"e,omitempty"`
	StartLine       int32           `json:"s,omitempty"`
	EndLine         int32           `json:"f,omitempty"`
	NumContextLines int             `json:"c,omitempty"`
	ContentModifier ContentModifier `json:"t,omitempty"`
}

func (c CustomID) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		// this should never happen
		fmt.Printf("failed to encode customID: %s\n", err.Error())
		return ""
	}
	return string(data)
}

func (c CustomID) withOption(options ...customIDOpt) CustomID {
	clone := &CustomID{
		EpisodeID:       c.EpisodeID,
		StartLine:       c.StartLine,
		EndLine:         c.EndLine,
		NumContextLines: c.NumContextLines,
		ContentModifier: c.ContentModifier,
	}
	for _, v := range options {
		v(clone)
	}
	return *clone
}

type ContentModifier uint8

const (
	ContentModifierNone ContentModifier = iota
	ContentModifierTextOnly
	ContentModifierAudioOnly
	ContentModifierGifOnly
)

func NewBot(
	logger *zap.Logger,
	session *discordgo.Session,
	guildID string,
	webUrl string,
	archiveStore *archive.Store,
	transcriptApiClient api.TranscriptServiceClient,
	searchApiClient api.SearchServiceClient,
) *Bot {
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentMessageContent)

	bot := &Bot{
		logger:              logger,
		session:             session,
		guildID:             guildID,
		webUrl:              webUrl,
		archiveStore:        archiveStore,
		transcriptApiClient: transcriptApiClient,
		searchApiClient:     searchApiClient,
		rewindThreadCache:   &sync.Map{},
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
			{
				Name:        "rewind",
				Description: "Create a rewind thread",
				Type:        discordgo.ChatApplicationCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:         "episode",
						Description:  "The episode ID",
						Type:         discordgo.ApplicationCommandOptionString,
						Required:     true,
						Autocomplete: true,
					},
				},
			},
			{
				Name: "scrimp-archive",
				Type: discordgo.MessageApplicationCommand,
			},
		},
	}
	bot.commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"scrimp":         bot.queryBegin,
		"rewind":         bot.rewindBegin,
		"scrimp-archive": bot.quickArchiveModalOpen,
	}
	bot.buttonHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, suffix string){
		"cfm":            bot.queryComplete,
		"up":             bot.updatePreview,
		"rewind-start":   bot.rewindStart,
		"episode-rating": bot.rateEpisode,
	}
	bot.modalHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, suffix string){
		"quick-archive-modal-save": bot.quickArchiveModalSave,
	}

	return bot
}

type Bot struct {
	logger              *zap.Logger
	session             *discordgo.Session
	guildID             string
	webUrl              string
	archiveStore        *archive.Store
	transcriptApiClient api.TranscriptServiceClient
	searchApiClient     api.SearchServiceClient
	commands            []*discordgo.ApplicationCommand
	commandHandlers     map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	buttonHandlers      map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, customIdPayload string)
	modalHandlers       map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, customIdPayload string)
	createdCommands     []*discordgo.ApplicationCommand
	rewindStateLock     sync.RWMutex
	rewindThreadCache   *sync.Map
}

func (b *Bot) Start() error {

	if err := b.initRewindCache(); err != nil {
		return err
	}

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
			b.respondError(s, i, fmt.Errorf("unknown customID format: %s", i.ModalSubmitData().CustomID))
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
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		b.handleMessageCreate(s, m)
	})

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		b.handleMessageReactAdd(s, r)
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

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if _, isValidThread := b.rewindThreadCache.Load(m.ChannelID); isValidThread {
		spew.Dump(m)
	}
}

func (b *Bot) handleMessageReactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if _, isValidThread := b.rewindThreadCache.Load(r.ChannelID); isValidThread {
		fmt.Println("reaction to message: ", r.MessageID, " channel: ", r.ChannelID)
	}
}

func (b *Bot) queryBegin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		selection := i.ApplicationCommandData().Options[0].StringValue()
		if selection == "" {
			return
		}
		customID, err := decodeCustomIDPayload(selection)
		if err != nil {
			b.respondError(s, i, err)
			return
		}
		if err := b.beginAudioResponse(s, i, customID); err != nil {
			b.respondError(s, i, err)
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
				Name: util.TrimToN(fmt.Sprintf("%s: %s", v.Actor, v.Line), 100),
				Value: (&CustomID{
					EpisodeID:       v.Epid,
					StartLine:       v.Pos,
					EndLine:         v.Pos,
					NumContextLines: defaultContext,
					ContentModifier: ContentModifierTextOnly,
				}).String(),
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
func (b *Bot) rewindBegin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		selection := i.ApplicationCommandData().Options[0].StringValue()
		if selection == "" {
			return
		}

		b.confirmRewindStart(s, i, selection)
		return
	case discordgo.InteractionApplicationCommandAutocomplete:
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
			b.logger.Error("Failed to respond with autocomplete options", zap.Error(err))
		}
		return
	}
	b.respondError(s, i, fmt.Errorf("unknown command type"))
}

func (b *Bot) updatePreview(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	customIDPayload string,
) {
	customID, err := decodeCustomIDPayload(customIDPayload)
	if err != nil {
		b.respondError(s, i, err)
		return
	}
	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}

	interactionResponse, maxDialogOffset, err, cleanup := b.audioFileResponse(customID, username)
	if err != nil {
		b.respondError(s, i, err)
		return
	}
	defer cleanup()

	interactionResponse.Data.Components = b.buttons(customID, maxDialogOffset)
	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: interactionResponse.Data,
	}); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) beginAudioResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	customID CustomID,
) error {
	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}

	interactionResponse, maxDialogOffset, err, cleanup := b.audioFileResponse(customID, username)
	if err != nil {
		b.respondError(s, i, err)
		return err
	}
	defer cleanup()

	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral
	interactionResponse.Data.Components = b.buttons(customID, maxDialogOffset)
	err = s.InteractionRespond(i.Interaction, interactionResponse)
	if err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
	}
	return nil
}

func (b *Bot) buttons(customID CustomID, maxDialogOffset int32) []discordgo.MessageComponent {

	editRow1 := []discordgo.MessageComponent{}
	if customID.StartLine > 0 {
		editRow1 = append(editRow1, discordgo.Button{
			// Label is what the user will see on the button.
			Label: "Shift Dialog Backwards",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏪",
			},
			// Style provides coloring of the button. There are not so many styles tho.
			Style: discordgo.SecondaryButton,
			// CustomID is a thing telling Discord which data to send when this button will be pressed.
			CustomID: encodeCustomIDForAction(
				"up",
				customID.withOption(
					withStartLine(customID.StartLine-1),
					withEndLine(customID.EndLine-1),
				),
			),
		})
	}
	if customID.StartLine+1 < maxDialogOffset {
		editRow1 = append(editRow1, discordgo.Button{
			// Label is what the user will see on the button.
			Label: "Shift Dialog Forward",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏩",
			},
			// Style provides coloring of the button. There are not so many styles tho.
			Style: discordgo.SecondaryButton,
			// CustomID is a thing telling Discord which data to send when this button will be pressed.
			CustomID: encodeCustomIDForAction(
				"up",
				customID.withOption(
					withStartLine(customID.StartLine+1),
					withEndLine(min(maxDialogOffset, customID.EndLine+1)),
				),
			),
		})
	}
	if customID.EndLine-customID.StartLine < 25 && customID.ContentModifier != ContentModifierGifOnly {
		if customID.StartLine > 0 {
			editRow1 = append(editRow1, discordgo.Button{
				// Label is what the user will see on the button.
				Label: "Add Previous Line",
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				// Style provides coloring of the button. There are not so many styles tho.
				Style: discordgo.SecondaryButton,
				// CustomID is a thing telling Discord which data to send when this button will be pressed.
				CustomID: encodeCustomIDForAction(
					"up",
					customID.withOption(
						withStartLine(customID.StartLine-1),
					),
				),
			})
		}
		if customID.EndLine+1 < maxDialogOffset {
			editRow1 = append(editRow1, discordgo.Button{
				// Label is what the user will see on the button.
				Label: "Add Next Line",
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				// Style provides coloring of the button. There are not so many styles tho.
				Style: discordgo.SecondaryButton,
				// CustomID is a thing telling Discord which data to send when this button will be pressed.
				CustomID: encodeCustomIDForAction(
					"up",
					customID.withOption(
						withEndLine(customID.EndLine+1),
					),
				),
			})
		}
	}

	editRow2 := []discordgo.MessageComponent{}
	if customID.EndLine-customID.StartLine > 0 {
		editRow2 = append(editRow2, discordgo.Button{
			// Label is what the user will see on the button.
			Label: "Trim First Line",
			Emoji: &discordgo.ComponentEmoji{
				Name: "✂",
			},
			// Style provides coloring of the button. There are not so many styles tho.
			Style: discordgo.SecondaryButton,
			// CustomID is a thing telling Discord which data to send when this button will be pressed.
			CustomID: encodeCustomIDForAction(
				"up",
				customID.withOption(
					withStartLine(customID.StartLine+1),
				),
			),
		})
		editRow2 = append(editRow2, discordgo.Button{
			// Label is what the user will see on the button.
			Label: "Trim Last Line",
			Emoji: &discordgo.ComponentEmoji{
				Name: "✂",
			},
			// Style provides coloring of the button. There are not so many styles tho.
			Style: discordgo.SecondaryButton,
			// CustomID is a thing telling Discord which data to send when this button will be pressed.
			CustomID: encodeCustomIDForAction(
				"up",
				customID.withOption(
					withEndLine(customID.EndLine-1),
				),
			),
		})
	}

	buttons := []discordgo.MessageComponent{}
	if len(editRow1) > 0 {
		buttons = append(buttons, discordgo.ActionsRow{
			Components: editRow1,
		})
	}
	if len(editRow2) > 0 {
		buttons = append(buttons, discordgo.ActionsRow{
			Components: editRow2,
		})
	}

	postButtons := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Post",
				Style:    discordgo.PrimaryButton,
				CustomID: encodeCustomIDForAction("cfm", customID),
			},
		},
	}
	if customID.ContentModifier != ContentModifierGifOnly {
		if customID.ContentModifier != ContentModifierNone {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label:    "Audio & Text",
				Style:    discordgo.SecondaryButton,
				CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierNone))),
			})
		}
		if customID.ContentModifier != ContentModifierAudioOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Audio Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "🔊",
				},
				Style:    discordgo.SecondaryButton,
				CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierAudioOnly))),
			})
		}
		if customID.ContentModifier != ContentModifierTextOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Text Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "📄",
				},
				Style:    discordgo.SecondaryButton,
				CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierTextOnly))),
			})
		}
	}
	if customID.StartLine == customID.EndLine && customID.NumContextLines == 0 {
		if customID.ContentModifier != ContentModifierGifOnly {
			postButtons.Components = append(postButtons.Components,
				discordgo.Button{
					Label: "GIF mode",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📺",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierGifOnly))),
				},
			)
		} else {
			postButtons.Components = append(postButtons.Components,
				discordgo.Button{
					Label: "Normal mode",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📻",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierTextOnly))),
				},
				discordgo.Button{
					Label: "Randomize image",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📺",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierGifOnly))),
				},
			)
		}
	}

	buttons = append(buttons, postButtons)

	return buttons
}

func (b *Bot) queryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	// can we get the files of the existing message?
	var files []*discordgo.File
	if len(i.Message.Attachments) > 0 {
		attachment := i.Message.Attachments[0]
		image, err := http.Get(attachment.URL)
		if err != nil {
			b.respondError(s, i, fmt.Errorf("failed to get original message attachment: %w", err))
			return
		}
		defer image.Body.Close()

		files = append(files, &discordgo.File{
			Name:        attachment.Filename,
			Reader:      image.Body,
			ContentType: attachment.ContentType,
		})
	}

	interactionResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:     i.Message.Content,
			Files:       files,
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}

	if err := s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) audioFileResponse(customID CustomID, username string) (*discordgo.InteractionResponse, int32, error, func()) {

	dialog, err := b.transcriptApiClient.GetTranscriptDialog(context.Background(), &api.GetTranscriptDialogRequest{
		Epid: customID.EpisodeID,
		Range: &api.DialogRange{
			Start: customID.StartLine,
			End:   customID.EndLine,
		},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch selected line"), func() {}
	}

	dialogFormatted := strings.Builder{}
	for _, d := range dialog.Dialog {
		switch d.Type {
		case api.Dialog_CHAT:
			if d.Actor == "" {
				dialogFormatted.WriteString(fmt.Sprintf("\n> *%s*", d.Content))
			} else {
				if d.IsMatchedRow {
					dialogFormatted.WriteString(fmt.Sprintf("\n> **%s:** %s", d.Actor, d.Content))
				} else {
					dialogFormatted.WriteString(fmt.Sprintf("\n> **%s:** %s", d.Actor, d.Content))
				}
			}
		case api.Dialog_NONE:
			dialogFormatted.WriteString(fmt.Sprintf("\n> *%s*", d.Content))
		case api.Dialog_SONG:
			dialogFormatted.WriteString(fmt.Sprintf("\n> **SONG:** %s", d.Content))
		}
	}

	var content string
	var files []*discordgo.File
	cancelFunc := func() {}

	if customID.ContentModifier == ContentModifierGifOnly {
		audioFileURL := fmt.Sprintf(
			"%s/dl/media/%s.gif?ts=%d-%d",
			b.webUrl,
			dialog.TranscriptMeta.ShortId,
			dialog.Dialog[0].OffsetMs,
			dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs,
		)
		resp, err := http.Get(audioFileURL)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to fetch selected line"), func() {}
		}
		if resp.StatusCode != http.StatusOK {
			b.logger.Error("failed to fetch gif", zap.Error(err), zap.String("url", audioFileURL), zap.Int("status_code", resp.StatusCode))
			return nil, 0, fmt.Errorf("failed to fetch gif: %s", resp.Status), func() {}
		}
		files = append(files, &discordgo.File{
			Name:        createFileName(dialog, "gif"),
			ContentType: "image/gif",
			Reader:      resp.Body,
		})
		cancelFunc = func() {
			resp.Body.Close()
		}

		if customID.ContentModifier != ContentModifierAudioOnly {
			content = fmt.Sprintf(
				"`%s` @ `%s - %s` | [%s](%s) | Posted by %s",
				dialog.TranscriptMeta.Id,
				(time.Duration(dialog.Dialog[0].OffsetMs)).String(),
				(time.Duration(dialog.Dialog[len(dialog.Dialog)-1].OffsetMs + dialog.Dialog[len(dialog.Dialog)-1].DurationMs)).String(),
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d-%d", b.webUrl, customID.EpisodeID, customID.StartLine, customID.EndLine),
				username,
			)
		} else {
			content = fmt.Sprintf("Posted by %s", username)
		}

	} else {
		if customID.ContentModifier != ContentModifierTextOnly {
			audioFileURL := fmt.Sprintf(
				"%s/dl/media/%s.mp3?ts=%d-%d",
				b.webUrl,
				dialog.TranscriptMeta.ShortId,
				dialog.Dialog[0].OffsetMs,
				dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs,
			)
			resp, err := http.Get(audioFileURL)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to fetch selected line"), func() {}
			}
			if resp.StatusCode != http.StatusOK {
				b.logger.Error("failed to fetch audio", zap.Error(err), zap.String("url", audioFileURL), zap.Int("status_code", resp.StatusCode))
				return nil, 0, fmt.Errorf("failed to fetch audio: %s", resp.Status), func() {}
			}
			files = append(files, &discordgo.File{
				Name:        createFileName(dialog, "mp3"),
				ContentType: "audio/mpeg",
				Reader:      resp.Body,
			})
			cancelFunc = func() {
				resp.Body.Close()
			}
		}

		if customID.ContentModifier != ContentModifierAudioOnly {
			content = fmt.Sprintf(
				"%s\n\n %s",
				dialogFormatted.String(),
				fmt.Sprintf(
					"`%s` @ `%s - %s` | [%s](%s) | Posted by %s",
					dialog.TranscriptMeta.Id,
					(time.Duration(dialog.Dialog[0].OffsetMs)).String(),
					(time.Duration(dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs)).String(),
					strings.TrimPrefix(b.webUrl, "https://"),
					fmt.Sprintf("%s/ep/%s#pos-%d-%d", b.webUrl, customID.EpisodeID, customID.StartLine, customID.EndLine),
					username,
				),
			)
		} else {
			content = fmt.Sprintf("Posted by %s", username)
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:     content,
			Files:       files,
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}, dialog.MaxDialogPosition, nil, cancelFunc
}

func (b *Bot) quickArchiveModalOpen(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var originalMessageID string
	if typed, ok := i.Interaction.Data.(discordgo.ApplicationCommandInteractionData); ok {
		originalMessageID = typed.TargetID
	}
	if originalMessageID == "" {
		b.respondError(s, i, fmt.Errorf("failed to find original message ID"))
		return
	}

	interactionData, ok := i.Interaction.Data.(discordgo.ApplicationCommandInteractionData)
	if !ok {
		b.respondError(s, i, fmt.Errorf("failed load target message"))
		return
	}

	modalContent := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "description",
					Label:       "Description",
					Style:       discordgo.TextInputParagraph,
					Required:    true,
					Placeholder: "Describe the content of the images",
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{discordgo.TextInput{
				CustomID:    "episode",
				Label:       "Optional Related Episode (format xfm-S01E01)",
				Style:       discordgo.TextInputShort,
				Required:    false,
				MaxLength:   128,
				Placeholder: "e.g. xfm-S01E01",
			}},
		},
	}

	userWarnings := []string{}
	if len(interactionData.Resolved.Messages[interactionData.TargetID].Attachments) == 0 {
		userWarnings = append(userWarnings, "- Message contained no media. Nothing will be submitted.")
	}
	for _, v := range interactionData.Resolved.Messages[interactionData.TargetID].Attachments {
		warning, err := b.validateAttachmentForArchive(v)
		if err != nil {
			b.logger.Error("failed to validate file", zap.Error(err))
			b.respondError(s, i, fmt.Errorf("failed to valid file"))
			return
		}
		if warning != "" {
			userWarnings = append(userWarnings, warning)
		}
	}
	if len(userWarnings) == len(interactionData.Resolved.Messages[interactionData.TargetID].Attachments) {
		userWarnings = append(userWarnings, "- All attachments were invalid, submitting this will do nothing.")
	}
	if len(userWarnings) > 0 {
		modalContent = append(modalContent, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{discordgo.TextInput{
				CustomID: "warning",
				Label:    "WARNING",
				Style:    discordgo.TextInputParagraph,
				Required: false,
				Value:    strings.Join(userWarnings, "\n"),
			}},
		})
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   "quick-archive-modal-save:" + originalMessageID,
			Title:      "Add To Archive",
			Components: modalContent,
		},
	}); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) quickArchiveModalSave(s *discordgo.Session, i *discordgo.InteractionCreate, customIDPayload string) {

	msg, err := s.ChannelMessage(i.ChannelID, customIDPayload)
	if err != nil {
		b.respondError(s, i, err)
		return
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		b.respondError(s, i, err)
		return
	}

	fileNames := []string{}
	for _, v := range msg.Attachments {
		warning, err := b.validateAttachmentForArchive(v)
		if err != nil {
			b.logger.Error("failed to validate file", zap.Error(err))
			b.followupError(s, i, fmt.Errorf("failed to valid file"))
			return
		}
		if warning != "" {
			continue
		}

		if err := b.archiveStore.ArchiveFile(v.Filename, v.URL); err != nil {
			b.followupError(s, i, err)
			return
		}

		fileNames = append(fileNames, v.Filename)
	}

	if len(fileNames) == 0 {
		b.followupError(s, i, fmt.Errorf("no valid/new files found in message"))
		return
	}

	archiveMeta := models.ArchiveMeta{
		OriginalMessageID: customIDPayload,
		CreatedAt:         time.Now(),
		Files:             fileNames,
		Description:       i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value,
	}
	if ep := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value; meta.IsValidEpisodeID(ep) {
		if publication, series, episode, err := models.ParseEpID(ep); err == nil {
			archiveMeta.Episode = models.ShortEpID(publication, series, episode)
		}
	}

	if len(archiveMeta.Files) > 0 {
		if err := b.createArchiveMeta(archiveMeta); err != nil {
			b.followupError(s, i, err)
			return
		}
	} else {
		b.followupError(s, i, fmt.Errorf("no new files were added"))
		return
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "Media was archived. Thanks!",
	}); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) confirmRewindStart(s *discordgo.Session, i *discordgo.InteractionCreate, epid string) {

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Confirm Rewind",
			Content: fmt.Sprintf("Are you sure you want to start a rewind thread?"),
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Confirm",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("rewind-start:%s", epid),
						},
					},
				},
			},
		},
	}); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) rewindStart(s *discordgo.Session, i *discordgo.InteractionCreate, epid string) {

	content, err := b.getEpisodeSummary(epid)
	if err != nil {
		b.respondError(s, i, err)
		return
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
						CustomID: fmt.Sprintf("episode-rating:%s:1", epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "2️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("episode-rating:%s:2", epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "3️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("episode-rating:%s:3", epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "4️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("episode-rating:%s:4", epid),
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							Name: "5️⃣",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: fmt.Sprintf("episode-rating:%s:5", epid),
					},
				},
			},
		},
	})
	if err != nil {
		b.respondError(s, i, err)
		return
	}

	thread, err := s.MessageThreadStartComplex(initialMessage.ChannelID, initialMessage.ID, &discordgo.ThreadStart{
		Name: fmt.Sprintf("%s REWIND", epid),
		Type: discordgo.ChannelTypeGuildPublicThread,
	})
	if err != nil {
		b.respondError(s, i, err)
		return
	}

	if err := b.createRewindState(RewindState{
		OriginalMessageID:      initialMessage.ID,
		OriginalMessageChannel: initialMessage.ChannelID,
		AnswerThreadID:         thread.ID,
		EpisodeID:              epid,
	}); err != nil {
		b.respondError(s, i, err)
		return
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Done!",
		},
	}); err != nil {
		b.respondError(s, i, err)
		return
	}
}

func (b *Bot) rateEpisode(s *discordgo.Session, i *discordgo.InteractionCreate, epidAndRating string) {

	idAndRatingParts := strings.Split(epidAndRating, ":")
	rating, err := strconv.ParseFloat(idAndRatingParts[1], 32)
	if err != nil {
		b.respondError(s, i, fmt.Errorf("failed to parse rating: %w", err))
		return
	}

	if _, err := b.transcriptApiClient.BulkSetTranscriptRatingScore(context.Background(), &api.BulkSetTranscriptRatingScoreRequest{
		Epid:        idAndRatingParts[0],
		OauthSource: "discord",
		Scores: map[string]float32{
			i.Interaction.Member.User.Username: float32(rating),
		},
	}); err != nil {
		b.respondError(s, i, err)
		return
	}

	if _, err := s.ChannelMessageSend(i.Interaction.Message.ID, fmt.Sprintf("%s rated the episode %s/5", i.Interaction.Member.DisplayName(), idAndRatingParts[1])); err != nil {
		b.respondError(s, i, err)
		return
	}

	b.respondConfirm(s, i, "Submitted!")
	return
}

func (b *Bot) respondError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	b.logger.Error("Error response was sent", zap.Error(err))
	responseErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Request failed with error: %s", err.Error()),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if responseErr != nil {
		b.logger.Error("failed to respond", zap.Error(responseErr), zap.String("original_error", err.Error()))
		return
	}
}

func (b *Bot) respondConfirm(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	responseErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if responseErr != nil {
		b.logger.Error("failed to respond", zap.Error(responseErr))
		return
	}
}

func (b *Bot) followupError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Failed: %s", err.Error()),
	}); err != nil {
		b.logger.Error("Followup error failed", zap.Error(err))
		return
	}
}

func (b *Bot) createArchiveMeta(meta models.ArchiveMeta) error {
	if err := b.archiveStore.CreateMetadata(meta); err != nil {
		if errors.Is(err, os.ErrExist) {
			// todo: could merge the image into the old meta or create a new file
			return fmt.Errorf("metadata for this message ID already exists, but some of the files do not exist. Perhaps the message was edited. Missing files: %s", strings.Join(meta.Files, ", "))
		}
		return err
	}
	return nil
}

func (b *Bot) validateAttachmentForArchive(v *discordgo.MessageAttachment) (string, error) {
	if v == nil {
		return "", nil
	}
	if !util.InStrings(v.ContentType, "image/png", "image/jpg", "image/jpeg", "image/webp") {
		return fmt.Sprintf("- SKIPPED %s was not a valid image", v.Filename), nil
	}
	exists, err := b.archiveStore.FileExists(v.Filename)
	if err != nil {
		b.logger.Error("failed to check file exists", zap.Error(err))
		return "", fmt.Errorf("failed to check file exists")
	}
	if exists {
		return fmt.Sprintf("- SKIPPED %s already exists", v.Filename), nil
	}
	return "", nil
}

func (b *Bot) getEpisodeSummary(epid string) (string, error) {
	transcript, err := b.transcriptApiClient.GetTranscript(context.Background(), &api.GetTranscriptRequest{
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

func (b *Bot) initRewindCache() error {
	entires, err := os.ReadDir("var/rewind/")
	if err != nil {
		return err
	}
	for _, v := range entires {
		if !strings.HasSuffix(v.Name(), ".json") || v.IsDir() {
			continue
		}
		threadID := strings.TrimSuffix(path.Base(v.Name()), ".json")
		b.rewindThreadCache.Store(threadID, struct{}{})
	}
	return nil
}

func (b *Bot) createRewindState(state RewindState) error {
	_, err := os.Stat(fmt.Sprintf("var/rewind/%s.json", state.AnswerThreadID))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return nil
	}
	f, err := os.Create(fmt.Sprintf("var/rewind/%s.json", state.AnswerThreadID))
	if err != nil {
		return err
	}

	if err := json.NewEncoder(f).Encode(state); err != nil {
		return err
	}

	b.rewindThreadCache.Store(state.AnswerThreadID, struct{}{})

	return nil

}

func (b *Bot) openRewindStateForReading(channelID string, cb func(cw *RewindState) error) error {
	b.rewindStateLock.RLock()
	defer b.rewindStateLock.RUnlock()

	f, err := os.Open(fmt.Sprintf("var/rewind/%s.json", channelID))
	if err != nil {
		return err
	}
	defer f.Close()

	s := RewindState{}
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return err
	}

	return cb(&s)
}

func (b *Bot) openRewindStateForWriting(channelID string, cb func(cw *RewindState) (*RewindState, error)) error {
	b.rewindStateLock.Lock()
	defer b.rewindStateLock.Unlock()

	f, err := os.OpenFile(fmt.Sprintf("var/rewind/%s.json", channelID), os.O_RDWR|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	s := &RewindState{}
	if err := json.NewDecoder(f).Decode(s); err != nil {
		return err
	}
	s, err = cb(s)
	if err != nil || s == nil {
		return err
	}

	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

func encodeCustomIDForAction(action string, customID CustomID) string {
	return fmt.Sprintf("%s:%s", action, customID.String())
}

func decodeCustomIDPayload(data string) (CustomID, error) {
	decoded := &CustomID{}
	return *decoded, json.Unmarshal([]byte(data), decoded)
}

func createFileName(dialog *api.TranscriptDialog, suffix string) string {
	if contentFilename := contentToFilename(dialog.Dialog[0].Content); contentFilename != "" {
		return fmt.Sprintf("%s.%s", contentFilename, suffix)
	}
	return fmt.Sprintf("%s-%d.%s", dialog.TranscriptMeta.Id, dialog.Dialog[0].Pos, suffix)
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

type RewindState struct {
	OriginalMessageID      string
	OriginalMessageChannel string
	AnswerThreadID         string
	EpisodeID              string
}
