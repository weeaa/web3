package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  common.Address
}

type GetAllTxnParams struct {
}
