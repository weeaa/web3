package friendtech

import (
	"context"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
)

const moduleName = "Friend Tech"
const FRIEND_TECH_CONTRACT_V1 = "0xcf205808ed36593aa40a44f10c7f7c2f67d4a4d4"

const (
	sellMethod = "0xb51d0534"
	buyMethod  = "0x6945b123"
)

const (
	prodBaseApi = "prod-api.kosetto.com"
	iphoneUA    = "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1"

	mainnet_rpc = "https://mainnet-sequencer.base.org/"
)

type Settings struct {
	Discord    *discord.Client
	Handler    *handler.Handler
	Context    context.Context
	WSSClient  *ethclient.Client
	HTTPClient *ethclient.Client
	Verbose    bool
	ABI        abi.ABI
	//Sniper     Sniper
	DB *db.DB

	LatestUserID  int
	isLatestFound bool

	BuysFilterWebhook  string
	SellsFilterWebhook string
	BuysWebhook        string
	SellsWebhook       string
}

type Indexer struct {
	DB          *db.DB
	UserCounter int
	ProxyList   []string
	Client      tls_client.HttpClient
}

type Account struct {
	Email    string
	Password string
	Bearer   string
	client   tls_client.HttpClient
}

type UserInformation struct {
	Id                         int    `json:"id"`
	Address                    string `json:"address"`
	TwitterUsername            string `json:"twitterUsername"`
	TwitterName                string `json:"twitterName"`
	TwitterPfpUrl              string `json:"twitterPfpUrl"`
	TwitterUserId              string `json:"twitterUserId"`
	LastOnline                 int64  `json:"lastOnline"`
	LastMessageTime            string `json:"lastMessageTime"`
	HolderCount                int    `json:"holderCount"`
	HoldingCount               int    `json:"holdingCount"`
	WatchlistCount             int    `json:"watchlistCount"`
	ShareSupply                int    `json:"shareSupply"`
	DisplayPrice               string `json:"displayPrice"`
	LifetimeFeesCollectedInWei string `json:"lifetimeFeesCollectedInWei"`
}

// UserByIDResponse is returned by the /users/by-id/ endpoint
type UserByIDResponse struct {
	Id                         int    `json:"id"`
	Address                    string `json:"address"`
	TwitterUsername            string `json:"twitterUsername"`
	TwitterName                string `json:"twitterName"`
	TwitterPfpUrl              string `json:"twitterPfpUrl"`
	TwitterUserId              string `json:"twitterUserId"`
	LastOnline                 int    `json:"lastOnline"`
	LifetimeFeesCollectedInWei string `json:"lifetimeFeesCollectedInWei"`
}

type Importance string
type ImpType string

var (
	Whale  Importance = "high"
	Fish   Importance = "medium"
	Shrimp Importance = "low"
)

var (
	Balance   ImpType = "balance"
	Followers ImpType = "followers"
)
