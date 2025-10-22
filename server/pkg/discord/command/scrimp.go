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
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type stateOpt func(c *ScrimpState)

func withModifier(mod ContentModifier) stateOpt {
	return func(c *ScrimpState) {
		c.ContentModifier = mod
	}
}

func withStartLine(pos int32) stateOpt {
	return func(c *ScrimpState) {
		c.StartLine = pos
	}
}
func withEndLine(pos int32) stateOpt {
	return func(c *ScrimpState) {
		c.EndLine = pos
	}
}

func withAudioShift(duration time.Duration) stateOpt {
	return func(c *ScrimpState) {
		c.AudioShift = duration
	}
}

func withAudioExtendOrTrim(duration time.Duration) stateOpt {
	return func(c *ScrimpState) {
		c.AudioExtendOrTrim = duration
	}
}

type ScrimpState struct {
	EpisodeID         string          `json:"e,omitempty"`
	StartLine         int32           `json:"s,omitempty"`
	EndLine           int32           `json:"f,omitempty"`
	NumContextLines   int             `json:"c,omitempty"`
	AudioShift        time.Duration   `json:"as,omitempty"`
	AudioExtendOrTrim time.Duration   `json:"ae,omitempty"`
	ContentModifier   ContentModifier `json:"t,omitempty"`
}

func (c ScrimpState) String() string {
	return mustEncodeState(c)
}

func (c ScrimpState) withOption(options ...stateOpt) ScrimpState {
	clone := &ScrimpState{
		EpisodeID:         c.EpisodeID,
		StartLine:         c.StartLine,
		EndLine:           c.EndLine,
		NumContextLines:   c.NumContextLines,
		ContentModifier:   c.ContentModifier,
		AudioShift:        c.AudioShift,
		AudioExtendOrTrim: c.AudioExtendOrTrim,
	}
	for _, v := range options {
		v(clone)
	}
	return *clone
}

type Action string

const (
	ActionNone               Action = ""
	ActionUpdateState        Action = "sta"
	ActionOpenAudioEditModal Action = "oaem"
	ActionCompleteQuery      Action = "cfm"
)

type ContentModifier uint8

const (
	ContentModifierNone ContentModifier = iota
	ContentModifierTextOnly
	ContentModifierAudioOnly
	ContentModifierGifOnly
)

type StateUpdateType string

const StateUpdateShiftDialogBackwards StateUpdateType = "sdb"
const StateUpdateShiftDialogForwards StateUpdateType = "sdf"
const StateUpdateAddPreviousLine StateUpdateType = "apl"
const StateUpdateAddNextLine StateUpdateType = "anl"
const StateUpdateTrimFirstLine StateUpdateType = "tfl"
const StateUpdateTrimLastLine StateUpdateType = "tll"
const StateUpdateAudioShift StateUpdateType = "as"
const StateUpdateAudioExtendTrim StateUpdateType = "aet"
const StateUpdateSetContentModifier StateUpdateType = "scm"

type StateUpdate struct {
	Type  StateUpdateType `json:"t"`
	Value any             `json:"v"`
}

func (s StateUpdate) CustomID() string {
	enc, err := json.Marshal(s)
	if err != nil {
		panic(fmt.Sprintf("failed to encode state update: %s", err.Error()))
	}
	return fmt.Sprintf("%s:%s", ActionUpdateState, string(enc))
}

type searchResult struct {
	EpisodeID string `json:"e,omitempty"`
	StartLine int32  `json:"s,omitempty"`
	EndLine   int32  `json:"f,omitempty"`
}

func (s searchResult) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s searchResult) ToState() ScrimpState {
	return ScrimpState{
		EpisodeID:       s.EpisodeID,
		StartLine:       s.StartLine,
		EndLine:         s.EndLine,
		ContentModifier: ContentModifierTextOnly,
	}
}

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
		string(ActionCompleteQuery):      b.queryComplete,
		string(ActionOpenAudioEditModal): b.handleOpenAudioEditModal,
		string(ActionUpdateState):        b.handleStateUpdate,
	}
}

