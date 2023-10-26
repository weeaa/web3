package watcher

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/modules/friendtech/constants"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/modules/twitter"
	"github.com/weeaa/nft/pkg/cache"
	"github.com/weeaa/nft/pkg/files"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewFriendTech(db *db.DB, bot *bot.Bot, proxyFilePath, nodeURL string) (*Watcher, error) {
	wssClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	// switch to file read
	friendTechABI, _ := abi.JSON(strings.NewReader("[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"trader\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"subject\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isBuy\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shareAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolEthAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subjectEthAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"}],\"name\":\"Trade\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"buyShares\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBuyPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getBuyPriceAfterFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"supply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getSellPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"getSellPriceAfterFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeeDestination\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeePercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sharesSubject\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"sellShares\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeDestination\",\"type\":\"address\"}],\"name\":\"setFeeDestination\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feePercent\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeePercent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_feePercent\",\"type\":\"uint256\"}],\"name\":\"setSubjectFeePercent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sharesBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sharesSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"subjectFeePercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"))

	return &Watcher{
		DB:            db,
		ABI:           friendTechABI,
		Bot:           bot,
		Addresses:     make(map[string]string),
		WatcherClient: tls.New(tls.RandProxyFromList(proxyList)),
		NitterClient:  twitter.NewClient("", "", proxyList),
		OutStreamData: make(chan BroadcastData),
		WSSClient:     wssClient,
		NewUsersCtx:   context.Background(),
		PendingDepCtx: context.Background(),
		Pool:          make(chan string),
		ProxyList:     proxyList,
		//Counter:       counter,
	}, nil
}

func (w *Watcher) StartAllWatchers(counter int) {
	w.Counter = counter
	log.Info().Str("mod", watcher)
	defer func() {
		if r := recover(); r != nil {
			log.Warn().Str("panic", "recovered")
			w.StartAllWatchers(w.Counter)
		}
	}()

	if w.EnablePool {
		go w.WatchNewUsersPool()
	}
	go w.WatchAddNewUsers()

	go func() {
		go w.WatchNewUsersPool()
		for !w.WatchNewUsers() {
			select {
			case <-w.NewUsersCtx.Done():

				return
			default:
				continue
			}
		}
	}()

	go func() {
		w.SubscribeFriendTech()
	}()
}

// WatchNewUsers monitors new sign-ups, & will ping even if
// the user has not deposited any amount of ETH. You can change
// this setting by enabling the 'Pool'.
func (w *Watcher) WatchNewUsers() bool {
	logger.LogInfo(watcher, "checking id "+fmt.Sprint(w.Counter))

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: constants.ProdBaseApi, Path: "/users/by-id/" + fmt.Sprint(w.Counter)},
		Host:   constants.ProdBaseApi,
		Header: http.Header{},
	}

	resp, err := w.WatcherClient.Do(req)
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			time.Sleep(1 * time.Second)
			return false
		} else if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
			tls.HandleRateLimit(w.WatcherClient, w.ProxyList, watcher)
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

	var u fren_utils.UserInformation
	if err = json.Unmarshal(body, &u); err != nil {
		logger.LogError(watcher, err)
		return false
	}

	uInfo, err := fren_utils.GetUserInformation(u.Address, w.WatcherClient)
	if err != nil {

		//logger.LogInfo(watcher, fmt.Sprintf("%s | %d signed up but didn't deposit, adding to pool...", u.TwitterUsername, w.Counter))
		//w.Pool <- u.Address
		//w.Counter++
		//return false
	}

	var nitter twitter.NitterResponse
	var importance fren_utils.Importance
	{
		nitter, err = w.NitterClient.FetchNitter(u.TwitterUsername)
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

	balance, err := utils.GetEthWalletBalance(w.WSSClient, common.HexToAddress(uInfo.Address))
	if err != nil {
		logger.LogError(watcher, err)
		return false
	}

	wei, ok := new(big.Int).SetString(uInfo.DisplayPrice, 10)
	if !ok {
		return false
	}

	displayedPrice := utils.WeiToEther(wei)
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

	if err = files.WriteJSON("id.json", map[string]int{"id": u.Id}); err != nil {
		logger.LogError(watcher, err)
	}

	w.Counter++
	return false
}

