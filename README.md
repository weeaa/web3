# web3

A project encompassing NFT monitors, Backend Infra, Web3 utils and Snipers/Minters. **I am frequently updating this repo & it is WIP**.

Should you want to reach out, please do so on Discord at **weeaa**. 🤙🏻

![image](https://github.com/weeaa/web3/assets/108926252/e03cf484-d00c-48df-9665-e75b6a4c94b9)

## 🐰 Features

- Discord Bot with Slash & Buttons features
- Postgres Database with CRUD API
- Friend Tech
    - [x] Indexer
    - [x] Buy/Sells w/ filters
    - [x] New Users w/ filters
    - [x] Deposits w/ filters
    - [x] Pending Deposits w/ filters
    - [x] Invite Redeemer
    - [x] Watchlist Adder
    - [x] [Sniper](https://www.friend.tech/rooms/0xe5d60f8324d472e10c4bf274dbb7371aa93034a0)
- Stars Arena
    - [ ] Monitors
    - [ ] Sniper
- DeFi
    - Uniswap
        - [ ] V2 Swap
        - [ ] V3 Swap
        - [ ] Pair Audit
        - [ ] Utils
     - Raydium
         - [ ] Swap 
- Etherscan Monitoring
  - [x] New Verified Contracts
- ExchangeArt Monitoring
  - [ ] New Drops by Artist (need to update to gql)
- LMNFT Monitoring Top Drops
  - [x] Solana
  - [x] Polygon
  - [x] Ethereum
  - [x] Binance
  - [x] Aptos
  - [x] Avalanche
  - [x] Fantom
  - [x] Stacks
- OpenSea Monitoring
  - [x] Sales
  - [x] Listings
- Premint Monitoring
  - [ ] Hype Weekly/Daily Raffles (Premint NFT Required) – (needs fixes)
- Bitcoin
  - [x] Unisat BRC20 Hype Mint Monitor (Discontinued due to header encryption, cba)
  - [x] Fees Monitor
  - [ ] Unisat BRC20 Minter (thoon)
- Wallet Watchers
    - [ ] Ethereum (thoon)
    - [ ] Base (thoon)
    - [ ] Solana
    - [ ] Polygon
    - [ ] Bitcoin (thoon)
- [x] Twitter Scraper

## 👀 Demo
Below is a demo of our Friend Tech monitor running for Machi, where we had a large pool of new users who hadn't deposited ETH at that time. It's running on localhost, hence the latency. (It is normal for the response status to be 404)

https://github.com/weeaa/web3/assets/108926252/3e997152-29af-4bfb-93db-ee217a22180b

## ⚒️ Project Setup

### Environment

Here is how your `.env` file should be looking like, these values are mainly used for testing purposes.

```ini
NODE_WSS_URL=
NODE_HTTP_URL=
BASIC_USERNAME=
BASIC_PASSWORD=
BOT_TOKEN= <mandatory>
PSQL_PORT=
PSQL_USERNAME=
PSQL_PASSWORD=
PSQL_DB_NAME=
FT_BEARER_TOKEN=
```

Within the scripts directory, you will find a db.sh Bash script that, upon request, generates a database and a table. Please refer to the instructions provided below.

unix 
```bash
$ chmod +x ./scripts/run.sh
$ ./scripts/run.sh
```

windows ⊞
```bat
soon
```


## 🫶🏻 Tips
- Please be aware that for new users of Friend Tech, the use of proxies is essential. Friend Tech tends to impose temporary bans on the same IP address after approximately 90 requests. In my current configuration, I have 1K ISP proxies in place, with bans typically resolved within one second, as demonstrated in the video above.

- You need WSS & HTTP RPCs (commonly named nodes) to monitor on-chain, free ones work well but are subject to rate limits.
    - [Base RPCs](https://docs.base.org/tools/node-providers/)

### Examples

There are various examples which can be found in the [/examples](https://github.com/weeaa/web3/tree/main/examples) folder. In order to build a binary, [Go](https://go.dev/doc/install) 1.20 or higher is required.

After Go is installed, `git clone` the repository and `cd` in `examples/~` (wherever you want) and execute `go build yourfilename.go`.

## Credits

s/o **DALL-E** for the image 😸
