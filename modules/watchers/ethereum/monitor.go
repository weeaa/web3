package ethereum

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/utils"
	"time"
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

// doesnt work
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

// monitorBalance monitors any balance changes on a specific address.
func (s *Settings) monitorBalance(wallet common.Address) {
	go func(address common.Address) {
		for {
			balance, err := utils.GetEthWalletBalance(s.Client, address)
			if err != nil {
				logger.LogError(moduleName, err)
				continue
			}

			b, ok := s.BalanceWatcher.Balances[address]
			if !ok {
				s.BalanceWatcher.Balances[address] = balance
				continue
			}

			if b != balance {
				s.BalanceWatcher.Balances[address] = balance

				var status string
				var userInfo database.User
				_ = userInfo
				if balance.Int64() > b.Int64() {
					status = "Increased"
				} else {
					status = "Decreased"
				}

				ethBalanceBefore := utils.WeiToEther(b)
				ethBalanceAfter := utils.WeiToEther(balance)

				/*
					userInfo, err = s.DB.GetUser(address.String())
					if err != nil {
						logger.LogError(moduleName, err)
					}
				*/

				s.Bot.BotWebhook(&discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: fmt.Sprintf("%s | Balance %s", address.String(), status),

							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Address",
									Value:  address.String(),
									Inline: true,
								},
								{
									Name:   "Ξ Balance BF/AF",
									Value:  fmt.Sprintf("%2.f | %2.f", ethBalanceBefore, ethBalanceAfter),
									Inline: true,
								},
								{
									Name:   "QuickLink",
									Value:  bot.BundleQuickLinks(address.String()),
									Inline: false,
								},
							},
						},
					},
				}, nil, bot.BalanceChange)
			}

			time.Sleep(2500 * time.Millisecond)
		}
	}(wallet)
}

// todo finish handling & rename func mby
func (s *Settings) handleTxnInfo(log types.Log) error {

	return nil
}

// OK
func (t *Transactions) GetLatestBlock() error {
	for {
		latestBlock, err := t.Client.BlockByNumber(context.Background(), nil)
		if err != nil {
			return err
		}

		if latestBlock.NumberU64() == t.LatestBlock {
			continue
		}

		t.LatestBlock = latestBlock.NumberU64()
		t.LatestBlockChan <- latestBlock
	}
}

// OK todo finish func
func (t *Transactions) ParseTransactions() {
	block := <-t.LatestBlockChan
	go func(latestBlock *types.Block) {
		for _, txns := range latestBlock.Transactions() {

		}
	}(block)
}
