package ethereum

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/utils"
	"github.com/weeaa/nft/pkg/utils/ethereum"
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

func (s *Settings) StartMonitor() {
	logger.LogStartup(moduleName)
	go func() {

		defer func() {
			if r := recover(); r != nil {
				logger.LogInfo(moduleName, fmt.Sprintf("program panicked! [%v]", r))
				s.StartMonitor()
				return
			}
		}()

		go s.Transactions.ParseTransactions()
		s.Transactions.GetLatestBlock(0)
	}()
}

func (s *Settings) monitor() {

}

func (s *Settings) WatchPendingTransactions(address common.Address, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}

func getPendingTxns() {}

// monitorBalance monitors any balance changes on a specific address.
func (s *Settings) monitorBalance(wallet common.Address, retryDelay time.Duration, ctx context.Context) {
	s.PromMetrics.GoroutineCount.Inc()
	go func(address common.Address, delay time.Duration, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				balanceNow, err := ethereum.GetEthWalletBalance(s.Client, address)
				if err != nil {
					logger.LogError(moduleName, err)
					continue
				}

				balanceBefore, ok := s.BalanceWatcher.Balances[address]
				if !ok {
					s.BalanceWatcher.Balances[address] = balanceNow
					continue
				}

				if balanceBefore != balanceNow {
					s.BalanceWatcher.Balances[address] = balanceNow

					var status string
					var user *models.FriendTechMonitor
					_ = user

					if balanceNow.Int64() > balanceBefore.Int64() {
						status = "Increased ↖︎"
					} else {
						status = "Decreased ↘︎"
					}

					ethBalanceBefore := ethereum.WeiToEther(balanceBefore)
					ethBalanceAfter := ethereum.WeiToEther(balanceNow)

					user, err = s.DB.Monitor.GetUserByAddress(address.String(), context.Background())
					if err != nil {
						// check if it's stored
						switch err.Error() {
						case "":
						case "d":
						}

					}

					s.Bot.BotWebhook(&discordgo.MessageSend{
						Components: bot.BundleQuickTaskComponents("", ""),
						Embeds: []*discordgo.MessageEmbed{
							{
								Color:       bot.Purple,
								Title:       fmt.Sprintf("%s Balance's %s", utils.FirstLastFour(wallet.String()), status),
								Description: fmt.Sprintf(""),
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:   "Balance Status",
										Value:  status,
										Inline: true,
									},
									{
										Name:   "Balance After",
										Value:  fmt.Sprintf("%3.f", ethBalanceAfter),
										Inline: true,
									},
									{
										Name:   "Balance Before",
										Value:  fmt.Sprintf("%3.f", ethBalanceBefore),
										Inline: true,
									},
									{},
								},
							},
						},
					}, "")
				}
				time.Sleep(retryDelay)
			}
		}
	}(wallet, retryDelay, ctx)
}

// todo finish handling & rename func mby
func (s *Settings) handleTxnInfo(log types.Log) error {

	return nil
}

func (t *Transactions) GetLatestBlock(delay time.Duration) {
	for {
		latestBlock, err := t.Client.BlockByNumber(context.Background(), nil)
		if err != nil {
			continue
		}

		if latestBlock.NumberU64() == t.LatestBlock {
			continue
		}

		t.LatestBlock = latestBlock.NumberU64()
		t.LatestBlockChan <- latestBlock

		time.Sleep(delay)
	}
}

func (t *Transactions) ParseTransactions() {
	for {
		block := <-t.LatestBlockChan
		go func(latestBlock *types.Block) {
			for _, tx := range latestBlock.Transactions() {
				go t.handleTransaction(tx)
			}
		}(block)
	}
}

func (t *Transactions) handleTransaction(tx *types.Transaction) {}

func (t *Transactions) isGreater(input any) {

}

func (s *Settings) SetParam(address common.Address, txnParams *TransactionsParams) {
	s.Transactions.Params.Set(address, txnParams)
}
