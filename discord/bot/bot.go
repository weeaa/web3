package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/database/db"
	"os"
)

type Bot struct {
	s  *discordgo.Session
	db *db.DB
}

func New(db *db.DB) (*Bot, error) {
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, err
	}

	log.Info("Discord Bot Ready ✅")

	bot := &Bot{s, db}

	s.AddHandler(bot.onReady)
	s.AddHandler(bot.onRoleReactionAdd)
	s.AddHandler(bot.onRoleReactionRemove)
	s.AddHandler(bot.onSlashCommand)
	s.AddHandler(bot.onInteractionCreate)

	s.Identify.Intents = discordgo.IntentsGuildMessageReactions

	return bot, nil
}

func (b *Bot) Start() error {
	return b.s.Open()
}

func (b *Bot) Stop() error {
	return b.s.Close()
}

func (b *Bot) onReady(s *discordgo.Session, r *discordgo.Event) {
	if err := s.UpdateListeningStatus("rugging 🖕"); err != nil {
		log.Error(err)
		return
	}

	b.checkIfMsgSent()
	b.registerCommands()
}

func (b *Bot) checkIfMsgSent() {
	if !b.getMessages(RolesChannel, "👤 — Roles") {
		b.messageRoleChannel()
	}
}

func (b *Bot) getMessages(channel, expected string) bool {
	messages, err := b.s.ChannelMessages(channel, 10, "", "", "")
	if err != nil {
		return false
	}

	for _, message := range messages {
		for _, embed := range message.Embeds {
			if embed.Title == expected {
				return true
			}
		}
	}

	return false
}

func (b *Bot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := b.handleInteraction(s, i); err != nil {
		log.Error(err)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("something went wrong: %s", err),
			},
		})
		if err != nil {
			log.Error(err)
		}
	}
}

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		return b.onSlashCommand(s, i)
	default:
		return fmt.Errorf("invalid interaction type")
	}
}
