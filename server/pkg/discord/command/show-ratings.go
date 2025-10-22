package command

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/chart"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/filter"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type State struct {
	Mine        bool       `json:"m"`
	Sort        bool       `json:"s"`
	Kind        chart.Kind `json:"k"`
	Publication string     `json:"p"`
	Series      string     `json:"sr"`
}

func NewShowRatingsCommand(
	logger *zap.Logger,
	transcriptApiClient api.TranscriptServiceClient) *ShowRatingsCommand {
	return &ShowRatingsCommand{
		logger:              logger,
		transcriptApiClient: transcriptApiClient,
	}
}

type ShowRatingsCommand struct {
	logger              *zap.Logger
	transcriptApiClient api.TranscriptServiceClient
}

func (r *ShowRatingsCommand) Kind() discordgo.ApplicationCommandOptionType {
	return discordgo.ApplicationCommandOptionString
}

func (r *ShowRatingsCommand) Name() string {
	return "scrimp-show-ratings"
}

func (r *ShowRatingsCommand) Description() string {
	return "Show ratings scatter plot"
}

func (r *ShowRatingsCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (r *ShowRatingsCommand) ButtonHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"post":               r.handlePost,
		"load":               r.handleLoadChart,
		"publication-filter": r.handlePublicationFilter,
		"series-filter":      r.handleSeriesFilter,
		"toggle-mine":        r.handleToggleMine,
		"toggle-sorting":     r.handleToggleSorting,
		"set-kind":           r.handleSetKind,
	}
}

func (r *ShowRatingsCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{}
}

func (r *ShowRatingsCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		r.Name(): r.handleInitialInvocation,
	}
}

func (r *ShowRatingsCommand) AutoCompleteHandler() discord.InteractionHandler {
	return nil
}

func (r *ShowRatingsCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{}
}

func (r *ShowRatingsCommand) handleToggleSorting(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}
	state.Sort = !state.Sort

	return r._handleLoadChart(s, i, state, "")
}

func (r *ShowRatingsCommand) handleToggleMine(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}
	state.Mine = !state.Mine
	// force the chart to avg if it's a user's ratings
	if state.Mine {
		state.Kind = chart.RatingAvg
	}

	return r._handleLoadChart(s, i, state, "")
}

func (r *ShowRatingsCommand) handleSetKind(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}

	switch args[0] {
	case "avg":
		state.Kind = chart.RatingAvg
		state.Mine = false
	case "count":
		state.Kind = chart.RatingCounts
		state.Mine = false
	case "breakdown":
		state.Kind = chart.RatingBreakdown
		state.Mine = false
	}

	return r._handleLoadChart(s, i, state, "")
}

func (r *ShowRatingsCommand) handlePublicationFilter(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}
	switch args[0] {
	case "all":
		state.Publication = ""
	case "xfm":
		state.Publication = "xfm"
	case "guide":
		state.Publication = "guide"
	case "podcast":
		state.Publication = "podcast"
	case "special":
		state.Publication = "special"
	}
	state.Series = ""

	return r._handleLoadChart(s, i, state, "")
}

func (r *ShowRatingsCommand) handleSeriesFilter(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}

	state.Series = args[0]

	return r._handleLoadChart(s, i, state, "")
}

func (r *ShowRatingsCommand) handleLoadChart(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}
	return r._handleLoadChart(s, i, state, args...)
}

func (r *ShowRatingsCommand) _handleLoadChart(s *discordgo.Session, i *discordgo.InteractionCreate, state *State, args ...string) error {

	var buff *bytes.Buffer
	var err error
	var author *string

	if state.Mine {
		author = util.ToPtr(i.Member.User.Username)
	}

	var f *filter.Filter
	if state.Publication != "" {
		if state.Publication == "special" {
			f = util.ToPtr(filter.Eq("special", filter.Bool(true)))
		} else {
			f = util.ToPtr(filter.Eq("publication", filter.String(state.Publication)))
			if state.Series != "" {
				seriesInt, err := strconv.ParseInt(state.Series, 10, 32)
				if err != nil {
					return err
				}
				f = util.ToPtr(filter.And(*f, filter.Eq("series", filter.Int(seriesInt))))
			}
		}

	}

	buff, err = r.ratingsChart(f, author, state.Sort, state.Kind)
	if err != nil {
		return err
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:       discordgo.MessageFlagsEphemeral,
			Content:     mustEncodeState[State](*state),
			Files:       []*discordgo.File{{Name: "ratings.png", ContentType: "image/png", Reader: buff}},
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
			Components:  r.buttons(*state),
		},
	})
}

func (r *ShowRatingsCommand) handlePost(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	if i.Type != discordgo.InteractionMessageComponent {
		return nil
	}

	state, err := extractStateFromBody[State](i.Message)
	if err != nil {
		return err
	}

	bodyContent := ""
	if state.Mine {
		bodyContent = fmt.Sprintf("Ratings by %s", i.Member.DisplayName())
	} else {
		if state.Kind == chart.RatingCounts {
			bodyContent = "Episode Ratings Count"
		}
		if state.Kind == chart.RatingAvg {
			bodyContent = "Episode Average Ratings"
		}
		if state.Kind == chart.RatingBreakdown {
			bodyContent = "Episode Ratings Breakdown"
		}
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

	interactionResponse := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:     bodyContent,
			Files:       files,
			Attachments: util.ToPtr([]*discordgo.MessageAttachment{}),
		},
	}

	if err := s.InteractionRespond(i.Interaction, interactionResponse); err != nil {
		return err
	}

	return nil
}