func (b *SearchCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"se": b.handleAudioEdit,
	}
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

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, v := range res.Predictions {
		if v.Actor == "" {
			continue
		}
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name: util.TrimToN(fmt.Sprintf("[%s] %s: %s", v.Epid, v.Actor, v.Line), 100),
			Value: (&searchResult{
				EpisodeID: v.Epid,
				StartLine: v.Pos,
				EndLine:   v.Pos,
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

	result := &searchResult{}
	if err := json.Unmarshal([]byte(selection), result); err != nil {
		return fmt.Errorf("failed to decode selected result: %w", err)
	}

	return b.beginAudioResponse(s, i, result.ToState())
}

func (b *SearchCommand) queryComplete(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	state, err := extractStateFromBody[ScrimpState](i.Message)
	if err != nil {
		return err
	}

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
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(image.Body)

		files = append(files, &discordgo.File{
			Name:        attachment.Filename,
			Reader:      image.Body,
			ContentType: attachment.ContentType,
		})
	}

	content := i.Message.Content
	if state.ContentModifier == ContentModifierAudioOnly {
		username := "unknown"
		if i.Member != nil {
			username = i.Member.DisplayName()
		}
		content = fmt.Sprintf("Posted by %s", username)
	}

	interactionResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:     deleteStateFromContent(content),
			Files:       files,
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}

	if err := s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		return err
	}

	return nil
}

func (b *SearchCommand) updatePreview(s *discordgo.Session, i *discordgo.InteractionCreate, state ScrimpState) error {
	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:     "Loading...",
			Components:  make([]discordgo.MessageComponent, 0),
			Files:       make([]*discordgo.File, 0),
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}); err != nil {
		return err
	}

	interactionResponse, maxDialogOffset, cleanup, err := b.audioFileResponse(state, username, true)
	if err != nil {
		return err
	}
	defer cleanup()

	interactionResponse.Data.Components = b.buttons(state, maxDialogOffset)
	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &interactionResponse.Data.Content,
		Files:      interactionResponse.Data.Files,
		Components: &interactionResponse.Data.Components,
	})

	if err != nil {
		return fmt.Errorf("failed to update preview: %w", err)
	}
	return nil
}

func (b *SearchCommand) beginAudioResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	state ScrimpState,
) error {
	username := "unknown"
	if i.Member != nil {
		username = i.Member.DisplayName()
	}

	interactionResponse, maxDialogOffset, cleanup, err := b.audioFileResponse(state, username, true)
	if err != nil {
		common.RespondError(b.logger, s, i, err)
		return err
	}
	defer cleanup()

	interactionResponse.Data.Flags = discordgo.MessageFlagsEphemeral
	interactionResponse.Data.Components = b.buttons(state, maxDialogOffset)
	err = s.InteractionRespond(i.Interaction, interactionResponse)
	if err != nil {
		b.logger.Error("failed to respond", zap.Error(err))
	}
	return nil
}

