package solana

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/utils"
)

func NewClient(discord *discord.Client, nodeUrl string) (*Settings, error) {
	client, err := ws.Connect(context.Background(), nodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to node > %s : %w", nodeUrl, err)
	}
	return &Settings{
		Client: client,
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
		ch := make(chan error)
		go func() {
			if err := s.monitorWallets(wallets, ch); err != nil {
				logger.LogError(moduleName, fmt.Errorf("%w", err))
			}
		}()
		// trying smth, may not work tho
		for err := range ch {
			logger.LogError(moduleName, fmt.Errorf("%w", err))
		}
		logger.LogShutDown(moduleName)
	}()
}

func (s *Settings) monitorWallets(wallets []string, ch chan error) error {
	programs := utils.SliceToPrograms(wallets)

	for _, program := range programs {
		go func(address solana.PublicKey) {
			var log *ws.LogResult
			sub, err := s.Client.LogsSubscribeMentions(
				address,
				rpc.CommitmentRecent,
			)
			if err != nil {
				ch <- err
			}
			defer sub.Unsubscribe()

			for {
				select {
				case <-s.Context.Done():
					ch <- nil
				default:
					log, err = sub.Recv()
					if err != nil {
						ch <- err
					}
				}
			}
		}(program)
	}
	return nil
}
