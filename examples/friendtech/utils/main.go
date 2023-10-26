package main

import (
	"fmt"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/pkg/tls"
	"sync"
	"time"
)

// modify 'bearer' with your bearer token, which you
// can find by using DevTools/capturing traffic.
const bearer = "quoicoubeh"

// !!!!! doesn't work – need to update code w/ client

// basic program to fetch your invite codes & adds to wishlist users you want
func main() {

	list := []string{
		"0xe5d60f8324d472e10c4bf274dbb7371aa93034a0",
	}

	client := tls.NewProxyLess()

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		for _, user := range list {
			if err := fren_utils.AddWishList(user, bearer, client); err != nil {
				fmt.Println("AddWishList", err)
			}
			fmt.Println("added to wishlist", user)
			// delay as you may get rate limited by FT if you follow 200+ persons
			// if lower than 100 you can run 0 delay
			time.Sleep(5 * time.Second)
		}
		wg.Done()
	}()

	codes, err := fren_utils.RedeemCodes(bearer, client)
	if err != nil {
		fmt.Println("RedeemCodes", err)
	}

	for _, code := range codes {
		fmt.Println(code)
	}

	wg.Wait()
	fmt.Println("program done, exiting...")
	time.Sleep(3 * time.Second)
}
