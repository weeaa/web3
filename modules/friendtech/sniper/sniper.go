package sniper

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/safemap"
	"github.com/weeaa/nft/pkg/tls"
	ethereum2 "github.com/weeaa/nft/pkg/utils/ethereum"
	"math/big"
	"strings"
	"sync"
)

// default testing data
const (
	gasLimit = 3_000_000
	gasPrice = 300 * params.GWei
)

// New initializes a sniping instance.
func New(privateKey, HTTPNodeUrl string, maxEthInput *big.Float, maxShareBuy int) (*Sniper, error) {
	httpClient, err := ethclient.Dial(HTTPNodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to http node: %w", err)
	}

	friendTechABI, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		return nil, fmt.Errorf("error reading abi: %w", err)
	}

	contract := bind.NewBoundContract(FRIEND_TECH_CONTRACT, friendTechABI, httpClient, httpClient, httpClient)

	return &Sniper{
		Bind:        contract,
		ABI:         friendTechABI,
		PrivateKey:  privateKey,
		Wallet:      ethereum2.InitWallet(privateKey),
		HttpClient:  tls.NewProxyLess(),
		Client:      httpClient,
		MaxEthInput: maxEthInput,
		MaxShareBuy: maxShareBuy,
	}, nil
}

func NewSnipeTransactions() *SnipeTransactions {
	return &SnipeTransactions{Txns: safemap.New[*types.Transaction, error]()}
}

// Snipe logic is...
// We fetch the share price of the ft_user we aim to snipe
// & compare it with max eth input; if higher we abort the snipe.
func (s *Sniper) Snipe(user common.Address, amountShares string, functionName Function, buyTimes int, mode Mode) (*SnipeTransactions, error) {

	userInfo, errUserInfo := fren_utils.GetUserInformation(user.String(), s.HttpClient)
	if errUserInfo != nil {
		return nil, errUserInfo
	}

	if ok := isSharePriceOk(s.MaxEthInput, userInfo.DisplayPrice); !ok {
		return nil, fmt.Errorf("error validating share price, trade aborted [sharePrice: %s| maxEthInput: %f]", userInfo.DisplayPrice, s.MaxEthInput)
	}

	snipeTransactions := NewSnipeTransactions()

	if mode == Normal {
		// Normal mode only sends one transaction.

		txn, err := s.buildTx(user, functionName, amountShares, userInfo.DisplayPrice)
		if err != nil {
			return nil, err
		}

		err = s.Client.SendTransaction(context.Background(), txn)
		snipeTransactions.Txns.Set(txn, err)
	} else {
		// If spam mode, we spam transactions per user requests.
		// Transactions broadcasted that are not in the same block
		// will be reverted.

		wg := sync.WaitGroup{}
		var counter int
		var blockNumber uint64
		_ = blockNumber
		for i := 0; i < buyTimes; i++ {
			go func() {
				wg.Add(1)
				defer wg.Done()

				txn, err := s.buildTx(user, functionName, amountShares, userInfo.DisplayPrice)
				if err != nil {
					logger.LogError(sniper, err)
					return
				}

				err = s.Client.SendTransaction(context.Background(), txn)
				snipeTransactions.Txns.Set(txn, err)
				counter++
			}()
		}
		wg.Wait()
	}

	return snipeTransactions, nil
}

// buildTx assembles a FriendTech transaction using the provided inputs.
func (s *Sniper) buildTx(address common.Address, functionName Function, qty, price string) (*types.Transaction, error) {
	var err error
	var totalValueDecimal decimal.Decimal
	var totalValue *big.Int

	funcStr := functionStr(functionName)

	function, ok := s.ABI.Methods[funcStr]
	if !ok {
		return nil, err
	}

	quantity, err := decimal.NewFromString(qty)
	if err != nil {
		return nil, fmt.Errorf("error str to dec. qty: %w", err)
	}

	if function.IsPayable() {
		if functionName == Sell {
			totalValue = big.NewInt(0)
		} else {
			totalValueDecimal, err = decimal.NewFromString(price)
			if err != nil {
				return nil, err
			}

			totalValue = big.NewInt((totalValueDecimal.IntPart() * 5) * quantity.IntPart())
		}
	}

	var filters []map[string]string
	for _, method := range s.ABI.Methods {
		if method.StateMutability == "view" {
			if len(method.Outputs) == 1 && len(method.Inputs) == 0 {
				if method.Outputs[0].Type.String() == "address" {
					var OwnerAddress []interface{}
					if err = s.Bind.Call(nil, &OwnerAddress, method.Name); err != nil {
						return nil, err
					}

					filters = append(filters, map[string]string{
						"from":                    OwnerAddress[0].(common.Address).String(),
						"to":                      FRIEND_TECH_CONTRACT.String(),
						"contractCall.methodName": funcStr,
					})

				}
			}
		}
	}

	if len(filters) == 0 {
		return nil, fmt.Errorf("error filters equal to nil")
	}

	var args []interface{}
	for _, method := range s.ABI.Methods[funcStr].Inputs {
		switch method.Type.String() {
		case "address":
			args = append(args, address)
		case "uint256":
			args = append(args, quantity.BigInt())
		default:
			return nil, fmt.Errorf("unknown method")
		}
	}

	contractCallData, err := s.ABI.Pack(funcStr, args...)
	if err != nil {
		return nil, err
	}

	nonce, err := s.Client.NonceAt(context.Background(), s.Wallet.PublicKey, nil)
	if err != nil {
		return nil, err
	}

	chainID, err := s.Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	gasTipCap, err := s.Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}

	pendingHeader, err := s.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	gasFeeCap := new(big.Int).Add(gasTipCap, new(big.Int).Mul(pendingHeader.BaseFee, big.NewInt(2)))

	gas, err := s.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From: s.Wallet.PublicKey,
		To:   &FRIEND_TECH_CONTRACT,
		Data: contractCallData,
	})

	return types.SignNewTx(s.Wallet.PrivateKey, types.NewLondonSigner(chainID), types.TxData(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       gas * params.Wei,
		To:        &FRIEND_TECH_CONTRACT,
		Value:     totalValue,
		Data:      contractCallData,
	}))
}

// revertTxn builds a tx with same nonce but with 0 ETH.
func (s *Sniper) revertTxn(nonce string, tx *types.Transaction) error {

	return nil
}

// ReadOrders reads orders created and subscribe to the socket.
func (s *Sniper) ReadOrders() {

}

func (s *Sniper) RemoveOrder(orderID string) error {
	// find order ID
	return nil
}

// CreateLimitBuyOrder creates and stores a limit order buy.
func (s *Sniper) CreateLimitBuyOrder(address string, shareMaxPrice *big.Int, expiration uint32) error {
	if expiration < 0 {
		return fmt.Errorf("invalid expiration time < 0 [%d]", expiration)
	}

	userInfo, err := validateAddress(address)
	if err != nil {
		return err
	}

	if userInfo.Address != address {
		return fmt.Errorf("CreateLimitBuyOrder mismatching address: expected %s, got %s", address, userInfo.Address)
	}

	// store db

	return nil
}

// why? cause some persons may run "snipers" wallets, and these don't exist on Friend Tech.
func validateAddress(address string) (fren_utils.UserInformation, error) {
	return fren_utils.GetUserInformation(address, tls.NewProxyLess())
}

func (s *Sniper) watchTxn(txHash common.Hash, blockNumber *big.Int) error {
	txReceipt, err := s.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return err
	}
	_ = txReceipt
	//if txReceipt.BlockNumber
	return nil
}
