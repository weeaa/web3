package sniper

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/weeaa/nft/pkg/logger"
	"math/big"
	"os"
	"testing"
)

func TestSnipe(t *testing.T) {

	m := new(big.Float)
	m.SetString("1")

	s, err := New(os.Getenv("FT_PRIVATE_KEY"), os.Getenv("NODE_HTTP_URL"), m, 1)
	if err != nil {
		assert.Error(t, err)
	}

	sniped, err := s.Snipe(common.HexToAddress("0x9777ea3684e58fcf80734222d49fe57a9c5302da"), "1", Sell, 1, Normal)
	if err != nil {
		assert.Error(t, fmt.Errorf("error sniping: %w", err))
	}

	sniped.Txns.ForEach(func(transaction *types.Transaction, err error) {
		if err != nil {
			assert.Error(t, err)
		}
	})

	logger.LogInfo(sniper, fmt.Sprint(sniped))

	assert.NoError(t, nil)
}
