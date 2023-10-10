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
	"github.com/holiman/uint256"
	"github.com/weeaa/nft/discord/bot"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"log"
	"math/big"
	"net/url"
	"os"
	"strings"
	"time"
)

func New(bot *bot.Bot, addresses []string, proxyFilePath, HttpNode, WssNode string) (*Watcher, error) {

	HttpClient, err := ethclient.Dial(HttpNode)
	if err != nil {
		return nil, err
	}

	WssClient, err := ethclient.Dial(WssNode)
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	abiJSON, err := abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"))

	return &Watcher{
		ABI:        abiJSON,
		Bot:        bot,
		Addresses:  addresses,
		Client:     tls.New(tls.RandProxyFromList(proxyList)),
		HTTPClient: HttpClient,
		WSSClient:  WssClient,
		ProxyList:  proxyList,
		Handler:    handler.New(),
	}, nil
}

func (w *Watcher) StartWatchers() {
	w.ListenOnChain()
	for !w.WatchNewUsers() {
		select {
		default:
			continue
		}
	}
}

func (w *Watcher) ListenOnChain() {
	ch := make(chan types.Log)
	sub, err := utils.CreateSubscription(w.WSSClient, []common.Address{common.HexToAddress("0xA481B139a1A654cA19d2074F174f17D7534e8CeC")}, ch)
	if err != nil {
		logger.LogError(watcher, err)
		return
	}

	for {
		select {
		case subErr := <-sub.Err():
			logger.LogError(watcher, subErr)
			return
		case txn := <-ch:
			go w.dispatchLog(txn)
		}
	}
}

func (w *Watcher) dispatchLog(txn types.Log) {
	var tx *types.Transaction
	var sender common.Address
	var balance *big.Int
	var isPending bool
	var err error

	_ = balance
	tx, isPending, err = w.HTTPClient.TransactionByHash(context.Background(), txn.TxHash)
	if err != nil {
		logger.LogError("moduleName", err)
	}

	if isPending {

	}

	sender, err = utils.GetSender(tx)
	if err != nil {
		logger.LogError("moduleName", err)
	}

	//log.Println(hex.EncodeToString(tx.Data()))
	//	log.Println(string(tx.Data()))
	log.Println(hex.EncodeToString(tx.Data()))
	log.Println(utils.DecodeTransactionInputData(w.ABI, hex.EncodeToString(tx.Data())))
	//	log.Println("sender", sender.String())

	os.Exit(0)
	balance, err = utils.GetEthWalletBalance(w.HTTPClient, sender)
	if err != nil {
		logger.LogError("moduleName", err)
	}

	if strings.Contains(tx.Hash().Hex(), "sellMethod") {
		if err = w.handleSell(tx, sender.String()); err != nil {
			logger.LogError("moduleName", err)
		}
	} else if strings.Contains(tx.Hash().Hex(), buyMethod) {
		log.Println("is buy")
		w.handleBuy(tx, sender.String())
	} else {

	}
}

func (w *Watcher) handleSell(tx *types.Transaction, sender string) error {
	if isSelf(hex.EncodeToString(tx.Data()), sender) {
		//user, err := w.DB.GetUser(sender)
		//if err != nil {
		// use the imp status to display in the right channel what the person bought

		//return s.Discord.SendNotification(buildWebhook(), "")
	} else {
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

		go w.Bot.BotWebhook(&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					//bot.BundleQuickTaskComponents(sender)
					Title: fmt.Sprintf("%s bought %s", sender, recipient),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Share Amount",
							Value:  fmt.Sprint(sharesAmount),
							Inline: true,
						},
						{
							Name:   "IsBot",
							Value:  fmt.Sprint(isBot),
							Inline: true,
						},
					},
				},
			},
		}, "")
		return nil
	}

	return nil
}

func (w *Watcher) handleBuy(tx *types.Transaction, sender string) {
	if isSelf(hex.EncodeToString(tx.Data()), sender) {

	} else {

	}
}

