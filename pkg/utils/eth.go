package utils

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

func InitWallet(privateStrKey string) *Wallet {
	privateKey, err := crypto.HexToECDSA(privateStrKey)
	if err != nil {
		return nil
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).SetPrec(236).SetMode(big.ToNearestEven).Quo(new(big.Float).SetPrec(236).SetMode(big.ToNearestAway).SetInt(wei), big.NewFloat(params.Ether))
}

func CreateSubscription(client *ethclient.Client, addresses []common.Address, ch chan types.Log) (ethereum.Subscription, error) {
	return client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{Addresses: addresses}, ch)
}

func GetSender(tx *types.Transaction) (common.Address, error) {
	return types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
}

func GetEthWalletBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	return client.BalanceAt(context.Background(), address, nil)
}

func SliceToAddresses(wallets []string) []common.Address {
	var addresses []common.Address
	for _, wallet := range wallets {
		addresses = append(addresses, common.HexToAddress(wallet))
	}
	return addresses
}

// todo finish func
func GetAllTransactions(client *ethclient.Client, address common.Address, params GetAllTxnParams) error {
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	for i := latestBlock.NumberU64() - 300; i <= latestBlock.NumberU64(); i++ {
		for _, tx := range latestBlock.Transactions() {
			_ = tx
		}
	}
	return nil
}