func (b *SearchCommand) buttons(state ScrimpState, maxDialogOffset int32) []discordgo.MessageComponent {

	editRow1 := []discordgo.MessageComponent{}
	if state.StartLine > 0 {
		editRow1 = append(editRow1, discordgo.Button{
			Label: "Shift Dialog Backwards",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏪",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateShiftDialogBackwards},
			),
		})
	}
	if state.StartLine+1 < maxDialogOffset {
		editRow1 = append(editRow1, discordgo.Button{
			Label: "Shift Dialog Forward",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏩",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateShiftDialogForwards},
			),
		})
	}
	if state.EndLine-state.StartLine < 25 && state.ContentModifier != ContentModifierGifOnly {
		if state.StartLine > 0 {
			editRow1 = append(editRow1, discordgo.Button{
				// Label is what the user will see on the button.
				Label: "Add Previous Line",
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				// Style provides coloring of the button. There are not so many styles tho.
				Style: discordgo.SecondaryButton,
				// CustomID is a thing telling Discord which data to send when this button will be pressed.
				CustomID: b.encodeCustomIDForStateUpdate(
					StateUpdate{Type: StateUpdateAddPreviousLine},
				),
			})
		}
		if state.EndLine+1 < maxDialogOffset {
			editRow1 = append(editRow1, discordgo.Button{
				// Label is what the user will see on the button.
				Label: "Add Next Line",
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				// Style provides coloring of the button. There are not so many styles tho.
				Style: discordgo.SecondaryButton,
				// CustomID is a thing telling Discord which data to send when this button will be pressed.
				CustomID: b.encodeCustomIDForStateUpdate(
					StateUpdate{Type: StateUpdateAddNextLine},
				),
			})
		}
	}

	editRow2 := []discordgo.MessageComponent{}
	if state.EndLine-state.StartLine > 0 {
		editRow2 = append(editRow2, discordgo.Button{
			Label: "Trim First Line",
			Emoji: &discordgo.ComponentEmoji{
				Name: "✂",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateTrimFirstLine},
			),
		})
		editRow2 = append(editRow2, discordgo.Button{
			Label: "Trim Last Line",
			Emoji: &discordgo.ComponentEmoji{
				Name: "✂",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateTrimLastLine},
			),
		})
	}

	// audio is only enabled with these modifiers
	editRow3 := []discordgo.MessageComponent{}
	if state.ContentModifier == ContentModifierAudioOnly || state.ContentModifier == ContentModifierNone {

		editRow3 = append(editRow3, discordgo.Button{
			Label: "Shift Audio Backwards 1s",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏪",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateShiftDialogBackwards, Value: "1s"},
			),
		})
		editRow3 = append(editRow3, discordgo.Button{
			Label: "Shift Audio Forward 1s",
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏩",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateShiftDialogForwards, Value: "1s"},
			),
		})

		editRow3 = append(editRow3, discordgo.Button{
			Label: "Trim Audio 1s",
			Emoji: &discordgo.ComponentEmoji{
				Name: "✂",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateAudioExtendTrim, Value: "-1s"},
			),
		})

		editRow3 = append(editRow3, discordgo.Button{
			Label: "Extend Audio 1s",
			Emoji: &discordgo.ComponentEmoji{
				Name: "➕",
			},
			Style: discordgo.SecondaryButton,
			CustomID: b.encodeCustomIDForStateUpdate(
				StateUpdate{Type: StateUpdateAudioExtendTrim, Value: "1s"},
			),
		})

		editRow3 = append(editRow3, discordgo.Button{
			Label:    "Custom",
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("%s:oaem:%s", b.Name(), state),
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
	if len(editRow3) > 0 {
		buttons = append(buttons, discordgo.ActionsRow{
			Components: editRow3,
		})
	}
	postButtons := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Post",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("%s:%s", b.Name(), ActionCompleteQuery),
			},
		},
	}
	if state.ContentModifier != ContentModifierGifOnly {
		if state.ContentModifier != ContentModifierNone {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Audio & Text",
				Style: discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForStateUpdate(
					StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierNone)},
				),
			})
		}
		if state.ContentModifier != ContentModifierAudioOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Audio Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "🔊",
				},
				Style: discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForStateUpdate(
					StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierAudioOnly)},
				),
			})
		}
		if state.ContentModifier != ContentModifierTextOnly {
			postButtons.Components = append(postButtons.Components, discordgo.Button{
				Label: "Text Only",
				Emoji: &discordgo.ComponentEmoji{
					Name: "📄",
				},
				Style: discordgo.SecondaryButton,
				CustomID: b.encodeCustomIDForStateUpdate(
					StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierTextOnly)},
				),
			})
		}
	}
	if state.StartLine == state.EndLine && state.NumContextLines == 0 {
		if state.ContentModifier != ContentModifierGifOnly {
			postButtons.Components = append(postButtons.Components,
				discordgo.Button{
					Label: "GIF mode",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📺",
					},
					Style: discordgo.SecondaryButton,
					CustomID: b.encodeCustomIDForStateUpdate(
						StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierGifOnly)},
					),
				},
			)
		} else {
			postButtons.Components = append(postButtons.Components,
				discordgo.Button{
					Label: "Normal mode",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📻",
					},
					Style: discordgo.SecondaryButton,
					CustomID: b.encodeCustomIDForStateUpdate(
						StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierTextOnly)},
					),
				},
				discordgo.Button{
					Label: "Randomize image",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📺",
					},
					Style: discordgo.SecondaryButton,
					CustomID: b.encodeCustomIDForStateUpdate(
						StateUpdate{Type: StateUpdateSetContentModifier, Value: fmt.Sprint(ContentModifierGifOnly)},
					),
				},
			)
		}
	}

	buttons = append(buttons, postButtons)

	return buttons
}

