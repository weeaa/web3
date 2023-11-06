package discord_utils

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

var (
	ErrNotEnoughRights = errors.New("error: you don't have enough rights to perform this command")
)

// IsAllowed checks if the user performing the request has the right role(s).
func IsAllowed(i *discordgo.InteractionCreate, s *discordgo.Session, rolesExpected []string) error {
	member, err := s.GuildMember("GuildID", i.User.ID)
	if err != nil {
		return err
	}

	for _, userRole := range member.Roles {
		for _, role := range rolesExpected {
			if role == userRole {
				return nil
			}
		}
	}

	return ErrNotEnoughRights
}
