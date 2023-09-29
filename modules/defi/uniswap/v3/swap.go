package v3

import (
	"context"
	"fmt"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func InitSwapInstance() {
	// call eth client b4
	f, err := contract.NewUniswapv3Factory(common.HexToAddress(ContractV3Factory), client)
	if err != nil {

	}
}

func (s *SwapInstance) sendTransaction(toAddress common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasLimit, err := s.client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     s.Wallet.PublicKey,
		To:       &toAddress,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("gasLimit=%d,  gasPrice=%d\n", gasLimit, gasPrice.Uint64())
	nounc, err := s.client.NonceAt(context.Background(), s.Wallet.PublicKey, nil)
	if err != nil {
		return nil, err
	}

	tx := types.NewTransaction(nounc, toAddress, value,
		gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), w.PrivateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
