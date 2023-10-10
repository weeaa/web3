package utils

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

// WeiToEther actually works with most currencies
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

func GenerateEthWallet() *Wallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  common.HexToAddress(crypto.PubkeyToAddress(*publicKeyECDSA).Hex()),
	}
}

// DecodeTransactionInputData returns you the tx data input (like address, amount, etc)
func DecodeTransactionInputData(contractABI abi.ABI, txData string) (map[string]any, error) {
	sig, err := hex.DecodeString(txData[2:10])
	if err != nil {
		return nil, err
	}

	method, err := contractABI.MethodById(sig)
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(txData[10:])
	if err != nil {
		return nil, err
	}

	inputs := make(map[string]any, 0)
	err = method.Inputs.UnpackIntoMap(inputs, data)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func DisperseEthFunds(privateKey string, addresses []common.Address, amount *big.Int, client *ethclient.Client) ([]*types.Transaction, error) {
	var transactions []*types.Transaction
	wallet := InitWallet(privateKey)

	for _, address := range addresses {

		nonce, err := client.PendingNonceAt(context.Background(), wallet.PublicKey)
		if err != nil {
			return transactions, err
		}

		//value := big.NewInt(1000000000000000000) // in wei (1 eth)
		//gasLimit := uint64(21000)                // in units
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return transactions, err
		}

		//var data []byte
		//txx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, data)

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			return transactions, err
		}

		tx, err := types.SignNewTx(wallet.PrivateKey, types.NewLondonSigner(chainID), &types.DynamicFeeTx{
			ChainID: chainID,
			Nonce:   nonce,
			To:      &address,
			Value:   amount,
			Gas:     gasPrice.Uint64(),
		})
		if err != nil {
			return transactions, err
		}

		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, tx)
	}
	return transactions, nil
}

func ConsolidateEthFunds(privateKeys []string, address common.Address, amount *big.Int, client *ethclient.Client) ([]*types.Transaction, error) {
	var transactions []*types.Transaction
	for _, privateKey := range privateKeys {
		wallet := InitWallet(privateKey)

		if amount.Uint64() == 0 {
			var err error
			amount, err = GetEthWalletBalance(client, wallet.PublicKey)
			if err != nil {
				return transactions, err
			}
		}

		nonce, err := client.PendingNonceAt(context.Background(), wallet.PublicKey)
		if err != nil {
			return transactions, err
		}

		//gasLimit := uint64(21000) // in units
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			return transactions, err
		}

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			return transactions, err
		}

		tx, err := types.SignNewTx(wallet.PrivateKey, types.NewLondonSigner(chainID), &types.DynamicFeeTx{
			ChainID: chainID,
			Nonce:   nonce,
			To:      &address,
			Value:   amount,
			Gas:     gasPrice.Uint64(),
		})
		if err != nil {
			return transactions, err
		}

		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, tx)
	}
	return transactions, nil
}