func isSelf(txData string, sender string) bool {
	return strings.Contains(strings.ToLower(txData), strings.ToLower(sender[2:]))
}

func (w *Watcher) WatchNewUsers() bool {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "api.starsarena.com", Path: "/user/page"},
		Header: http.Header{
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"dnt":                {"1"},
			"sec-ch-ua-mobile":   {"?0"},
			"authorization":      {"Bearer its not working anymore :>"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
			"accept":             {"application/json"},
			"referer":            {""},
			"sec-ch-ua-platform": {"\"macOS\""},
		},
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

	var newUsers StarsArenaNewUsersResponse
	if err = json.Unmarshal(body, &newUsers); err != nil {
		logger.LogError(watcher, err)
		return false
	}

	for _, user := range newUsers.Users {
		_, ok := w.Handler.M.Get(user.Address)
		if !ok {
			w.Handler.M.Set(user.Address, user.Id)

			i := new(big.Int)
			i.SetString(user.KeyPrice, 10)
			keyPrice := utils.WeiToEther(i)

			var channel, roleID string
			status := AssertImportance(user.TwitterFollowers, Followers)
			switch status {
			case "none":
				channel = bot.StarsArenaNewUsers
			case Shrimp:
				//roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.ShrimpEmote])
				channel = bot.StarsArenaNewUsers5k
			case Whale:
				roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.WhaleEmote])
				channel = bot.StarsArenaNewUsers50k
			case Fish:
				//roleID = fmt.Sprintf("<@&%s>", bot.EmojiRoleMap[fren_utils.FishEmote])
				channel = bot.StarsArenaNewUsers10k
			}

			avaxBalance, _ := utils.GetEthWalletBalance(w.HTTPClient, common.HexToAddress(user.Address))
			var buys, sells int
			var supply string

			if user.Stats != nil {
				supply = user.Stats.Supply
				buys = user.Stats.Buys
				sells = user.Stats.Sells
			} else {
				supply = "n/a"
				buys = 0
				sells = 0
			}

			w.Bot.BotWebhook(&discordgo.MessageSend{
				Components: bot.BundleQuickTaskComponents(user.Address, "starsArena"),
				Content:    roleID,
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       user.TwitterHandle,
						URL:         fmt.Sprintf("https://www.starsarena.com/%s/", user.TwitterHandle),
						Description: fmt.Sprintf("[%s](https://snowtrace.io/address/%s)", user.Address, user.Address),
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: user.TwitterPicture,
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "@starsarena – new users [" + string(status) + "]",
							IconURL: "https://www.starsarena.com/assets/logo/starshares-sm.png",
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Price",
								Value:  fmt.Sprintf("%f AVAX", keyPrice),
								Inline: true,
							},
							{
								Name:   "Followers | Supply",
								Value:  fmt.Sprintf("%d | %s", user.TwitterFollowers, supply),
								Inline: true,
							},
							{
								Name:   "Balance",
								Value:  fmt.Sprintf("%f AVAX", utils.WeiToEther(avaxBalance)),
								Inline: true,
							},
							{
								Name:   "Buys | Sells",
								Value:  fmt.Sprintf("```%d | %d```", buys, sells),
								Inline: true,
							},
							{
								Name:   "Twitter Name",
								Value:  fmt.Sprintf("[%s](https://x.com/%s)", user.TwitterName, user.TwitterHandle),
								Inline: true,
							},
							{
								Name:   "Twitter Username",
								Value:  fmt.Sprintf("```%s```", user.TwitterHandle),
								Inline: true,
							},
						},
					},
				},
			}, channel)
		}
	}
	time.Sleep(1 * time.Second)
	return false
}

