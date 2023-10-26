package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/friendtech/watcher"
	"log"
	"os"
)

func main() {
	c := make(chan struct{})
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	discBot, err := bot.New(pg)
	if err != nil {
		log.Fatal(err)
	}

	if err = discBot.Start(); err != nil {
		log.Fatal(err)
	}

	friendTechWatcher, err := watcher.NewFriendTech(pg, discBot, "proxies.txt", os.Getenv("NODE_WSS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	friendTechWatcher.StartAllWatchers(1)
	<-c
}