func (b *SearchCommand) audioFileResponse(
	state ScrimpState,
	username string,
	isPreview bool,
) (*discordgo.InteractionResponse, int32, func(), error) {

	dialog, err := b.transcriptApiClient.GetTranscriptDialog(context.Background(), &api.GetTranscriptDialogRequest{
		Epid: state.EpisodeID,
		Range: &api.DialogRange{
			Start: state.StartLine,
			End:   state.EndLine,
		},
	})
	if err != nil {
		return nil, 0, func() {}, fmt.Errorf("failed to fetch selected line")
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

	if state.ContentModifier == ContentModifierGifOnly {
		audioFileURL := fmt.Sprintf(
			"%s/dl/media/%s.gif?ts=%d-%d",
			b.webUrl,
			dialog.TranscriptMeta.ShortId,
			dialog.Dialog[0].OffsetMs,
			dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs,
		)
		resp, err := http.Get(audioFileURL)
		if err != nil {
			return nil, 0, func() {}, fmt.Errorf("failed to fetch selected line")
		}
		if resp.StatusCode != http.StatusOK {
			b.logger.Error("failed to fetch gif", zap.Error(err), zap.String("url", audioFileURL), zap.Int("status_code", resp.StatusCode))
			return nil, 0, func() {}, fmt.Errorf("failed to fetch gif: %s", resp.Status)
		}
		files = append(files, &discordgo.File{
			Name:        createFileName(dialog, "gif"),
			ContentType: "image/gif",
			Reader:      resp.Body,
		})
		cancelFunc = func() {
			_ = resp.Body.Close()
		}
		content = fmt.Sprintf(
			"`%s` @ `%s - %s` | [%s](%s) | Posted by %s",
			dialog.TranscriptMeta.Id,
			(time.Duration(dialog.Dialog[0].OffsetMs)).String(),
			(time.Duration(dialog.Dialog[len(dialog.Dialog)-1].OffsetMs + dialog.Dialog[len(dialog.Dialog)-1].DurationMs)).String(),
			strings.TrimPrefix(b.webUrl, "https://"),
			fmt.Sprintf("%s/ep/%s#pos-%d-%d", b.webUrl, state.EpisodeID, state.StartLine, state.EndLine),
			username,
		)

	} else {
		if state.ContentModifier != ContentModifierTextOnly {
			audioFileURL := fmt.Sprintf(
				"%s/dl/media/%s.mp3?ts=%d-%d&extend=%s&shift=%s",
				b.webUrl,
				dialog.TranscriptMeta.ShortId,
				dialog.Dialog[0].OffsetMs,
				dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs,
				state.AudioExtendOrTrim.String(),
				state.AudioShift.String(),
			)
			resp, err := http.Get(audioFileURL)
			if err != nil {
				return nil, 0, func() {}, fmt.Errorf("failed to fetch selected line")
			}
			if resp.StatusCode != http.StatusOK {
				b.logger.Error("failed to fetch audio", zap.Error(err), zap.String("url", audioFileURL), zap.Int("status_code", resp.StatusCode))
				return nil, 0, func() {}, fmt.Errorf("failed to fetch audio: %s", resp.Status)
			}
			files = append(files, &discordgo.File{
				Name:        createFileName(dialog, "mp3"),
				ContentType: "audio/mpeg",
				Reader:      resp.Body,
			})
			cancelFunc = func() {
				_ = resp.Body.Close()
			}
		}

		audioModifierDescriptions := []string{}
		if state.AudioShift != 0 {
			audioModifierDescriptions = append(audioModifierDescriptions, fmt.Sprintf("[>> %s]", state.AudioShift))
		}
		if state.AudioExtendOrTrim != 0 {
			audioModifierDescriptions = append(audioModifierDescriptions, fmt.Sprintf("[++ %s]", state.AudioExtendOrTrim))
		}

		dialogText := fmt.Sprintf("%s\n\n", dialogFormatted.String())
		if state.ContentModifier == ContentModifierAudioOnly {
			dialogText = ""
		}
		content = fmt.Sprintf(
			"%s%s",
			dialogText,
			fmt.Sprintf(
				"`%s` @ `%s - %s %s` | [%s](%s) | Posted by %s",
				dialog.TranscriptMeta.Id,
				(time.Duration(dialog.Dialog[0].OffsetMs)).String(),
				(time.Duration(dialog.Dialog[len(dialog.Dialog)-1].OffsetMs+dialog.Dialog[len(dialog.Dialog)-1].DurationMs)).String(),
				strings.Join(audioModifierDescriptions, " "),
				strings.TrimPrefix(b.webUrl, "https://"),
				fmt.Sprintf("%s/ep/%s#pos-%d-%d", b.webUrl, state.EpisodeID, state.StartLine, state.EndLine),
				username,
			),
		)
	}

	if isPreview {
		content = fmt.Sprintf("%s\n%s", content, state.String())
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:     content,
			Files:       files,
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}, dialog.MaxDialogPosition, cancelFunc, nil
}

func (b *SearchCommand) encodeCustomIDForStateUpdate(opt StateUpdate) string {
	return fmt.Sprintf("%s:%s", b.Name(), opt.CustomID())
}

func (b *SearchCommand) handleOpenAudioEditModal(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	state, err := extractStateFromBody[ScrimpState](i.Message)
	if err != nil {
		return err
	}

	fields := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "shift",
					Label:       "Shift Audio",
					Placeholder: "-1s",
					Style:       discordgo.TextInputShort,
					Required:    true,
					MaxLength:   10,
					Value:       state.AudioShift.String(),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "extend",
					Label:       "Extend/Trim Audio",
					Placeholder: "2s",
					Style:       discordgo.TextInputShort,
					Required:    true,
					MaxLength:   10,
					Value:       state.AudioExtendOrTrim.String(),
				},
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("%s:se", b.Name()),
			Title:      "Edit Audio",
			Components: fields,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *SearchCommand) handleStateUpdate(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[ScrimpState](i.Message)
	if err != nil {
		return fmt.Errorf("Failed decode state: %w", err)
	}

	update := &StateUpdate{}
	if err := json.Unmarshal(
		[]byte(strings.TrimPrefix(args[0], fmt.Sprintf("%s:", ActionUpdateState))),
		update,
	); err != nil {
		return err
	}

	switch update.Type {
	case StateUpdateShiftDialogBackwards:
		return b.updatePreview(
			s,
			i,
			state.withOption(withEndLine(state.EndLine-1), withStartLine(state.StartLine-1)),
		)
	case StateUpdateShiftDialogForwards:
		return b.updatePreview(
			s,
			i,
			state.withOption(withEndLine(state.EndLine+1), withStartLine(state.StartLine+1)),
		)
	case StateUpdateAddPreviousLine:
		return b.updatePreview(
			s,
			i,
			state.withOption(
				withStartLine(state.StartLine-1),
			),
		)
	case StateUpdateAddNextLine:
		return b.updatePreview(
			s,
			i,
			state.withOption(
				withEndLine(state.EndLine+1),
			),
		)
	case StateUpdateTrimFirstLine:
		return b.updatePreview(
			s,
			i,
			state.withOption(withStartLine(state.StartLine-1)),
		)
	case StateUpdateTrimLastLine:
		return b.updatePreview(
			s,
			i,
			state.withOption(withEndLine(state.EndLine-1)),
		)
	case StateUpdateAudioShift:
		strVal, ok := update.Value.(string)
		if !ok {
			return fmt.Errorf("update had wrong type %T", update.Value)
		}
		duration, err := time.ParseDuration(strVal)
		if err != nil {
			return err
		}
		return b.updatePreview(
			s,
			i,
			state.withOption(withAudioShift(duration)),
		)
	case StateUpdateAudioExtendTrim:
		strVal, ok := update.Value.(string)
		if !ok {
			return fmt.Errorf("update had wrong type %T", update.Value)
		}
		duration, err := time.ParseDuration(strVal)
		if err != nil {
			return err
		}
		return b.updatePreview(
			s,
			i,
			state.withOption(withAudioExtendOrTrim(duration)),
		)
	case StateUpdateSetContentModifier:
		// encode as a string to stop it being interpreted as a float
		strVal, ok := update.Value.(string)
		if !ok {
			return fmt.Errorf("update had wrong type %T", update.Value)
		}
		modifier, err := strconv.ParseInt(strVal, 10, 8)
		if err != nil {
			return err
		}
		return b.updatePreview(
			s,
			i,
			state.withOption(withModifier(ContentModifier(uint8(modifier)))),
		)
	}

	return b.updatePreview(s, i, *state)
}

