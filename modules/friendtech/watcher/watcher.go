package watcher

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/friendtech"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/modules/twitter"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func NewFriendTech(bot *bot.Bot, proxyFilePath, nodeURL string, counter int) (*Watcher, error) {
	wssClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		Bot:           bot,
		Addresses:     make(map[string]string),
		Client:        tls.NewProxyLess(),
		WSSClient:     wssClient,
		NewUsersCtx:   context.Background(),
		PendingDepCtx: context.Background(),
		ProxyList:     proxyList,
		Counter:       counter,
	}, nil
}

func (w *Watcher) StartAllWatchers() {
	go func() {
		defer func() {
			// write to json the latest user id taken
		}()
		for !w.WatchNewUsers() {
			select {
			case <-w.NewUsersCtx.Done():
				return
			default:

			}
		}
	}()
}

func (w *Watcher) WatchNewUsers() bool {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: fren_utils.ProdBaseApi, Path: "/users/by-id/" + fmt.Sprint(w.Counter)},
		Host:   fren_utils.ProdBaseApi,
		Header: http.Header{},
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			logger.LogError(watcher, fmt.Errorf("status not found for id: %d", w.Counter))
			time.Sleep(1 * time.Second)
			return false
		} else if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
			time.Sleep(1 * time.Second)
			tls.HandleRateLimit(w.Client, w.ProxyList, watcher)
			return false
		}
		logger.LogError(watcher, fmt.Errorf("status %s for id: %d", resp.Status, w.Counter))
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	var u friendtech.UserInformation
	if err = json.Unmarshal(body, &u); err != nil {
		logger.LogError(watcher, err)
		return false
	}

	uInfo, err := fren_utils.GetUserInformation(u.Address, w.Client)
	if err != nil {
		w.Pool <- u.Address
		w.Counter++
		return false
	}

	var nitter twitter.NitterResponse
	var importance fren_utils.Importance
	var displayedPrice *big.Float
	{
		nitter, err = twitter.FetchNitter(u.TwitterUsername, w.Client)
		if err != nil {
			logger.LogError(watcher, err)
			return false
		}

		followers, _ := strconv.Atoi(nitter.Followers)

		importance = fren_utils.AssertImportance(followers, fren_utils.Followers)
		if err != nil {
			logger.LogError(watcher, err)
		}
	}

	balance, err := utils.GetEthWalletBalance(w.HTTPClient, common.HexToAddress(uInfo.Address))
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	wei := new(big.Int)
	wei.SetString(uInfo.DisplayPrice, 10)
	displayedPrice = utils.WeiToEther(wei)

	balanceEth := utils.WeiToEther(balance)

	var channelID string
	var roleID string
	switch importance {
	case fren_utils.Shrimp:
		roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.ShrimpEmote])
		channelID = bot.FriendTechNewUsers5
	case fren_utils.Whale:
		roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.WhaleEmote])
		channelID = bot.FriendTechNewUsers50
	case fren_utils.Fish:
		roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.FishEmote])
		channelID = bot.FriendTechNewUsers10
	case "none":
		channelID = bot.FriendTechNewUsers
	}

	w.Bot.BotWebhook(&discordgo.MessageSend{
		Content: roleID,
		Embeds: []*discordgo.MessageEmbed{
			{
				Description: fmt.Sprintf("[%s](https://basescan.org/address/%s)", u.Address, u.Address),
				Color:       bot.Purple,
				Title:       u.TwitterUsername,
				URL:         fmt.Sprintf("https://www.friend.tech/rooms/%s", u.Address),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: u.TwitterPfpUrl,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("friend.tech – new users [%s]", roleID),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Price",
						Value:  fmt.Sprintf("%.4f", displayedPrice) + "Ξ",
						Inline: true,
					},
					{
						Name:   "Holders | Followers",
						Value:  fmt.Sprintf("%d | %s", uInfo.HolderCount, nitter.Followers),
						Inline: true,
					},
					{
						Name:   "Balance",
						Value:  fmt.Sprintf("%.3f", balanceEth) + "Ξ",
						Inline: true,
					},
					{
						Name:   "Age",
						Value:  nitter.AccountAge,
						Inline: true,
					},
					{
						Name:   "Twitter Name",
						Value:  u.TwitterName,
						Inline: true,
					},
					{
						Name:   "Twitter Username",
						Value:  fmt.Sprintf("[%s](https://x.com/%s)", u.TwitterUsername, u.TwitterUsername),
						Inline: true,
					},
					{
						Name:   "QuickTask",
						Value:  fmt.Sprintf(""),
						Inline: true,
					},
				},
			},
		},
	}, channelID)
	logger.LogInfo(watcher, fmt.Sprintf("%d | %s", w.Counter, u.TwitterUsername))

	w.Counter++
	return false
}