func (r *ShowRatingsCommand) handleInitialInvocation(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	defaultState := State{
		Mine:        false,
		Sort:        false,
		Kind:        chart.RatingAvg,
		Publication: "",
		Series:      "",
	}

	buff, err := r.ratingsChart(nil, nil, false, chart.RatingAvg)
	if err != nil {
		return err
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:      "Chart",
			Content:    mustEncodeState[State](defaultState),
			Flags:      discordgo.MessageFlagsEphemeral,
			Files:      []*discordgo.File{{Name: "ratings.png", ContentType: "image/png", Reader: buff}},
			Components: r.buttons(defaultState),
		},
	}); err != nil {
		return err
	}

	return nil
}

func (r *ShowRatingsCommand) buttons(state State) []discordgo.MessageComponent {
	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Sort",
					Emoji: &discordgo.ComponentEmoji{
						Name: "↘️",
					},
					Style:    buttonStyleIf(state.Sort, discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:toggle-sorting", r.Name()),
				},
				discordgo.Button{
					Label: "My Ratings",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📊",
					},
					Style:    buttonStyleIf(state.Mine, discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:toggle-mine", r.Name()),
				},
				discordgo.Button{
					Label: "Average",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📊",
					},
					Style:    buttonStyleIf(state.Kind == chart.RatingAvg && !state.Mine, discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:set-kind:avg", r.Name()),
				},
				discordgo.Button{
					Label: "Count",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📊",
					},
					Style:    buttonStyleIf(state.Kind == chart.RatingCounts, discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:set-kind:count", r.Name()),
				},
				discordgo.Button{
					Label: "Breakdown",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📊",
					},
					Style:    buttonStyleIf(state.Kind == chart.RatingBreakdown, discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:set-kind:breakdown", r.Name()),
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "All",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📂",
					},
					Style:    buttonStyleIf(state.Publication == "", discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:publication-filter:all", r.Name()),
				},
				discordgo.Button{
					Label: "Xfm",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📂",
					},
					Style:    buttonStyleIf(state.Publication == "xfm", discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:publication-filter:xfm", r.Name()),
				},
				discordgo.Button{
					Label: "Podcast",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📂",
					},
					Style:    buttonStyleIf(state.Publication == "podcast", discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:publication-filter:podcast", r.Name()),
				},
				discordgo.Button{
					Label: "Guide",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📂",
					},
					Style:    buttonStyleIf(state.Publication == "guide", discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:publication-filter:guide", r.Name()),
				},
				discordgo.Button{
					Label: "Special",
					Emoji: &discordgo.ComponentEmoji{
						Name: "📂",
					},
					Style:    buttonStyleIf(state.Publication == "special", discordgo.SuccessButton, discordgo.SecondaryButton),
					CustomID: fmt.Sprintf("%s:publication-filter:special", r.Name()),
				},
			},
		},
	}

	switch state.Publication {
	case "xfm":
		row := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{},
		}
		for i := 0; i <= 4; i++ {
			row.Components = append(row.Components, discordgo.Button{
				Label:    fmt.Sprintf("%d", i),
				Style:    buttonStyleIf(state.Series == fmt.Sprintf("%d", i), discordgo.SuccessButton, discordgo.SecondaryButton),
				CustomID: fmt.Sprintf("%s:series-filter:%d", r.Name(), i),
			})
		}
		buttons = append(buttons, row)
	case "podcast":
		row := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{},
		}
		for i := 1; i <= 4; i++ {
			row.Components = append(row.Components, discordgo.Button{
				Label:    fmt.Sprintf("%d", i),
				Style:    buttonStyleIf(state.Series == fmt.Sprintf("%d", i), discordgo.SuccessButton, discordgo.SecondaryButton),
				CustomID: fmt.Sprintf("%s:series-filter:%d", r.Name(), i),
			})
		}
		buttons = append(buttons, row)

	}

	buttons = append(buttons, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Post",
				Style:    discordgo.PrimaryButton,
				CustomID: fmt.Sprintf("%s:post", r.Name()),
			},
		},
	})

	return buttons
}

func (r *ShowRatingsCommand) ratingsChart(f *filter.Filter, author *string, sort bool, kind chart.Kind) (*bytes.Buffer, error) {

	var canvas *gg.Context
	if kind == chart.RatingBreakdown {
		var err error
		canvas, err = chart.GenerateBreakdownChart(
			context.Background(),
			r.transcriptApiClient,
			f,
		)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		canvas, err = chart.GenerateRatingsChart(
			context.Background(),
			r.transcriptApiClient,
			f,
			author,
			sort,
			kind,
		)
		if err != nil {
			return nil, err
		}
	}

	buff := &bytes.Buffer{}
	if err := canvas.EncodePNG(buff); err != nil {
		return nil, err
	}
	return buff, nil

}

func buttonStyleIf(cond bool, style discordgo.ButtonStyle, def discordgo.ButtonStyle) discordgo.ButtonStyle {
	if cond {
		return style
	}
	return def
}
