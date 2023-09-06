package etherscan

import (
	"context"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
)

const (
	moduleName = "Etherscan Verified Contract"
	retryDelay = 3000
)

type Settings struct {
	Discord *discord.Client
	Handler *handler.Handler
	Context context.Context
	Verbose bool
}

type Contract struct {
	Address string
	Name    string
	Link    string
}