func (w *Watcher) WatchNewUsersPool() {
	wg := sync.WaitGroup{}
	for {
		user := <-w.Pool

		logger.LogInfo(watcher, "pool add: "+user)
		startTime := time.Now()

		wg.Add(1)
		go func(address string) {
			for {
				if time.Since(startTime).Hours() >= w.Deadline {
					break
				}

				userInfo, err := fren_utils.GetUserInformation(address, w.WatcherClient)
				if err != nil {
					continue
				}

				var nitter twitter.NitterResponse
				var importance fren_utils.Importance
				var displayedPrice *big.Float

				{
					nitter, err = twitter.FetchNitter(userInfo.TwitterUsername, w.WatcherClient)
					if err != nil {
						logger.LogError(watcher, err)
						continue
					}

					followers, _ := strconv.Atoi(nitter.Followers)

					importance = fren_utils.AssertImportance(followers, fren_utils.Followers)
				}

				balance, err := utils.GetEthWalletBalance(w.WSSClient, common.HexToAddress(userInfo.Address))
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

				// tood modify
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
			logger.LogError(watcher, err)
		}

		for _, address := range addresses {
			if _, ok := w.Addresses[address]; !ok {
				w.Addresses[address] = uuid.NewString()
			}
		}

		time.Sleep(10 * time.Second)
	}
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

	return w.WatcherClient.Do(req)
}

func (w *Watcher) SubscribeFriendTech() {
	ch := make(chan types.Log)
	sub, err := utils.CreateSubscription(w.WSSClient, []common.Address{common.HexToAddress(constants.FRIEND_TECH_CONTRACT_V1)}, ch)
	if err != nil {
		logger.LogFatal(watcher, err.Error())
		return
	}

	for {
		select {
		case err = <-sub.Err():
			logger.LogError(watcher, fmt.Errorf("subscription stopped: %w | attempting to restore", err))
			go w.SubscribeFriendTech()
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

	logger.LogInfo(watcher, "tx found: "+txn.TxHash.String())

	tx, _, err = w.WSSClient.TransactionByHash(context.Background(), txn.TxHash)
	if err != nil {
		logger.LogError(watcher, err)
		return
	}

	sender, err = utils.GetSender(tx)
	if err != nil {
		logger.LogError(watcher, err)
		return
	}

	txData := utils.HexEncodeTxData(tx.Data())

	if strings.Contains(txData, sellMethod) {
		if err = w.handleSell(tx, sender.String(), txData); err != nil {
			logger.LogError(watcher, fmt.Errorf("handleSell: %w", err))
		}
	} else if strings.Contains(txData, buyMethod) {
		if err = w.handleBuy(tx, sender.String(), txData); err != nil {
			logger.LogError(watcher, fmt.Errorf("handleBuy: %w", err))
		}
	} else {
		logger.LogError(watcher, fmt.Errorf("unknown/unsupported tx type: %s", txData))
	}
}

func (w *Watcher) handleSell(tx *types.Transaction, sender, txData string) error {
	var channel string
	var err error
	var isBot bool

	data, err := utils.DecodeTransactionInputData(w.ABI, txData)
	if err != nil {
		return err
	}

	recipient := data["sharesSubject"].(common.Address)
	sharesAmount := data["amount"].(*big.Int)
	if sharesAmount.Int64() > 1 {
		isBot = true
	}

	senderInfo, _ := fren_utils.GetUserInformation(sender, w.WatcherClient)

	sharePrice := new(big.Int)
	sharePrice.SetString(senderInfo.DisplayPrice, 10)

	recipientInfo, err := fren_utils.GetUserInformation(recipient.String(), w.WatcherClient)
	if err != nil {
		return err
	}

	balance, err := utils.GetEthWalletBalance(w.WSSClient, common.HexToAddress(sender))
	if err != nil {
		return err
	}

	if isSelf(hex.EncodeToString(tx.Data()), sender) {

		/*
			var item *memcache.Item
			item, err = w.Cache.RetrieveData(sender)
			if err != nil {
				return err
			}

			if err = w.Cache.InsertData(sender, "self"); err != nil {
				return err
			}
		*/

		// goes to selfbuy chan
		if _, ok := w.Addresses[strings.ToLower(sender)]; !ok { // means it's not on our list & doesn't go to filtered
			var user *models.FriendTechMonitorAll

			user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
			if err != nil { // means we can't find it & not stored in our db, we fetch the data
				var nitter twitter.NitterResponse

				nitter, _ = twitter.FetchNitter(senderInfo.TwitterUsername, w.NitterClient)

				followers, _ := strconv.Atoi(nitter.Followers)

				status := fren_utils.AssertImportance(followers, fren_utils.Followers)

				if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
					BaseAddress:     sender,
					Status:          string(status),
					Followers:       fmt.Sprint(followers),
					TwitterUsername: senderInfo.TwitterUsername,
					TwitterName:     senderInfo.TwitterName,
					TwitterURL:      senderInfo.TwitterPfpUrl,
					UserID:          senderInfo.Id,
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
				Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s sold %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Seller](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – unfiltered sells",
							IconURL: bot.Image,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Buyer Balance",
								Value:  fmt.Sprintf("%v", utils.WeiToEther(balance)) + "Ξ",
								Inline: true,
							},
							{
								Name:   "Sh. Amt. | Price",
								Value:  fmt.Sprintf("%v | %v", sharesAmount, utils.WeiToEther(sharePrice)) + "Ξ",
								Inline: true,
							},
							{
								Name:   "IsBotPurchase",
								Value:  fmt.Sprint(isBot),
								Inline: true,
							},
							{
								Name:  "Transaction Hash",
								Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
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
				Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s sold %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Seller](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – filtered sells",
							IconURL: bot.Image,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Buyer Balance",
								Value:  fmt.Sprintf("%v Ξ", utils.WeiToEther(balance)),
								Inline: true,
							},
							{
								Name:   "Sh. Amt. | Price",
								Value:  fmt.Sprintf("%v | %v Ξ", sharesAmount, utils.WeiToEther(sharePrice)),
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
							{
								Name:  "Transaction Hash",
								Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
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

			nitter, err = twitter.FetchNitter(senderInfo.TwitterUsername, w.NitterClient)
			if err != nil {
				return err
			}

			followers, _ := strconv.Atoi(nitter.Followers)

			status := fren_utils.AssertImportance(followers, fren_utils.Followers)

			if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
				BaseAddress:     sender,
				Status:          string(status),
				Followers:       fmt.Sprint(followers),
				TwitterUsername: senderInfo.TwitterUsername,
				TwitterName:     senderInfo.TwitterName,
				TwitterURL:      senderInfo.TwitterPfpUrl,
				UserID:          senderInfo.Id,
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

		w.Bot.BotWebhook(&discordgo.MessageSend{
			Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       bot.Purple,
					Title:       fmt.Sprintf("%s sold %v key(s) of %s", senderInfo.TwitterUsername, sharesAmount, recipientInfo.TwitterUsername),
					Description: fmt.Sprintf("[Seller](https://www.friend.tech/rooms/%s)", sender),
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "friendtech – unfiltered sells",
						IconURL: bot.Image,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Buyer Balance",
							Value:  fmt.Sprintf("%v Ξ", utils.WeiToEther(balance)),
							Inline: true,
						},
						{
							Name:   "Sh. Amt. | Price",
							Value:  fmt.Sprintf("%v | %v Ξ", sharesAmount, utils.WeiToEther(sharePrice)),
							Inline: true,
						},
						{
							Name:   "IsBotPurchase",
							Value:  fmt.Sprint(isBot),
							Inline: true,
						},
						{
							Name:  "Transaction Hash",
							Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
						},
					},
				},
			},
		}, bot.FriendTechAllLogs)
	}
	return nil
}

// handleBuy
func (w *Watcher) handleBuy(tx *types.Transaction, sender, txData string) error {
	var channel string
	var err error
	var isBot bool

	data, err := utils.DecodeTransactionInputData(w.ABI, txData)
	if err != nil {
		return err
	}

	recipient := data["sharesSubject"].(common.Address)
	sharesAmount := data["amount"].(*big.Int)
	if sharesAmount.Int64() > 1 {
		isBot = true
	}

	senderInfo, err := fren_utils.GetUserInformation(sender, w.WatcherClient)
	if err != nil {
		return fmt.Errorf("bot purchase")
	}

	sharePrice := new(big.Int)
	sharePrice.SetString(senderInfo.DisplayPrice, 10)

	recipientInfo, _ := fren_utils.GetUserInformation(recipient.String(), w.WatcherClient)

	balance, err := utils.GetEthWalletBalance(w.WSSClient, common.HexToAddress(sender))
	if err != nil {
		return err
	}

	if isSelf(hex.EncodeToString(tx.Data()), sender) {
		// goes to self-buy channel.
		if _, ok := w.Addresses[strings.ToLower(sender)]; !ok {
			// it's not on our list & doesn't go to filtered.
			var user *models.FriendTechMonitorAll

			user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
			if err != nil {
				// we can't find it stored in our database, we fetch the data & add it to the database.
				var nitter twitter.NitterResponse

				nitter, _ = twitter.FetchNitter(senderInfo.TwitterUsername, w.NitterClient)

				followers, _ := strconv.Atoi(nitter.Followers)
				status := fren_utils.AssertImportance(followers, fren_utils.Followers)

				if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
					BaseAddress:     sender,
					Status:          string(status),
					Followers:       nitter.Followers,
					TwitterUsername: senderInfo.TwitterUsername,
					TwitterName:     senderInfo.TwitterName,
					TwitterURL:      senderInfo.TwitterPfpUrl,
					UserID:          senderInfo.Id,
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
				channel = bot.FriendTechWhalesBuys
			case "medium":
				channel = bot.FriendTechFishBuys
			case "low":
				channel = bot.FriendTechShrimpBuys
			case "none":
				channel = bot.FriendTechAllLogs
			}

			w.Bot.BotWebhook(&discordgo.MessageSend{
				Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s bought %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – unfiltered buys",
							IconURL: bot.Image,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Buyer Balance",
								Value:  fmt.Sprintf("%v Ξ", utils.WeiToEther(balance)),
								Inline: true,
							},
							{
								Name:   "Sh. Amt. | Price",
								Value:  fmt.Sprintf("%v | %v Ξ", sharesAmount, utils.WeiToEther(sharePrice)),
								Inline: true,
							},
							{
								Name:   "IsBotPurchase",
								Value:  fmt.Sprint(isBot),
								Inline: true,
							},
							{
								Name:  "Transaction Hash",
								Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
							},
						},
					},
				},
			}, channel)
		} else {
			// it's monitored from the w.Addresses list.
			var user *models.FriendTechMonitor

			user, _ = w.DB.Monitor.GetUserByAddress(sender, context.Background())

			w.Bot.BotWebhook(&discordgo.MessageSend{
				Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       fmt.Sprintf("%s bought %v key(s) of himself", senderInfo.TwitterUsername, sharesAmount),
						Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "friendtech – filtered buys",
							IconURL: bot.FriendTechImage,
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Buyer Balance",
								Value:  fmt.Sprintf("%v Ξ", utils.WeiToEther(balance)),
								Inline: true,
							},
							{
								Name:   "Sh. Amt. | Price",
								Value:  fmt.Sprintf("%v | %v Ξ", sharesAmount, utils.WeiToEther(sharePrice)),
								Inline: true,
							},
							{
								Name:   "IsBotPurchase",
								Value:  fmt.Sprint(isBot),
								Inline: true,
							},
							{
								Name:  "Transaction Hash",
								Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
							},
							{
								Name:   "User Importance",
								Value:  user.Status,
								Inline: true,
							},
							{
								Name:  "Buyer Twitter 🐰",
								Value: "",
							},
							{
								Name: "Seller Twitter",
							},
						},
					},
				},
			}, bot.FriendTechFilteredBuys)
		}
	} else { // is not a self buy
		var user *models.FriendTechMonitorAll
		user, err = w.DB.MonitorAll.GetUserByAddress(sender, context.Background())
		if err != nil { // means we can't find it & not stored in our db, we fetch the data
			var nitter twitter.NitterResponse

			nitter, err = twitter.FetchNitter(senderInfo.TwitterUsername, w.NitterClient)
			if err != nil {
				return err
			}

			followers, _ := strconv.Atoi(nitter.Followers)

			status := fren_utils.AssertImportance(followers, fren_utils.Followers)

			if err = w.DB.MonitorAll.InsertUser(&models.FriendTechMonitorAll{
				BaseAddress:     sender,
				Status:          string(status),
				Followers:       nitter.Followers,
				TwitterUsername: senderInfo.TwitterUsername,
				TwitterName:     senderInfo.TwitterName,
				TwitterURL:      senderInfo.TwitterPfpUrl,
				UserID:          senderInfo.Id,
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
			channel = bot.FriendTechWhalesBuys
		case "medium":
			channel = bot.FriendTechFishBuys
		case "low":
			channel = bot.FriendTechShrimpBuys
		case "none":
			channel = bot.FriendTechAllLogs
		}

		w.Bot.BotWebhook(&discordgo.MessageSend{
			Components: bot.BundleQuickTaskComponents(sender, "friendTech"),
			Embeds: []*discordgo.MessageEmbed{
				{
					Color:       bot.Purple,
					Title:       fmt.Sprintf("%s bought %v key(s) of %s", senderInfo.TwitterUsername, sharesAmount, recipientInfo.TwitterUsername),
					Description: fmt.Sprintf("[Buyer](https://www.friend.tech/rooms/%s)", sender),
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "friendtech – unfiltered buys",
						IconURL: bot.Image,
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Buyer Balance",
							Value:  fmt.Sprintf("%v Ξ", utils.WeiToEther(balance)),
							Inline: true,
						},
						{
							Name:   "Sh. Amt. | Price",
							Value:  fmt.Sprintf("%v | %v Ξ", sharesAmount, utils.WeiToEther(sharePrice)),
							Inline: true,
						},
						{
							Name:   "IsBotPurchase",
							Value:  fmt.Sprint(isBot),
							Inline: true,
						},
						{
							Name:  "Transaction Hash",
							Value: fmt.Sprintf("[%s](https://basescan.org/tx/%s)", tx.Hash().String(), tx.Hash().String()),
						},
					},
				},
			},
		}, channel)
	}
	return nil
}

// WatchRug watches self-sells, if 2+ self sells occurs within 5
// minutes we'll consider it as a rug.
func (w *Watcher) WatchRug() {
	w.Cache.InsertData("", "", cache.DefaultExpiration)
}

// serialize serializes data sent to the websocket.
// todo finish func
func (w *Watcher) serialize() {
	w.OutStreamData <- BroadcastData{}
}

func isSelf(txData string, sender string) bool {
	return strings.Contains(strings.ToLower(txData), strings.ToLower(sender[2:]))
}

func (tx *types.Transaction) isBotPurchaseNonRegistered() bool {
	return false
}
