package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/weeaa/nft/pkg/logger"
)

func (b *Bot) BotWebhook(webhook *discordgo.MessageSend, channelID string) {
	_, err := b.s.ChannelMessageSendComplex(channelID, webhook)
	if err != nil {
		logger.LogError(discord, err)
	}
}

func BundleQuickTaskComponents(target, module string) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Buy",
					Style: discordgo.LinkButton,
					URL:   fmt.Sprintf("http://localhost:3666/quickTask?module=%s&method=buy&target=%s", module, target),
					Emoji: discordgo.ComponentEmoji{
						Name: "↗️",
					},
				},
				discordgo.Button{
					Label: "Sell",
					Style: discordgo.LinkButton,
					URL:   fmt.Sprintf("http://localhost:3666/quickTask?module=%s&method=sell&target=%s", module, target),
					Emoji: discordgo.ComponentEmoji{
						Name: "↘️",
					},
				},
			},
		},
	}
}
func BundleQuickLinks(address string) string {
	return fmt.Sprintf("(BaseScan)[https://basescan.org/address/%s] | (FriendTech ChatRoom)[https://www.friend.tech/rooms/%s]", address, address)
}
