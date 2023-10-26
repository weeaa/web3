package base

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
)

type Settings struct {
	Discord       *discord.Client
	Handler       *handler.Handler
	Verbose       bool
	Context       context.Context
	Client        *ethclient.Client
	MonitorParams MonitorParams
}

type MonitorParams struct {
	BlacklistedTokens []string
}
