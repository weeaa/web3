package main

import (
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/modules/etherscan"
	"github.com/weeaa/nft/modules/exchangeArt"
	"github.com/weeaa/nft/modules/unisat"
	"github.com/weeaa/nft/pkg/tls"
)

func main() {
	c := make(chan struct{})

	proxyList, err := tls.ReadProxyFile("my/path/to/proxy/txt/file")
	if err != nil {
		log.Fatal(err)
	}

	ethscan := etherscan.NewClient(&discord.Client{
		ProfileName: "Etherscan",
		AvatarImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Color:       0x00000,
		FooterText:  "made by weeaa",
		FooterImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Webhook:     "https://discord.com/api/webhooks/1148667139415359539/yM8bK3k2DL7NmzVgR6KH_0PUz4cpbQs2rQCITENElRi4XFudyE7fIRRFd6Gp1mwubk7M",
	}, true)

	exchArt := exchangeArt.NewClient(&discord.Client{
		ProfileName: "Exchange Art",
		AvatarImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Color:       0x00000,
		FooterText:  "made by weeaa",
		FooterImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Webhook:     "https://discord.com/api/webhooks/1148667139415359539/yM8bK3k2DL7NmzVgR6KH_0PUz4cpbQs2rQCITENElRi4XFudyE7fIRRFd6Gp1mwubk7M",
	}, false, true, 1500)

	uni := unisat.NewClient(&discord.Client{
		ProfileName: "Unisat",
		AvatarImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Color:       0x00000,
		FooterText:  "made by weeaa",
		FooterImage: "https://camo.githubusercontent.com/a0d06e6da8dcc033e33c2694eb550ffb775a3f805c7e2edd55758275a0862dd4/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3638393036333238303335383036343135382f313133393533383030323034313839373034312f696d6167652e706e67",
		Webhook:     "https://discord.com/api/webhooks/1148667139415359539/yM8bK3k2DL7NmzVgR6KH_0PUz4cpbQs2rQCITENElRi4XFudyE7fIRRFd6Gp1mwubk7M",
	}, true, tls.New(tls.NewProxy(tls.RandProxyFromList(proxyList))), proxyList, false)

	exchangeArtArtists := &[]string{
		"WQh4eWseb7ObpFVGFxTov08HlSr2",
		"wnzwvrxRolf9qXSzWuC9rN4XbcB2",
	}

	ethscan.StartMonitor()
	exchArt.StartMonitor(exchangeArtArtists)
	uni.StartMonitor()

	<-c // blocks forever because code runs in goroutines
}
