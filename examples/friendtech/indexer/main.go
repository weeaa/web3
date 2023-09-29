package main

import (
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/db"
	"github.com/weeaa/nft/modules/friendtech"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	indexer, err := friendtech.NewIndexer(pg, "proxies.txt")
	if err != nil {
		log.Fatal(err)
	}

	indexer.Index()
}
