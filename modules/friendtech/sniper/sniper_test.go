package sniper

import (
	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"testing"
)

func TestSnipe(t *testing.T) {

	m := new(big.Float)
	m.SetString("1")

	s, err := New(os.Getenv("FT_PRIVATE_KEY"), os.Getenv("NODE_HTTP_URL"), m, 1)
	if err != nil {
		log.Error(err)
	}

	txns, err := s.Snipe(common.HexToAddress("0x9777ea3684e58fcf80734222d49fe57a9c5302da"), "1", Sell, 1, Normal)
	if err != nil {
		log.Error(err)
	}

	log.Info("sniped successful", txns)
}
