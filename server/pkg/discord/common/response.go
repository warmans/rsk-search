package common

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func RespondError(logger *zap.Logger, s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	logger.Error("Error response was sent", zap.Error(err))
	responseErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Request failed with error: %s", err.Error()),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if responseErr != nil {
		logger.Error("failed to respond", zap.Error(responseErr), zap.String("original_error", err.Error()))
		return
	}
}

func RespondConfirm(logger *zap.Logger, s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	responseErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if responseErr != nil {
		logger.Error("failed to respond", zap.Error(responseErr))
		return
	}
}
