package base

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
)

func NewClient(discordClient *discord.Client, verbose bool, nodeUrl string, monitorParams MonitorParams) (*Settings, error) {

	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to node > %s: %w", nodeUrl, err)
	}

	return &Settings{
		Discord:       discordClient,
		Handler:       handler.New(),
		Verbose:       verbose,
		Context:       context.Background(),
		Client:        client,
		MonitorParams: monitorParams,
	}, nil
}

func (s *Settings) Run(address common.Address) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.Run(address)
			}
		}()
		for !s.monitorBalance(address) {
			select {}
		}
	}()
}

func (s *Settings) monitorWallets() {}

func (s *Settings) monitorBalance(address common.Address) bool {
	balance, err := s.Client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return false
	}
	_ = balance
	return false
}
