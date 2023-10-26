package bot

import (
	"os"
)

const discord = "Discord Bot"

const GuildID = "914250717722734652"

const Image = "https://images-ext-2.discordapp.net/external/uAW-8FcuQgvrhN3PL1TxNl19LfP9nSR2KhXOhnQLVr8/https/camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67"

const (
	CommunityPings = ""
)

// Colors
const (
	Purple = 0x9796F0
	Blue   = 0x026CDF
	Green  = 0x4DFF94
	Red    = 0xFD6157
	Orange = 0xFFBD31
)

// Channels
const (
	RolesChannel = "1157975432432455741"

	Authenticate = "1156489358011019264"
)

var (
	PremintChannel = os.Getenv("PREMINT_WEBHOOK")
)

// FriendTech Channels
const (
	FriendTech      = "friendTech"
	FriendTechImage = "https://content.fortune.com/wp-content/uploads/2023/08/friend.tech-logo.jpg"

	FriendTechNewUsers      = "1156497474630991903"
	FriendTechNewUsers5     = "1158558173779734548"
	FriendTechNewUsers10    = "1158541907035693066"
	FriendTechNewUsers50    = "1158438539202146375"
	FriendTechBalanceChange = "1159441853519777824"

	FriendTechFeed             = "1159124587808837816"
	FriendTechFeedPingsChannel = "1159420636360364062"

	FriendTechFilteredBuys  = "1156497675471048744"
	FriendTechFilteredSells = "1156497736380715029"

	FriendTechFishBuys  = "1161255638270226454"
	FriendTechFishSells = "1161255672244097065"

	FriendTechShrimpBuys  = "1161255555202039989"
	FriendTechShrimpSells = "1161255595173748737"

	FriendTechWhalesBuys  = "1156884576556306473"
	FriendTechWhalesSells = "1156884614871261214"

	FriendTechAllLogs = "1161256468180381696"
)

// StarsArena Channels
const (
	StarsArena            = "starsArena"
	StarsArenaFeed        = "1159858801080541284"
	StarsArenaFeedPing    = "1159858814313570315"
	StarsArenaNewUsers    = "1159631144980074526"
	StarsArenaNewUsers5k  = "1159631169764216902"
	StarsArenaNewUsers10k = "1159631191129989140"
	StarsArenaNewUsers50k = "1159631211057127516"
)

// Hub3 Channels
const (
	Hub3 = "hub3"

	Hub3TwitterUsername = "1160625679377911890"
	Hub3Map             = "1160633631518236702"
)
