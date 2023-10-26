package main

import (
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/hub3/watcher"
	"log"
)

func main() {
	c := make(chan struct{})
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	discBot, err := bot.New(pg)
	if err != nil {
		log.Fatal(err)
	}

	discBot.Start()
	defer discBot.Stop()

	wa, err := watcher.New(discBot, "proxies.txt")
	usernames := []string{
		"OttoSuwenNFT",
		"HsakaTrades",
		"fewture",
		"saliencexbt",
		"clamggz",
		"const_phoenixed",
		"Anonymoux2311",
		"0xMakesy",
		"mooncat2878",
	}

	go func() {
		for _, user := range usernames {
			wa.WatchUser(user)
		}
	}()

	wa.Map()

	<-c

}
