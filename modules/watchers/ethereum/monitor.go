package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/utils"
)

func NewClient(discordClient *discord.Client, verbose bool, nodeUrl string, monitorParams MonitorParams) (*Settings, error) {
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to node: %w", err)
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

func (s *Settings) StartMonitor(wallets []string) {
	logger.LogStartup(moduleName)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogInfo(moduleName, fmt.Sprintf("program panicked! [%v]", r))
				s.StartMonitor(wallets)
				return
			}
		}()
		if err := s.monitorWallets(wallets); err != nil {
			logger.LogError(moduleName, err)
		}
		logger.LogShutDown(moduleName)
	}()
}

func (s *Settings) monitorWallets(wallets []string) error {
	ch := make(chan types.Log)
	sub, err := utils.CreateSubscription(s.Client, utils.SliceToAddresses(wallets), ch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-s.Context.Done():
			return nil
		case <-sub.Err():
			return err
		case log := <-ch:
			if err = s.handleTxnInfo(log); err != nil {

			}
		}
	}
}

func (s *Settings) monitorBalance(wallets []string) {
	go func() {}()
}

// todo finish handling & rename func mby
func (s *Settings) handleTxnInfo(log types.Log) error {

	return nil
}
