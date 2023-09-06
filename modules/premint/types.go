package premint

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
)

const (
	moduleName = "premint.xyz"
)

var (
	maxRetriesReached = errors.New("maximum retries reached, aborting function")
	RateLimited       = errors.New("you are rate limited :( you got to wait till you're unbanned, which is approx 5+ minutes")
)

type RaffleType string

/* you need to hold a Premint NFT in order to access those eps */
var (
	Daily  RaffleType = "https://www.premint.xyz/collectors/explore/"
	Weekly RaffleType = "https://www.premint.xyz/collectors/explore/top/"
	New    RaffleType = "https://www.premint.xyz/collectors/explore/new"
)

type Settings struct {
	Discord *discord.Client
	Handler *handler.Handler
	Context context.Context
	Verbose bool
	Profile Profile
}

type Profile struct {
	RetryDelay      int
	publicAddress   string
	privateKey      string
	sessionId       string
	csrfToken       string
	nonce           string
	Client          tls_client.HttpClient
	ProxyList       []string
	RotateProxyOnRL bool
	isLoggedIn      bool
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