// WatchNewHandle monitors Twitter usernames which signs up on Stars Arena.
func (w *Watcher) WatchNewHandle(handleTwitterName string) {
	go func(twitterUsername string) {
		f := func() bool {
			req := &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Scheme: "https", Host: "api.starsarena.com", Path: "/user/handle?handle=" + twitterUsername},
				Header: http.Header{},
			}

			resp, err := w.Client.Do(req)
			if err != nil {
				logger.LogError("", err)
				return false
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
					tls.HandleRateLimit(w.Client, w.ProxyList, watcher)
					return false
				}
				return false
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.LogError("", err)
				return false
			}

			var response map[string]any
			if err = json.Unmarshal(body, &response); err != nil {
				logger.LogError("", err)
				return false
			}

			user, ok := response["user"].(map[string]any)
			if !ok {
				return false
			}

			// double check lol
			if user == nil {
				return false
			}

			address := user["address"].(string)
			twitterFollowers := user["twitterFollowers"].(int)
			twitterName := user["twitterName"].(string)
			twitterPicture := user["twitterPicture"].(string)

			keyPrice := new(big.Int)
			keyPrice.SetString(user["keyPrice"].(string), 10)

			timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", user["createdOn"].(string))
			if err != nil {
				return false
			}

			status := AssertImportance(twitterFollowers, Followers)

			avaxBalance, _ := utils.GetEthWalletBalance(w.HTTPClient, common.HexToAddress(address))

			w.Bot.BotWebhook(&discordgo.MessageSend{
				Content:    fmt.Sprintf("<@&%s>", bot.EmojiRoleMap["🫱🏻‍🫲🏾"]),
				Components: bot.BundleQuickTaskComponents(address, bot.StarsArena),
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       bot.Purple,
						Title:       twitterUsername,
						Description: fmt.Sprintf("[%s](https://snowtrace.io/address/%s)\n\nJoined <t:%d>", address, address, timestamp.Unix()),
						URL:         "https://www.starsarena.com/" + twitterUsername,
						Thumbnail: &discordgo.MessageEmbedThumbnail{
							URL: twitterPicture,
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "@starsarena – new users [" + string(status) + "]",
							IconURL: "https://www.starsarena.com/assets/logo/starshares-sm.png",
						},
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Price",
								Value:  fmt.Sprintf("%f AVAX", utils.WeiToEther(keyPrice)),
								Inline: true,
							},
							{
								Name:   "Followers | Supply",
								Value:  fmt.Sprintf("%d | %s", twitterFollowers, "N/A"),
								Inline: true,
							},
							{
								Name:   "Balance",
								Value:  fmt.Sprintf("%f AVAX", utils.WeiToEther(avaxBalance)),
								Inline: true,
							},
							{
								Name:   "Twitter Name",
								Value:  fmt.Sprintf("[%s](https://x.com/%s)", twitterName, twitterUsername),
								Inline: true,
							},
							{
								Name:   "Twitter Username",
								Value:  fmt.Sprintf("```%s```", twitterUsername),
								Inline: true,
							},
						},
					},
				},
			}, bot.StarsArenaFeedPing)

			return true
		}

		for !f() {
			select {
			default:
				time.Sleep(5 * time.Second) // high delay as a prevention to handling lot of users & risking banning proxies too fast.
				continue
			}
		}
	}(handleTwitterName)
}

func (w *Watcher) WatchContractInternalTransactions() {

}

// GetAvalancheInternalTxns returns latest internal transactions from a given contract.
func (w *Watcher) GetAvalancheInternalTxns(address string) (AvalancheInternalTransactionsResponse, error) {
	var response AvalancheInternalTransactionsResponse

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "https",
			Host:   "glacier-api.avax.network",
			Path:   fmt.Sprintf("/v1/chains/43114/addresses/%s/transactions:listInternals?pageSize=100", address)},
		Header: http.Header{},
	}

	resp, err := w.Client.Do(req)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
			time.Sleep(1 * time.Second)
			tls.HandleRateLimit(w.Client, w.ProxyList, watcher)
			return response, err
		}
		logger.LogError(watcher, fmt.Errorf("status %s for id: %d", resp.Status, w.Counter))
		return response, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}

	return response, nil
}