func (w *Watcher) WatchNewUsersPool(deadline float64) {
	for {
		user := <-w.Pool
		startTime := time.Now()
		go func(address string) {
			for {
				if time.Since(startTime).Hours() >= deadline {
					break
				}

				userInfo, err := fren_utils.GetUserInformation(address, w.Client)
				if err != nil {
					logger.LogError(watcher, err)
					continue
				}

				var nitter twitter.NitterResponse
				var importance fren_utils.Importance
				var displayedPrice *big.Float

				{
					nitter, err = twitter.FetchNitter(userInfo.TwitterUsername, w.Client)
					if err != nil {
						logger.LogError(watcher, err)
						continue
					}

					followers, _ := strconv.Atoi(nitter.Followers)

					importance = fren_utils.AssertImportance(followers, fren_utils.Followers)
					if err != nil {
						logger.LogError(watcher, err)
						continue
					}

				}

				balance, err := utils.GetEthWalletBalance(w.HTTPClient, common.HexToAddress(userInfo.Address))
				if err != nil {
					logger.LogError(watcher, err)
					continue
				}

				wei := new(big.Int)
				wei.SetString(userInfo.DisplayPrice, 10)
				displayedPrice = utils.WeiToEther(wei)

				balanceEth := utils.WeiToEther(balance)

				var channelID string
				var roleID string
				switch importance {
				case fren_utils.Shrimp:
					roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.ShrimpEmote])
					channelID = bot.FriendTechNewUsers5
				case fren_utils.Fish:
					roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.FishEmote])
					channelID = bot.FriendTechNewUsers10
				case fren_utils.Whale:
					roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.WhaleEmote])
					channelID = bot.FriendTechNewUsers50
				case "none":
					channelID = bot.FriendTechNewUsers
				}

				w.OutStreamData <- BroadcastData{Event: NewSignup, Data: map[string]any{
					"user":   userInfo,
					"nitter": nitter,
				}}

				w.Bot.BotWebhook(&discordgo.MessageSend{
					Content: roleID,
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: fmt.Sprintf("[%s](https://basescan.org/address/%s)", userInfo.Address, userInfo.Address),
							Color:       bot.Purple,
							Title:       userInfo.TwitterUsername,
							URL:         fmt.Sprintf("https://www.friend.tech/rooms/%s", userInfo.Address),
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: userInfo.TwitterPfpUrl,
							},
							Footer: &discordgo.MessageEmbedFooter{
								Text: "friend.tech – new users " + roleID,
							},
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Price",
									Value:  fmt.Sprintf("%.4f", displayedPrice) + "Ξ",
									Inline: true,
								},
								{
									Name:   "Holders | Followers",
									Value:  fmt.Sprintf("%d | %s", userInfo.HolderCount, nitter.Followers),
									Inline: true,
								},
								{
									Name:   "Balance",
									Value:  fmt.Sprintf("%.3f", balanceEth) + "Ξ",
									Inline: true,
								},
								{
									Name:   "Age",
									Value:  nitter.AccountAge,
									Inline: true,
								},
								{
									Name:   "Twitter Name",
									Value:  userInfo.TwitterName,
									Inline: true,
								},
								{
									Name:   "Twitter Username",
									Value:  fmt.Sprintf("[%s](https://x.com/%s)", userInfo.TwitterUsername, userInfo.TwitterUsername),
									Inline: true,
								},
								{
									Name:   "QuickTask",
									Value:  fmt.Sprintf(""),
									Inline: true,
								},
							},
						},
					},
				}, channelID)
				logger.LogInfo(watcher, fmt.Sprintf("%d | %s", w.Counter, userInfo.TwitterUsername))
			}
		}(user)
	}
}

// WatchAddNewUsers adds new users fetched from the db
func (w *Watcher) WatchAddNewUsers() {
	for {
		addresses, err := w.DB.Monitor.GetAllAddresses(context.Background())
		if err != nil {

		}

		for _, address := range addresses {
			if _, ok := w.Addresses[address]; !ok {
				w.Addresses[address] = uuid.NewString()
			}
		}

		time.Sleep(10 * time.Second)
	}
}

