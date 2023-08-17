package premint

import (
	"github.com/PuerkitoBio/goquery"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
)

const moduleName = "Premint.xyz"

type RaffleType string

/* you need to hold a Premint NFT in order to access those eps */
var (
	Daily  RaffleType = "https://www.premint.xyz/collectors/explore/"
	Weekly RaffleType = "https://www.premint.xyz/collectors/explore/top/"
	New    RaffleType = "https://www.premint.xyz/collectors/explore/new"
)

type Profile struct {
	RetryDelay    int
	publicAddress string
	privateKey    string
	sessionId     string
	csrfToken     string
	nonce         string
	handler       handler.Handler
	Client        tls_client.HttpClient
	DiscordClient discord.Client
}

type Webhook struct {
	document *goquery.Document

	Title        string
	Slug         string
	Image        string
	Desc         string
	Price        string
	BalanceFall  string
	ETHtoHold    string
	TimeClose    string
	WinnerAmount string
	Status       string
	StatusImg    string

	Twitter TwitterReqs
	Discord DiscordReqs
	Misc    MiscReqs
	Custom  Custom
}

type TwitterReqs struct {
	Total   string
	Account string
	Tweet   string
}

type DiscordReqs struct {
	Total  string
	Server string
	Role   string
}

type MiscReqs struct {
	Total          string
	Spots          string
	OverAllocating string
	RegOut         string
	LinkOut        string
}

type Custom struct {
	Total string
}
