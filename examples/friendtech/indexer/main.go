package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/modules/friendtech/indexer"
	"log"
	"os"
	"time"
)

// ⚠️⚠️⚠️ proxies are mandatory as you'll get rate limited.
func main() {
	c := make(chan struct{})
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("PSQL_USERNAME"), os.Getenv("PSQL_PASSWORD"), os.Getenv("PSQL_PORT"), os.Getenv("PSQL_DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}

	index, err := indexer.New(pg, "proxies.txt", true, 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	index.StartIndexer()
	<-c
}
