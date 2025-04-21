package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/discord/common"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/searchterms"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strings"
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

func NewSearchCommand(
	logger *zap.Logger,
	session *discordgo.Session,
	webUrl string,
	transcriptApiClient api.TranscriptServiceClient,
	searchApiClient api.SearchServiceClient) *SearchCommand {
	return &SearchCommand{
		logger:              logger,
		session:             session,
		webUrl:              webUrl,
		transcriptApiClient: transcriptApiClient,
		searchApiClient:     searchApiClient,
	}

}

type SearchCommand struct {
	logger              *zap.Logger
	session             *discordgo.Session
	webUrl              string
	transcriptApiClient api.TranscriptServiceClient
	searchApiClient     api.SearchServiceClient
}

func (b *SearchCommand) Kind() discordgo.ApplicationCommandOptionType {
	return discordgo.ApplicationCommandOptionString
}

func (b *SearchCommand) Name() string {
	return "scrimp"
}

func (b *SearchCommand) Description() string {
	return "Search Scrimpton transcripts."
}

func (b *SearchCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Name:         "query",
			Description:  "enter a partial quote",
			Type:         discordgo.ApplicationCommandOptionString,
			Required:     true,
			Autocomplete: true,
		},
	}
}

func (b *SearchCommand) ButtonHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"cfm": b.queryComplete,
		"up":  b.updatePreview,
	}
}

func (b *SearchCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{}
}

func (b *SearchCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		b.Name(): b.queryBegin,
	}
}

func (b *SearchCommand) AutoCompleteHandler() discord.InteractionHandler {
	return b.handleAutocomplete
}

func (b *SearchCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{}
}

func (b *SearchCommand) handleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	data := i.ApplicationCommandData()

	rawTerms := strings.TrimSpace(data.Options[0].StringValue())

	terms, err := searchterms.Parse(rawTerms)
	if err != nil {
		return err
	}
	if len(terms) == 0 {
		if err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: []*discordgo.ApplicationCommandOptionChoice{},
			},
		}); err != nil {
			return err
		}
		return nil
	}

	filterString, err := filter.Print(searchterms.TermsToFilter(terms))
	if err != nil {
		return fmt.Errorf("failed to create filter: %w", err)
	}
	res, err := b.searchApiClient.PredictSearchTerm(
		context.Background(),
		&api.PredictSearchTermRequest{
			Query:          filterString,
			MaxPredictions: 25,
		},
	)
	if err != nil {
		return err
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

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func (b *SearchCommand) queryBegin(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	selection := i.ApplicationCommandData().Options[0].StringValue()
	if selection == "" {
		return nil
	}
	customID, err := decodeCustomIDPayload(selection)
	if err != nil {
		return err
	}
	if err := b.beginAudioResponse(s, i, customID); err != nil {
		return err
	}

	return nil
}

func (b *SearchCommand) queryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}
	// can we get the files of the existing message?
	var files []*discordgo.File
	if len(i.Message.Attachments) > 0 {
		attachment := i.Message.Attachments[0]
		image, err := http.Get(attachment.URL)
		if err != nil {
			return fmt.Errorf("failed to get original message attachment: %w", err)
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
		return err
	}

	return nil
}

func (b *SearchCommand) updatePreview(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	customID, err := decodeCustomIDPayload(args[0])
	if err != nil {
		return err
	}
	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}

	interactionResponse, maxDialogOffset, err, cleanup := b.audioFileResponse(customID, username)
	if err != nil {
		return err
	}
	defer cleanup()

	interactionResponse.Data.Components = b.buttons(customID, maxDialogOffset)
	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: interactionResponse.Data,
	}); err != nil {
		return err
	}
	return nil
}

func (b *SearchCommand) beginAudioResponse(
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
		common.RespondError(b.logger, s, i, err)
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

func (b *SearchCommand) buttons(customID CustomID, maxDialogOffset int32) []discordgo.MessageComponent {

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
			CustomID: b.encodeCustomIDForAction(
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
			CustomID: b.encodeCustomIDForAction(
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
				CustomID: b.encodeCustomIDForAction(
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
				CustomID: b.encodeCustomIDForAction(
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
			CustomID: b.encodeCustomIDForAction(
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
			CustomID: b.encodeCustomIDForAction(
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
				CustomID: b.encodeCustomIDForAction("cfm", customID),
			},
		},
	}
	if customID.ContentModifier != ContentModifierGifOnly {
		if customID.ContentModifier != ContentModifierNone {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label:    "Audio & Text",
				Style:    discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierNone))),
			})
		}
		if customID.ContentModifier != ContentModifierAudioOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Audio Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "🔊",
				},
				Style:    discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierAudioOnly))),
			})
		}
		if customID.ContentModifier != ContentModifierTextOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Text Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "📄",
				},
				Style:    discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierTextOnly))),
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
					CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierGifOnly))),
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
					CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierTextOnly))),
				},
				discordgo.Button{
					Label: "Randomize image",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📺",
					},
					Style:    discordgo.SecondaryButton,
					CustomID: b.encodeCustomIDForAction("up", customID.withOption(withModifier(ContentModifierGifOnly))),
				},
			)
		}
	}

	buttons = append(buttons, postButtons)

	return buttons
}

func (b *SearchCommand) audioFileResponse(customID CustomID, username string) (*discordgo.InteractionResponse, int32, error, func()) {

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

func (b *SearchCommand) encodeCustomIDForAction(action string, customID CustomID) string {
	return fmt.Sprintf("%s:%s:%s", b.Name(), action, customID.String())
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
