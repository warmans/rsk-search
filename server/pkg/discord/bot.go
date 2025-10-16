package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/warmans/rsk-search/pkg/discord/common"
	"go.uber.org/zap"
	"log"
	"strings"
)

type InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) error
type InteractionHandlers map[string]InteractionHandler
type MessageHandlers []func(s *discordgo.Session, m *discordgo.MessageCreate)

type Command interface {
	Name() string
	ButtonHandlers() InteractionHandlers
	ModalHandlers() InteractionHandlers
	CommandHandlers() InteractionHandlers
	MessageHandlers() MessageHandlers
}

type SlashCommand interface {
	Command
	Kind() discordgo.ApplicationCommandOptionType
	Description() string
	Options() []*discordgo.ApplicationCommandOption
	AutoCompleteHandler() InteractionHandler
}

func NewBot(
	logger *zap.Logger,
	session *discordgo.Session,
	slashCommands []SlashCommand,
	applicationCommands []Command,
) *Bot {
	session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsAllWithoutPrivileged | discordgo.IntentMessageContent | discordgo.IntentGuildMembers,
	)

	bot := &Bot{
		logger:               logger,
		session:              session,
		commands:             []*discordgo.ApplicationCommand{},
		buttonHandlers:       InteractionHandlers{},
		modalHandlers:        InteractionHandlers{},
		autoCompleteHandlers: InteractionHandlers{},
		commandHandlers:      map[string]InteractionHandlers{},
	}

	for _, c := range slashCommands {
		bot.commands = append(bot.commands, &discordgo.ApplicationCommand{
			Name:        c.Name(),
			Type:        discordgo.ChatApplicationCommand,
			Description: c.Description(),
			Options:     c.Options(),
		})
		for k, v := range c.CommandHandlers() {
			if bot.commandHandlers[c.Name()] == nil {
				bot.commandHandlers[c.Name()] = InteractionHandlers{}
			}
			bot.commandHandlers[c.Name()][k] = v
		}
		for k, v := range c.ButtonHandlers() {
			bot.buttonHandlers[fmt.Sprintf("%s:%s", c.Name(), k)] = v
		}
		for k, v := range c.ModalHandlers() {
			bot.modalHandlers[fmt.Sprintf("%s:%s", c.Name(), k)] = v
		}
		if c.AutoCompleteHandler() != nil {
			bot.autoCompleteHandlers[c.Name()] = c.AutoCompleteHandler()
		}

		bot.messageHandlers = append(bot.messageHandlers, c.MessageHandlers()...)
	}

	for _, c := range applicationCommands {
		bot.commands = append(
			bot.commands,
			&discordgo.ApplicationCommand{
				Name: c.Name(),
				Type: discordgo.MessageApplicationCommand,
			})

		for k, v := range c.ButtonHandlers() {
			bot.buttonHandlers[fmt.Sprintf("%s:%s", c.Name(), k)] = v
		}
		for k, v := range c.ModalHandlers() {
			bot.modalHandlers[fmt.Sprintf("%s:%s", c.Name(), k)] = v
		}
		for k, v := range c.CommandHandlers() {
			if bot.commandHandlers[c.Name()] == nil {
				bot.commandHandlers[c.Name()] = InteractionHandlers{}
			}
			bot.commandHandlers[c.Name()][k] = v
		}
		bot.messageHandlers = append(bot.messageHandlers, c.MessageHandlers()...)
	}

	return bot
}

type Bot struct {
	logger               *zap.Logger
	session              *discordgo.Session
	commands             []*discordgo.ApplicationCommand
	commandHandlers      map[string]InteractionHandlers
	autoCompleteHandlers InteractionHandlers
	buttonHandlers       InteractionHandlers
	modalHandlers        InteractionHandlers
	messageHandlers      MessageHandlers

	createdCommands []*discordgo.ApplicationCommand
}

func (b *Bot) Start() error {

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if err := b.handleRootCommand(s, i); err != nil {
				common.RespondError(b.logger, s, i, err)
			}
			return
		case discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := b.autoCompleteHandlers[i.ApplicationCommandData().Name]; ok {
				if err := h(s, i); err != nil {
					common.RespondError(b.logger, s, i, err)
				}
				return
			}
			common.RespondError(b.logger, s, i, fmt.Errorf("no handler for autocomplete action: %s", i.ApplicationCommandData().Name))
			return
		case discordgo.InteractionModalSubmit:
			for k, h := range b.modalHandlers {
				if strings.HasPrefix(i.ModalSubmitData().CustomID, k) {
					arg := strings.TrimPrefix(i.ModalSubmitData().CustomID, fmt.Sprintf("%s:", k))
					if err := h(s, i, arg); err != nil {
						b.logger.Error(fmt.Sprintf("failed to execute handler for modal: %s", k))
						common.RespondError(b.logger, s, i, err)
					}
					return
				}
			}
			common.RespondError(b.logger, s, i, fmt.Errorf("no handler for modal action: %s", i.ModalSubmitData().CustomID))
			return
		case discordgo.InteractionMessageComponent:
			for k, h := range b.buttonHandlers {
				if strings.HasPrefix(i.MessageComponentData().CustomID, k) {
					arg := strings.TrimPrefix(i.MessageComponentData().CustomID, fmt.Sprintf("%s:", k))
					if err := h(s, i, arg); err != nil {
						b.logger.Error(fmt.Sprintf("failed to execute handler for button: %s", k))
						common.RespondError(b.logger, s, i, err)
					}
					return
				}
			}
			common.RespondError(b.logger, s, i, fmt.Errorf("no handler for button action: %s", i.MessageComponentData().CustomID))
			return
		}
	})
	for _, v := range b.messageHandlers {
		b.session.AddHandler(v)
	}

	if err := b.session.Open(); err != nil {
		return fmt.Errorf("failed to open session: %w", err)
	}

	var err error
	b.createdCommands, err = b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", b.commands)
	if err != nil {
		return fmt.Errorf("cannot register commands: %w", err)
	}
	return nil
}

func (b *Bot) Close() error {
	// cleanup commands
	for _, cmd := range b.createdCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, "", cmd.ID)
		if err != nil {
			return fmt.Errorf("cannot delete %s command: %w", cmd.Name, err)
		}
	}
	return b.session.Close()
}

//func (b *Bot) handleMessageReactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
//	if _, isValidThread := b.rewindThreadCache.Load(r.ChannelID); isValidThread {
//		fmt.Println("reaction to message: ", r.MessageID, " channel: ", r.ChannelID)
//	}
//}

func (b *Bot) handleRootCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	command, ok := b.commandHandlers[i.ApplicationCommandData().Name]
	if !ok {
		return fmt.Errorf("unkown sub-command: %s", i.ApplicationCommandData().Options[0].Options[0].Name)
	}
	if commandOption, ok := command[i.ApplicationCommandData().Name]; ok {
		return commandOption(s, i)
	}
	return fmt.Errorf("unknown command: %s", i.ApplicationCommandData().Name)
}
