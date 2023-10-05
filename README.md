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
    - [x] [Sniper](https://www.friend.tech/rooms/0xe5d60f8324d472e10c4bf274dbb7371aa93034a0) (Access for FriendTech Holders only, gotta keep it competitive 🫶🏻)
- Stars Arena
    - [ ] thoon 🤓
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
  - [x] Unisat BRC20 Hype Mint Monitor
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

Here is how your `.env` file should be looking like, you can omit `NODE_WSS_URL` & `NODE_HTTP_URL`.

```ini
# Not Mandatory
NODE_WSS_URL=
NODE_HTTP_URL=
# API (Mandatory)
BASIC_USERNAME=email@gmail.com
BASIC_PASSWORD=nYp@SsWW0rD
# END
BOT_TOKEN=discordBotToken
PSQL_PORT=
PSQL_USERNAME=
PSQL_PASSWORD=
PSQL_DB_NAME=
```

In the `scripts` folder, you may find a `db.sh` bash script which will create a table and a database upon request. Follow instructions below.

```bash
$ chmod +x ./scripts/db.sh
$ ./scripts/db.sh
```

## 🫶🏻 Tips
Note that proxies are mandatory for Friend Tech New Users (off-chain stuff) as they ban you on average at the ~90th request you do on the same IP – they only ban temporarily tho. My current setup is 1k ISP and it runs perfectly, with bans being resolved in 1s as you can see on the demo. Residential proxies will be costly, I advise to have a pool of DCs or ISPs. I may add delay in future updates so run it at the pace you want.

You need WSS & HTTP (commonly named nodes) to monitor on-chain, free ones work well.

- [Base RPCs](https://docs.base.org/tools/node-providers/)

It is also advised to run on a server if you run the Friend Tech New Users Monitor as it gets network intensive sometimes.

### Examples

There are various examples which can be found in the [/examples](https://github.com/weeaa/web3/tree/main/examples) folder.

## Credits

s/o **DALL-E** for the image 😸
