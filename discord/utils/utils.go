package discord_utils

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func BundleUserInformation() {

}

// isAllowed checks if the user performing the request has the right roles.
func isAllowed(i *discordgo.InteractionCreate, s *discordgo.Session) error {
	member, err := s.GuildMember("GuildID", i.User.ID)
	if err != nil {
		return err
	}

	hasRequiredRole := false
	for _, userRole := range member.Roles {
		if userRole == "1156498768322117664" {
			hasRequiredRole = true
			break
		}
	}

	if !hasRequiredRole {
		return fmt.Errorf("invalid roles, you are not allowed to perform this command")
	}
	return nil
}
