package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/weeaa/nft/pkg/api"
	"github.com/weeaa/nft/pkg/logger"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "add_user",
			Description: "Adds a user to our monitors.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "base_address",
					Description: "Base Address you want to add.",
					Required:    true,
				},
			},
		},
		{
			Name:        "monitor_new_user",
			Description: "Monitors users joining Friend Tech.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "twitter_name",
					Description: "Twitter Name of the user you want to monitor (i.e: weea_a)",
					Required:    true,
				},
			},
		},
	}
)

// registerCommands registers slash commands.
func (b *Bot) registerCommands() {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		cmd, err := b.s.ApplicationCommandCreate(b.s.State.User.ID, GuildID, command)
		if err != nil {
			logger.LogError(discord, fmt.Errorf("cannot create [%s] command: %w", command.Name, err))
		}
		registeredCommands[i] = cmd
	}
}

// onSlashCommand is a handler: whenever a user performs a /slash
// command, it will execute it.
func (b *Bot) onSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Name {
	case "monitor_new_user":
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		args := make([]interface{}, 0, len(options))
		if option, ok := optionMap["twitter_name"]; ok {
			args = append(args, option.StringValue())
		}

		// add to database & start monitoring

		_, err := b.s.ChannelMessageSendComplex(FriendTechFeed, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "🔔 | New User Added",
					Description: fmt.Sprintf("**[%s](https://x.com/%s)** is now monitored on Friend Tech & waiting until he joins.", optionMap["twitter_name"].StringValue(), optionMap["twitter_name"].StringValue()),
					Color:       Purple,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "@friendtech — feed",
						IconURL: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
					},
				},
			},
		})
		if err != nil {
			return b.ReturnErrorInteraction(i, err)
		} else {
			return b.ReturnConfirmationInteraction(i)
		}

	case "add_user":
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		userInfo, err := api.AddUserToMonitor(optionMap["base_address"].StringValue(), "weeaa")
		if err != nil {
			return b.ReturnErrorInteraction(i, err)
		}

		if _, err = b.s.ChannelMessageSendComplex(FriendTechFeed, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "🎩 | add_user",
					Description: fmt.Sprintf("**[%s](https://x.com/%s)** is now monitored on Friend Tech.\n\n__Audit__\n > Imp. Status: **%s**\n> Followers: **%s**\n> ChatRoom: **[Link](https://www.friend.tech/rooms/%s)**", userInfo["twitter_name"], userInfo["twitter_username"], fmt.Sprint(userInfo["status"]), fmt.Sprint(userInfo["followers"]), optionMap["base_address"].StringValue()),
					Color:       Purple,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: fmt.Sprint(userInfo["image"]),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("@friendtech — feed [%s]", optionMap["base_address"].StringValue()),
						IconURL: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
					},
				},
			},
		}); err != nil {
			return b.ReturnErrorInteraction(i, err)
		}

		return b.ReturnConfirmationInteraction(i)
	default:
		return fmt.Errorf("unknown slash command")
	}

}
