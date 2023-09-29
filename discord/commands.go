package discord

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) onSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
	}

	for _, command := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, GuildID, command)
	}

	return nil
}
