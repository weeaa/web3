package main

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/weeaa/nft/modules/friendtech/sniper"
	"github.com/weeaa/nft/pkg/files"
	"github.com/weeaa/nft/pkg/utils"
	"log"
	"math/big"
	"os"
)

// sike, will be public when ft dies

// well...
func main() {

	maxEthInput, ok := new(big.Float).SetString("1")
	if !ok {
		log.Fatal("error setting maxEthInput")
	}

	sniperClient, err := sniper.New(os.Getenv("FT_PRIVATE_KEY"), os.Getenv("NODE_HTTP_URL"), maxEthInput, 2)
	if err != nil {
		log.Fatal(err)
	}

	clients := make(map[string]Task)

	// load tasks from csv file
	tasks, err := loadTasks()
	if err != nil {
		log.Fatal(err)
	}

	// generate clients
	for _, task := range tasks {
		clients[task.Wallet.PublicKey.String()] = task
	}

	// todo finish example

	snipe, err := sniperClient.Snipe(common.HexToAddress("0xe5d60f8324D472E10C4BF274dBb7371aa93034A0"), "1", sniper.Buy, 1, sniper.Normal)
	if err != nil {
		log.Fatal(err)
	}

	snipe.Txns.ForEach(func(transaction *types.Transaction, err error) {
		log.Println(transaction, "err=", err)
	})

	log.Println("exiting")
}

type Task struct {
	PrivateKey  string     `csv:"Private Key"`
	MaxEthInput *big.Float `csv:"Max ETH Input"`
	MaxBuys     uint       `csv:"Max Buys"`

	Wallet *utils.EthereumWallet
}

func loadTasks() ([]Task, error) {
	path := "tasks.csv"
	if _, err := os.Stat(path); errors.Is(err, os.ErrExist) {
		files.CreateCSV(path, [][]string{{
			"Private Key",
			"Max ETH Input",
			"Max Buys",
		}})
		return nil, fmt.Errorf("created file")
	}
	return files.ReadCSV[Task](path)
}
