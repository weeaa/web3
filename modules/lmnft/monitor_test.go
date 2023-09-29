package lmnft

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMonitorDrops(t *testing.T) {
	tests := []struct {
		name    string
		network Network
		link    string
	}{
		{
			name:    "solana",
			network: Solana,
			link:    "https://www.launchmynft.io/collections/9PTwLrfoYpY9Ao9cKwXTenVLngx8pAXeP841G6jQ7o7P/ny9QAw8DzGg8ClImG0oP",
		},
		{
			name:    "binance",
			network: Binance,
			link:    "",
		},
		{
			name:    "ethereum",
			network: Ethereum,
			link:    "",
		},
	}

	for _, test := range tests {
		resp, err := doRequest([]Network{test.network})
		if err != nil {
			assert.Error(t, err)
		}

		if resp.StatusCode != 200 {
			assert.Error(t, fmt.Errorf(""))
		}

		if err = resp.Body.Close(); err != nil {
			assert.Error(t, err)
		}

		switch test.network {
		case Solana:
			_, err = scrapeInformation[resSolana](test.link)
			if err != nil {
				assert.Error(t, err)
			}
		case Binance:
		case Sui:
		case Ethereum:
			//todo add others
		}
	}
}
