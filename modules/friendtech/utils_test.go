package friendtech

import (
	"fmt"
	"github.com/weeaa/nft/pkg/tls"
	"log"
	"testing"
)

func TestRedeemCodes(t *testing.T) {

	u := &Account{
		client: tls.NewProxyLess(),
	}

	codes, err := u.RedeemCodes()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("codes", codes)
}

func TestAssertImportance(t *testing.T) {
	var expectedImp = Whale

	imp, err := assertImportance(87000, Followers)
	if err != nil {
		t.Error(t, err)
	}

	if imp != expectedImp {
		t.Error(fmt.Errorf("expected %s, got %s", expectedImp, imp))
	}
}
