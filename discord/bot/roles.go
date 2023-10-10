package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/weeaa/nft/pkg/logger"
	"log"
)

var EmojiRoleMap = map[string]string{
	"🫱🏻‍🫲🏾": "1157981202171576360",
	"👶":     "1158029574563700757",
	"🐋":     "1157981248719945728",
	"🐠":     "1157981281167089684",
	"🦐":     "1157981304114139197",
	"🐰":     "1159193600886837371",
}

func (b *Bot) messageRoleChannel() {
	embed := &discordgo.MessageEmbed{
		Title:       "👤 — Roles",
		Description: "> \"\U0001FAF1🏻‍\U0001FAF2🏾\" designates the \"Community Pings\" role for members who ping others when they spot profitable opportunities within the monitor feed.\n\n    > \"👶\" bestows the \"New Users\" role upon users who sign up on Friend Tech with a substantial followers amount.\n\n   > \"🐋\" assigns the \"Whale\" role.\n\n    > \"🐠\" grants the \"Fish\" role.\n\n    > \"🦐\" assigns the \"Shrimp\" role.",
		Color:       Purple,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "@weeaa — roles",
			IconURL: "https://pbs.twimg.com/profile_images/1706780390210347008/dJSxjBGv_400x400.jpg",
		},
	}

	msgSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
		//Components: components,
	}

	m, err := b.s.ChannelMessageSendComplex(RolesChannel, msgSend)
	if err != nil {
		logger.LogError(discord, err)
	}

	for em := range EmojiRoleMap {
		b.s.MessageReactionAdd(RolesChannel, m.ID, em)
	}

}

func (b *Bot) onRoleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	if roleID, ok := EmojiRoleMap[r.Emoji.Name]; ok {
		err := s.GuildMemberRoleAdd(r.GuildID, r.UserID, roleID)
		if err != nil {
			log.Println("Error adding role to user:", err)
		}
	}
}

func (b *Bot) onRoleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {

	if roleID, ok := EmojiRoleMap[r.Emoji.Name]; ok {
		err := s.GuildMemberRoleRemove(r.GuildID, r.UserID, roleID)
		if err != nil {
			log.Println("Error removing role from user:", err)
		}
	}
}
