package watcher

import (
	"context"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"time"
)

const watcher = "Friend Tech Watcher"

const (
	sellMethod = "0xb51d0534"
	buyMethod  = "0x6945b123"
)

type Watcher struct {
	DB         *db.DB
	WSSClient  *ethclient.Client
	HTTPClient *ethclient.Client
	Client     tls_client.HttpClient
	Addresses  map[string]string // Base Addresses you want to monitor
	Counter    int

	OutStreamData chan BroadcastData

	ABI  abi.ABI
	Pool chan string

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

type L2TransactionsPendingResponse struct {
	Items []struct {
		L1BlockNumber    int       `json:"l1_block_number"`
		L1BlockTimestamp time.Time `json:"l1_block_timestamp"`
		L1TxHash         string    `json:"l1_tx_hash"`
		L1TxOrigin       string    `json:"l1_tx_origin"`
		L2TxGasLimit     string    `json:"l2_tx_gas_limit"`
		L2TxHash         string    `json:"l2_tx_hash"`
	} `json:"items"`
	NextPageParams struct {
		ItemsCount    int    `json:"items_count"`
		L1BlockNumber int    `json:"l1_block_number"`
		TxHash        string `json:"tx_hash"`
	} `json:"next_page_params"`
}
