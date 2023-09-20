package unisat

import (
	"encoding/hex"
	"github.com/weeaa/nft/pkg/test"
	"github.com/weeaa/nft/pkg/tls"
	"net/http"
	"testing"
)

func TestGetHolders(t *testing.T) {
	tests := []struct {
		name string
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
			s := Settings{Client: tls.NewProxyLess()}
			holders, ok := s.FetchHolders(hex.EncodeToString([]byte("1024")), 10240000)
			if !ok {
				t.Errorf("expected no error, but got %v", holders)
			}
		case test_utils.InvalidResponse:
		}
	}
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
		case test_utils.InvalidResponse:
			server := test_utils.NewServer(http.StatusForbidden, ResFees{})

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Errorf("expected an error, but got nil")
			}

			if resp.StatusCode == 200 {
				t.Errorf("expected an non 200 code response")
			}

			fees := test_utils.ReadUnmarshal[ResFees](resp)

			if len(fees.FastestFee) != 0 || len(fees.HalfHourFee) != 0 || len(fees.HourFee) != 0 {
				t.Errorf("expected fees to be empty, but got %v", fees)
			}
		}

	}
}
