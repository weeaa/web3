package unisat

import (
	"github.com/stretchr/testify/assert"
	"github.com/weeaa/nft/pkg/test"
	"github.com/weeaa/nft/pkg/tls"
	"log"
	"testing"
)

func TestGetHolders(t *testing.T) {
	ticker := "66736174"
	supply := 500

	client := NewClient(nil, false, tls.NewProxyLess(), nil, false)

	holders, err := client.FetchHolders(ticker, supply)
	log.Println("err", err)
	if err != nil {
		assert.Error(t, err)
	}

	//https://api.unisat.io/query-v4/brc20/66736174/holders?start=0&limit=20
	//https://api.unisat.io/query-v4/brc20/66736174/holders?start=0&limit=5

	log.Println(holders)
}

func TestGetFees(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name: test_utils.ValidResponse,
		},
		{
			name: test_utils.InvalidResponse,
		},
	}

	for _, test := range tests {
		switch test.name {
		case test_utils.ValidResponse:
			fees, err := GetFees()
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if len(fees.FastestFee) == 0 || len(fees.HalfHourFee) == 0 || len(fees.HourFee) == 0 {
				t.Errorf("expected fees to be valid, but got %v", fees)
			}

		}
	}
}
