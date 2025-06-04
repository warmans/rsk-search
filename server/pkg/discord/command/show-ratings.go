package command

import (
	"bytes"
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/gen/api"
	"github.com/warmans/rsk-search/pkg/chart"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/discord/common"
	"go.uber.org/zap"
)

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
	return discord.InteractionHandlers{}
}

func (r *ShowRatingsCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{}
}

func (r *ShowRatingsCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		r.Name(): r.handleShowRatingChart,
	}
}

func (r *ShowRatingsCommand) AutoCompleteHandler() discord.InteractionHandler {
	return nil
}

func (r *ShowRatingsCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{}
}

func (r *ShowRatingsCommand) handleShowRatingChart(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	canvas, err := chart.GenerateRatingsChart(context.Background(), r.transcriptApiClient)
	if err != nil {
		return err
	}
	buff := &bytes.Buffer{}
	if err := canvas.EncodePNG(buff); err != nil {
		return err
	}

	if _, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{{Name: "ratings.png", ContentType: "image/png", Reader: buff}},
	}); err != nil {
		return err
	}

	common.RespondConfirm(r.logger, s, i, "OK!")

	return nil
}
