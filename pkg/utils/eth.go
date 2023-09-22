package utils

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func CreateSubscription(client *ethclient.Client, addresses []common.Address, ch chan types.Log) (ethereum.Subscription, error) {
	return client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{Addresses: addresses}, ch)
}

func GetSender(tx *types.Transaction) (common.Address, error) {
	return types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
}

func GetWalletBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	return client.PendingBalanceAt(context.Background(), address)
}

func SliceToAddresses(wallets []string) []common.Address {
	var addresses []common.Address
	for _, wallet := range wallets {
		addresses = append(addresses, common.HexToAddress(wallet))
	}
	return addresses
}

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

type GetAllTxnParams struct {
}