func (b *SearchCommand) handleAudioEdit(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	shiftStrVal := i.Interaction.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	extendStrVal := i.Interaction.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	var shift time.Duration
	var extend time.Duration
	var err error
	if shiftStrVal != "" {
		shift, err = time.ParseDuration(shiftStrVal)
		if err != nil {
			return fmt.Errorf("invalid shift value: %s", shiftStrVal)
		}
	}
	if extendStrVal != "" {
		extend, err = time.ParseDuration(extendStrVal)
		if err != nil {
			return fmt.Errorf("invalid extend value: %s", extendStrVal)
		}
	}
	state, err := extractStateFromBody[ScrimpState](i.Message)
	if err != nil {
		return err
	}

	if err := b.updatePreview(
		s,
		i,
		state.withOption(
			withAudioShift(shift),
			withAudioExtendOrTrim(extend),
		),
	); err != nil {
		return fmt.Errorf("failed to update preview: %w", err)
	}

	return nil
}

func createFileName(dialog *api.TranscriptDialog, suffix string) string {
	if contentFilename := contentToFilename(dialog.Dialog[0].Content); contentFilename != "" {
		return fmt.Sprintf("%s.%s", contentFilename, suffix)
	}
	return fmt.Sprintf("%s-%d.%s", dialog.TranscriptMeta.Id, dialog.Dialog[0].Pos, suffix)
}
