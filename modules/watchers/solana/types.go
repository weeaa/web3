package solana

import (
	"context"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

const moduleName = "Solana Wallet Watcher"

type Settings struct {
	Client  *ws.Client
	Context context.Context
}