// WatchPendingDeposits self explicit bro
func (w *Watcher) WatchPendingDeposits() bool {
	resp, err := w.doRequestPendingTxn()
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	var r L2TransactionsPendingResponse
	if err = json.Unmarshal(body, &r); err != nil {
		return false
	}

	for _, item := range r.Items {
		for _, address := range w.Addresses {
			//if address == item.
			// extract tx data

		}
	}

	//curl 'https://eth.blockscout.com/api/v2/transactions/0xf1929a1096f1cf94b2e3b1ee79a70e77ce4c607edb77c934d305335dc827a0ef' \
	//  -H 'authority: eth.blockscout.com' \
	//  -H 'accept: */*' \
	//  -H 'accept-language: en-US,en;q=0.9' \
	//  -H 'dnt: 1' \
	//  -H 'referer: https://eth.blockscout.com/tx/0xf1929a1096f1cf94b2e3b1ee79a70e77ce4c607edb77c934d305335dc827a0ef' \
	//  -H 'sec-ch-ua: "Chromium";v="117", "Not;A=Brand";v="8"' \
	//  -H 'sec-ch-ua-mobile: ?0' \
	//  -H 'sec-ch-ua-platform: "macOS"' \
	//  -H 'sec-fetch-dest: empty' \
	//  -H 'sec-fetch-mode: cors' \
	//  -H 'sec-fetch-site: same-origin' \
	//  -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
	//  --compressed

	return false
}

func (w *Watcher) doRequestPendingTxn() (*http.Response, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "optimism.blockscout.com", Path: "/api/v2/optimism/deposits"},
		Header: http.Header{
			"authority":          {"optimism.blockscout.com"},
			"accept":             {"*/*"},
			"accept-language":    {"en-US,en;q=0.9"},
			"dnt":                {"1"},
			"referer":            {"https://optimism.blockscout.com/l2-deposits"},
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-origin"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
		},
	}

	return w.Client.Do(req)
}

func (w *Watcher) SubscribeFriendTech() {
	ch := make(chan types.Log)
	sub, err := utils.CreateSubscription(w.WSSClient, []common.Address{common.HexToAddress(fren_utils.FRIEND_TECH_CONTRACT_V1)}, ch)
	if err != nil {
		logger.LogFatal(watcher, err)
	}

	for {
		select {
		case err = <-sub.Err():
			logger.LogError(watcher, fmt.Errorf("subscription stopped: %w", err))
			return
		case tx := <-ch:
			go w.dispatchLog(tx)
		}
	}
}

func (w *Watcher) dispatchLog(txn types.Log) {
	var tx *types.Transaction
	var sender common.Address
	var err error

	tx, _, err = w.HTTPClient.TransactionByHash(context.Background(), txn.TxHash)
	if err != nil {
		logger.LogError(watcher, err)
	}

	sender, err = utils.GetSender(tx)
	if err != nil {
		logger.LogError(watcher, err)
	}

	if strings.Contains(tx.Hash().Hex(), sellMethod) {
		if err = w.handleSell(tx, sender.String()); err != nil {
			logger.LogError(watcher, fmt.Errorf("handleSell: %w", err))
		}
	} else if strings.Contains(tx.Hash().Hex(), buyMethod) {
		if err = w.handleBuy(tx, sender.String()); err != nil {
			logger.LogError(watcher, fmt.Errorf("handleBuy: %w", err))
		}
	} else {
		logger.LogError(watcher, fmt.Errorf("unknown/unsupported tx type: %s", tx.Hash().String()))
	}
}

