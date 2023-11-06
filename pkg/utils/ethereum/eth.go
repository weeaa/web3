package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"io"
	"math/big"
	"net/url"
)

var (
	ErrNotContract = errors.New("address is not a contract")
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  common.Address
}

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

func HexEncodeTxData(data []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(data))
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

func GenerateEthWallet() (*Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting ecdsa public key")
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  common.HexToAddress(crypto.PubkeyToAddress(*publicKeyECDSA).Hex()),
	}, nil
}

// DecodeTransactionInputData returns the transaction data input (address, amount, etc).
func DecodeTransactionInputData(contractABI abi.ABI, txData string) (map[string]any, error) {
	meth, err := hex.DecodeString(txData[2:10])
	if err != nil {
		return nil, err
	}

	method, err := contractABI.MethodById(meth)
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

func EthToUsd() (map[string]string, error) {
	rate, err := FetchRate()
	if err != nil {
		return nil, err
	}
	_ = rate
	return make(map[string]string), nil
}

func FetchRate() (map[string]any, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "api.coingecko.com", Path: "/api/v3/simple/price?ids=ethereum&vs_currencies=usd%2Ceur%2Cchf%2Ccad"},
		Header: http.Header{
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"accept":             {"application/json, text/javascript, */*; q=0.01"},
			"referer":            {"https://fr.beincrypto.com/"},
			"dnt":                {"1"},
			"sec-ch-ua-mobile":   {"?0"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
			"sec-ch-ua-platform": {"\"macOS\""},
		},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client error: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]any
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// IsContract pulls code from an address, if there is no code it is an account address.
func IsContract(client *ethclient.Client, address common.Address) ([]byte, error) {
	code, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	if len(code) == 0 {
		return nil, ErrNotContract
	} else {
		return code, nil
	}
}

func DisperseFunds(privateKey string, addresses []common.Address, amount *big.Int, client *ethclient.Client) ([]*types.Transaction, error) {
	var transactions []*types.Transaction
	wallet := InitWallet(privateKey)

	for _, address := range addresses {

		nonce, err := client.PendingNonceAt(context.Background(), wallet.PublicKey)
		if err != nil {
			return transactions, err
		}

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

func ConsolidateFunds(privateKeys []string, address common.Address, amount *big.Int, client *ethclient.Client) ([]*types.Transaction, error) {
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
