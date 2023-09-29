package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/db"
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

	bot := &Bot{s, db}

	s.AddHandler(bot.onReady)
	//	s.AddHandler(bot.onInteractionCreate)
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
	if r.Type != "READY" {
		log.Error("not ready")
		return
	}

	if err := s.UpdateListeningStatus("rugging"); err != nil {
		log.Error(err)
		return
	}

}
