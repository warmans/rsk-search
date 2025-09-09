package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/pkg/archive"
	"github.com/warmans/rsk-search/pkg/discord"
	"github.com/warmans/rsk-search/pkg/meta"
	"github.com/warmans/rsk-search/pkg/models"
	"github.com/warmans/rsk-search/pkg/util"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"
)

func NewArchiveCommand(logger *zap.Logger, archiveStore *archive.Store) *ArchiveCommand {
	return &ArchiveCommand{
		logger:       logger,
		archiveStore: archiveStore,
	}
}

type ArchiveCommand struct {
	logger       *zap.Logger
	archiveStore *archive.Store
}

func (a *ArchiveCommand) Name() string {
	return "scrimp-archive"
}

func (a *ArchiveCommand) ButtonHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{}
}

func (a *ArchiveCommand) ModalHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"quick-archive-modal-save": a.quickArchiveModalSave,
	}
}

func (a *ArchiveCommand) CommandHandlers() discord.InteractionHandlers {
	return discord.InteractionHandlers{
		"scrimp-archive": a.quickArchiveModalOpen,
	}
}

func (a *ArchiveCommand) AutoCompleteHandler() discord.InteractionHandler {
	return nil
}

func (a *ArchiveCommand) MessageHandlers() discord.MessageHandlers {
	return discord.MessageHandlers{}
}

func (a *ArchiveCommand) quickArchiveModalOpen(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {
	var originalMessageID string
	if typed, ok := i.Data.(discordgo.ApplicationCommandInteractionData); ok {
		originalMessageID = typed.TargetID
	}
	if originalMessageID == "" {
		return fmt.Errorf("failed to find original message ID")
	}

	interactionData, ok := i.Data.(discordgo.ApplicationCommandInteractionData)
	if !ok {
		return fmt.Errorf("failed load target message")
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
		warning, err := a.validateAttachmentForArchive(v)
		if err != nil {
			return fmt.Errorf("failed to valid file")
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
			CustomID:   fmt.Sprintf("%s:quick-archive-modal-save:%s", a.Name(), originalMessageID),
			Title:      "Add To Archive",
			Components: modalContent,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (a *ArchiveCommand) quickArchiveModalSave(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error {

	msg, err := s.ChannelMessage(i.ChannelID, args[0])
	if err != nil {
		return err
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	fileNames := []string{}
	for _, v := range msg.Attachments {
		warning, err := a.validateAttachmentForArchive(v)
		if err != nil {
			a.followupError(s, i, fmt.Errorf("failed to valid file"))
			return nil
		}
		if warning != "" {
			continue
		}

		if err := a.archiveStore.ArchiveFile(a.fileName(v), v.URL); err != nil {
			a.followupError(s, i, err)
			return nil
		}

		fileNames = append(fileNames, a.fileName(v))
	}

	if len(fileNames) == 0 {
		a.followupError(s, i, fmt.Errorf("no valid/new files found in message"))
		return nil
	}

	archiveMeta := models.ArchiveMeta{
		OriginalMessageID: args[0],
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
		if err := a.createArchiveMeta(archiveMeta); err != nil {
			a.followupError(s, i, err)
			return nil
		}
	} else {
		a.followupError(s, i, fmt.Errorf("no new files were added"))
		return nil
	}

	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "Media was archived. Thanks!",
	}); err != nil {
		return err
	}

	return nil
}

func (a *ArchiveCommand) followupError(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("Failed: %s", err.Error()),
	}); err != nil {
		a.logger.Error("Followup error failed", zap.Error(err))
		return
	}
}

func (a *ArchiveCommand) createArchiveMeta(meta models.ArchiveMeta) error {
	if err := a.archiveStore.CreateMetadata(meta); err != nil {
		if errors.Is(err, os.ErrExist) {
			// todo: could merge the image into the old meta or create a new file
			return fmt.Errorf("metadata for this message ID already exists, but some of the files do not exist. Perhaps the message was edited. Missing files: %s", strings.Join(meta.Files, ", "))
		}
		return err
	}
	return nil
}

func (a *ArchiveCommand) validateAttachmentForArchive(v *discordgo.MessageAttachment) (string, error) {
	if v == nil {
		return "", nil
	}
	if !util.InStrings(v.ContentType, "image/png", "image/jpg", "image/jpeg", "image/webp") {
		return fmt.Sprintf("- SKIPPED %s was not a valid image", a.fileName(v)), nil
	}
	exists, err := a.archiveStore.FileExists(a.fileName(v))
	if err != nil {
		a.logger.Error("failed to check file exists", zap.Error(err))
		return "", fmt.Errorf("failed to check file exists")
	}
	if exists {
		return fmt.Sprintf("- SKIPPED %s already exists", a.fileName(v)), nil
	}
	return "", nil
}

func (a *ArchiveCommand) fileName(v *discordgo.MessageAttachment) string {
	return fmt.Sprintf("%s-%s", v.ID, v.Filename)
}
