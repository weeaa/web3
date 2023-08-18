# NFT Monitors 

An NFT monitoring toolkit for tracking NFT drops, sales, listings, and more.

<div align="center">
    <img src="https://cdn.discordapp.com/attachments/689063280358064158/1139538002041897041/image.png" margin="auto" height="270"/>
</div>

## Features

- [x] Etherscan Monitoring
  - [x] New Verified Contracts
- [x] ExchangeArt Monitoring
  - [x] New Drops by Artist
- [x] LMNFT Monitoring Top Drops
  - [x] Solana
  - [ ] Polygon
  - [ ] Ethereum
  - [ ] Binance
  - [ ] Aptos
  - [ ] Avalanche
  - [ ] Fantom
  - [ ] Stacks
- [x] OpenSea Monitoring
  - [x] Sales
  - [x] Listings
- [x] Premint Monitoring
  - [x] Hype Weekly/Daily Raffles (Premint NFT Required)
- [ ] BRC20 Unisat

## Getting Started

These instructions will guide you through setting up and running the NFT Monitors project on your local machine.

### Prerequisites

- GoLang 1.20 or later

### Usage

1. Clone the repository:

```bash
$ cd my-project
$ go get github.com/weeaa/nft
```

### Example

```go
func main() {
	
	client := discord.NewClient(
		os.Getenv("EXCHANGEART_WEBHOOK"),
		os.Getenv("LMNFT_WEBHOOK"),
		os.Getenv("PREMINT_WEBHOOK"),
		os.Getenv("ETHERSCAN_WEBHOOK"),
		os.Getenv("BRC20_WEBHOOK"),
		"weeaa's monitor",
		"https://image.png",
		0xFFFFF,
	)
	
	etherscan.Monitor(client)
	exchangeArt.Monitor(client, exchangeArt.DefaultList, false, 1000)
	lmnft.Monitor(client, []lmnft.Network{lmnft.Solana, lmnft.Binance}, 1000)
	
	profile := premint.NewProfile(os.Getenv("PREMINT_PUB"), os.Getenv("PREMINT_PRIV"), "", 5000)
	profile.Monitor(client, []premint.RaffleType{premint.Daily, premint.Weekly})
}
```

## Credits

s/o **DALL-E** for the image 😸
