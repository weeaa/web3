package fren_utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/weeaa/nft/pkg/tls"
	"log"
	"os"
	"testing"
)

func TestGetUserInformation(t *testing.T) {
	address := "0xe5d60f8324d472e10c4bf274dbb7371aa93034a0"

	userInfo, err := GetUserInformation(address, tls.NewProxyLess())
	if err != nil {
		assert.Error(t, err)
	}

	if userInfo.Address != address {
		assert.Error(t, fmt.Errorf("expected %s, got %s", address, userInfo.Address))
	}

	assert.NoError(t, nil)
}

func TestW(t *testing.T) {
	followers := 133000
	i := AssertImportance(followers, "", Followers)
	log.Println(i)
}

func TestAssertImportanceFollowers(t *testing.T) {

	var tests = []struct {
		followers    int
		expectedImp  Importance
		expectedChan string
	}{
		{
			followers:   2000,
			expectedImp: "none",
		},
		{
			followers:   5000,
			expectedImp: Shrimp,
		},
		{
			followers:   18700,
			expectedImp: Fish,
		},
		{
			followers:   299000,
			expectedImp: Whale,
		},
	}

	for _, test := range tests {
		imp := AssertImportance(test.followers, Followers)

		if imp != test.expectedImp {
			t.Error(fmt.Errorf("expected %s, got %s", test.expectedImp, imp))
		}
	}

}

func TestAssertImportanceBalance(t *testing.T) {

}

func TestRedeemCodes(t *testing.T) {
	codes, err := RedeemCodes(os.Getenv("FT_BEARER_TOKEN"), tls.NewProxyLess())
	if err != nil {
		assert.Error(t, err)
	}

	if !(len(codes) > 0) {
		assert.Error(t, fmt.Errorf("expected len of codes higher than 0, got %v", len(codes)))
	}

	assert.NoError(t, nil)
}

func TestWatchList(t *testing.T) {
	if err := AddWishList("0xe5d60f8324d472e10c4bf274dbb7371aa93034a0", os.Getenv("FT_BEARER_TOKEN"), tls.NewProxyLess()); err != nil {
		assert.Error(t, err)
	}

	assert.NoError(t, nil)
}