func (w *Watcher) handleSell(tx *types.Transaction, sender string) error {
	var channel string
	var err error
	var isBot bool

	data, err := utils.DecodeTransactionInputData(w.ABI, string(tx.Data()))
	if err != nil {
		return err
	}

	recipient := data["sharesSubject"].(string)
	sharesAmount := data["amount"].(uint256.Int)
	if sharesAmount.Uint64() > 1 {
		isBot = true
	}

	senderInfo, err := fren_utils.GetUserInformation(sender, w.Client)
	if err != nil {
		return err
	}

	sharePrice := new(big.Int)
	sharePrice.SetString(senderInfo.DisplayPrice, 10)

	recipientInfo, err := fren_utils.GetUserInformation(recipient, w.Client)
	if err != nil {
		return err
	}

	balance, err := utils.GetEthWalletBalance(w.HTTPClient, common.HexToAddress(sender))
	if err != nil {
		return err
	}

	if isSelf(hex.EncodeToString(tx.Data()), sender) {
		// goes to selfbuy chan
		if _, ok := w.Addresses[sender]; !ok { // means it's not on our list & doesn't go to filtered

			var user *models.FriendTechMonitorAll

			user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
			if err != nil { // means we can't find it & not stored in our db, we fetch the data
				var nitter twitter.NitterResponse

				nitter, err = twitter.FetchNitter(senderInfo.TwitterUsername, w.Client)
				if err != nil {
					return err
				}

				status := fren_utils.AssertImportance(nitter.Followers, fren_utils.Followers)

				if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
					BaseAddress:     sender,
					Status:          string(status),
					Followers:       nitter.Followers,
					TwitterUsername: senderInfo.TwitterUsername,
					TwitterName:     senderInfo.TwitterName,
					TwitterURL:      senderInfo.TwitterPfpUrl,
				}, context.Background()); err != nil {
					return err
				}

				user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
				if err != nil {
					return err
				}
			}

			switch user.Status {
			case "high":
				channel = bot.FriendTechWhalesSells
			case "medium":
				channel = bot.FriendTechFishSells
			case "low":
				channel = bot.FriendTechShrimpSells
			case "none":
				channel = bot.FriendTechAllLogs
			}

			go w.Bot.BotWebhook(&discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s bought %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – unfiltered sells",
							IconURL: bot.Image,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Buyer Balance",
								Value: fmt.Sprintf("%2.f Ξ", utils.WeiToEther(balance)),
							},
							{
								Name:   "Share Amount/Price",
								Value:  fmt.Sprintf("%v | %s", sharesAmount, utils.WeiToEther(sharePrice)),
								Inline: true,
							},
							{
								Name:   "IsBotPurchase",
								Value:  fmt.Sprint(isBot),
								Inline: true,
							},
						},
					},
				},
			}, channel)
		} else {
			// is on our monitored list
			var user *models.FriendTechMonitor

			user, err = w.DB.Monitor.GetUserByAddress(sender, context.Background())
			if err != nil {
				return err
			}

			go w.Bot.BotWebhook(&discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s bought %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – filtered sells",
							IconURL: bot.Image,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "Buyer Balance",
								Value: fmt.Sprintf("%2.f Ξ", utils.WeiToEther(balance)),
							},
							{
								Name:   "Share Amount/Price",
								Value:  fmt.Sprintf("%v | %s", sharesAmount, utils.WeiToEther(sharePrice)),
								Inline: true,
							},
							{
								Name:   "IsBotPurchase",
								Value:  fmt.Sprint(isBot),
								Inline: true,
							},
							{
								Name:   "User Importance",
								Value:  user.Status,
								Inline: true,
							},
						},
					},
				},
			}, bot.FriendTechFilteredSells)
		}
	} else { // is not a self buy
		var user *models.FriendTechMonitorAll

		user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
		if err != nil { // means we can't find it & not stored in our db, we fetch the data
			var nitter twitter.NitterResponse

			nitter, err = twitter.FetchNitter(senderInfo.TwitterUsername, w.Client)
			if err != nil {
				return err
			}

			status := fren_utils.AssertImportance(nitter.Followers, fren_utils.Followers)

			if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
				BaseAddress:     sender,
				Status:          string(status),
				Followers:       nitter.Followers,
				TwitterUsername: senderInfo.TwitterUsername,
				TwitterName:     senderInfo.TwitterName,
				TwitterURL:      senderInfo.TwitterPfpUrl,
			}, context.Background()); err != nil {
				return err
			}

			user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
			if err != nil {
				return err
			}
		}

		switch user.Status {
		case "high":
			channel = bot.FriendTechWhalesSells
		case "medium":
			channel = bot.FriendTechFishSells
		case "low":
			channel = bot.FriendTechShrimpSells
		case "none":
			channel = bot.FriendTechAllLogs
		}

		go w.Bot.BotWebhook(&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       bot.Purple,
					Title:       fmt.Sprintf("%s bought %v key(s) of %s", senderInfo.TwitterUsername, sharesAmount, recipientInfo.TwitterUsername),
					Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "friendtech – unfiltered sells",
						IconURL: bot.Image,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Buyer Balance",
							Value: fmt.Sprintf("%2.f Ξ", utils.WeiToEther(balance)),
						},
						{
							Name:   "Share Amount/Price",
							Value:  fmt.Sprintf("%v | %s", sharesAmount, utils.WeiToEther(sharePrice)),
							Inline: true,
						},
						{
							Name:   "IsBotPurchase",
							Value:  fmt.Sprint(isBot),
							Inline: true,
						},
					},
				},
			},
		}, bot.FriendTechAllLogs)
	}
	return nil
}

func (w *Watcher) handleBuy(tx *types.Transaction, sender string) error {
	if isSelf(hex.EncodeToString(tx.Data()), sender) {
		return nil
	} else {
		return nil
	}
}

func isSelf(txData string, sender string) bool {
	return strings.Contains(strings.ToLower(txData), strings.ToLower(sender[2:]))
}
