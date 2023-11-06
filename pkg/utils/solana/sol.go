package solana

import (
	"github.com/gagliardetto/solana-go"
)

func SliceToPrograms(wallets []string) []solana.PublicKey {
	var addresses []solana.PublicKey
	for _, wallet := range wallets {
		addresses = append(addresses, solana.MustPublicKeyFromBase58(wallet))
	}
	return addresses
}

func LmpToSol() {}
