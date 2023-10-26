package watcher

import (
	"context"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/twitter"
	"github.com/weeaa/nft/pkg/cache"
)

const watcher = "Friend Tech Watcher"

const (
	sellMethod = "0xb51d0534"
	buyMethod  = "0x6945b123"
)

type Watcher struct {
	DB *db.DB

	// WSSClient is your Node Client Websocket conn.
	WSSClient     *ethclient.Client
	HTTPClient    *ethclient.Client
	NitterClient  *twitter.Client
	WatcherClient tls_client.HttpClient

	// Cache is used to store self-sells
	// and detect potential rugs/exploits.
	Cache *cache.Handler

	// Addresses stores Base addresses you want to monitor
	// fetched directly from the database.
	Addresses map[string]string

	// Counter represents the most recent FriendTech userID
	// that will be observed, as a starting ID, for new users.
	Counter int

	// OutStreamData is the data sent to our websocket.
	OutStreamData chan BroadcastData

	ABI abi.ABI

	// Pool is a 'watching' pool where new users that
	// have not deposited yet ETH to their wallet
	// will be 'watched' until they deposit, for x time.
	Pool       chan string
	EnablePool bool
	Deadline   float64

	ProxyList     []string
	Bot           *bot.Bot
	NewUsersCtx   context.Context
	PendingDepCtx context.Context
}

const (
	NewSignup    = "new_signup"
	BuyFiltered  = "buy_filtered"
	SellFiltered = "sell_filtered"
)

type BroadcastData struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
