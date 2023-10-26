package twitter

import (
	"github.com/weeaa/nft/pkg/tls"
	"testing"
)

func TestFetchNitter(t *testing.T) {
	expectedUsername := "weea_a"

	resp, err := FetchNitter(expectedUsername, tls.NewProxyLess())
	if err != nil {
		t.Errorf("unexpected error [%v]", err)
	}

	if resp.Followers == "0" || resp.Followers == "" {
		t.Errorf("expected followers to be not equal to nil or 0")
	}
}
