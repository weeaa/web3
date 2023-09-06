package main

import (
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/modules/etherscan"
)

func main() {
	c := make(chan struct{})

	etherscanClient := etherscan.NewClient(&discord.Client{ProfileName: "vuitton"}, true)
	etherscanClient.Discord.Webhook = "https://discord.com/api/webhooks/1148667139415359539/yM8bK3k2DL7NmzVgR6KH_0PUz4cpbQs2rQCITENElRi4XFudyE7fIRRFd6Gp1mwubk7M"

	etherscanClient.StartMonitor()

	<-c
}
