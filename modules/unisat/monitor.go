package unisat

import (
	"context"
	"fmt"
	"github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/bwmarrin/discordgo"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"net/url"
	"time"
)

func NewClient(bot *bot.Bot, verbose bool, client tls_client.HttpClient, proxyList []string, rotateOnProxyBan bool) *Settings {
	return &Settings{
		Bot:              bot,
		Verbose:          verbose,
		Client:           client,
		ProxyList:        proxyList,
		RotateProxyOnBan: rotateOnProxyBan,
		Context:          context.Background(),
		Handler:          handler.New(),
	}
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
		for !s.monitorDrops() {
			select {
			case <-s.Context.Done():
				logger.LogShutDown(moduleName)
				return
			default:
				time.Sleep(2 * time.Minute)
				continue
			}
		}
	}()
}

// monitorDrops monitors latest BRC20 Mints that are deducted as 'hype' mints.

func (s *Settings) monitorDrops() bool {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "unisat.io", Path: "/brc20-api-v2/brc20/status?ticker=&start=0&limit=40&complete=no&sort=minted"},
		Header: http.Header{
			"Authority":          {"unisat.io"},
			"Accept":             {"application/json, text/plain, */*"},
			"Accept-Language":    {"en-US,en;q=0.9"},
			"Dnt":                {"1"},
			"Referer":            {"https://unisat.io/brc20"},
			"Sec-Ch-Ua":          {"\"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
			"Sec-Ch-Ua-Mobile":   {"?0"},
			"Sec-Ch-Ua-Platform": {"\"macOS\""},
			"Sec-Fetch-Dest":     {"empty"},
			"Sec-Fetch-Mode":     {"cors"},
			"Sec-Fetch-Site":     {"same-origin"},
			"User-Agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		},
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	/*
		if resp.StatusCode != 200 {
			logger.LogInfo(moduleName, fmt.Sprintf("unexpected response status: %s", resp.Status))
			if s.RotateProxyOnBan && resp.StatusCode == http.StatusTooManyRequests {
				if ok := tls.HandleRateLimit(s.Client, s.ProxyList, moduleName); !ok {
					return ok
				}
			}
			return false
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}

		var tickers ResTickers
		if err = json.Unmarshal(body, &tickers); err != nil {
			return false
		}

		for _, brc := range tickers.Data.Detail {

			encodedTicker := hex.EncodeToString([]byte(brc.Ticker))
			supply, _ := strconv.Atoi(brc.Max)
			minted, _ := strconv.Atoi(brc.TotalMinted)
			rawPercentage := calculatePercentage(minted, supply)

			// If it was minted more than 500 times & supply is 40% minted, we go further.
			if brc.MintTimes >= 500 && rawPercentage >= 40 {
				var disc discord.Webhook
				var fees ResFees

				if fees, err = GetFees(); err != nil {
					logger.LogError(moduleName, err)
					continue
				}

				//disc.Username = s.Discord.ProfileName

				embed := disc.Embeds[0]
				embedsField := embed.Fields
				rawHoldersData, ok := s.FetchHolders(encodedTicker, supply)

				holders, balance := prettyPrintHolders(rawHoldersData)

				embed.Title = brc.Ticker
				embed.Description = fmt.Sprintf("token deployed at: <t:%d> – block: `%d`", brc.DeployBlocktime, brc.DeployHeight)
				embed.Url = "https://unisat.io/brc20/" + brc.Ticker
				//embed.Color = s.Discord.Color
				embed.Timestamp = discord.GetTimestamp()
				embed.Footer = discord.EmbedFooter{
					Text: fmt.Sprintf("⛽️ %s sats/byte – %s", fees.FastestFee, s.Discord.FooterText),
					//IconUrl: s.Discord.FooterImage,
				}

				{
					embedsField[0].Name = "Supply | Minted"
					embedsField[0].Value = fmt.Sprintf("%s | %s", brc.Max, brc.TotalMinted)
				}
				{
					embedsField[1].Name = "Minted | Percentage"
					embedsField[1].Value = fmt.Sprintf("%d | %s", brc.MintTimes, fmt.Sprintf("%.2f", rawPercentage))
				}
				{
					embedsField[2].Name = "Holders No."
					embedsField[2].Value = fmt.Sprint(brc.HoldersCount)
				}
				{
					embedsField[3].Name = "Top Holders"
					embedsField[3].Value = holders
				}
				{
					embedsField[3].Name = "Balance"
					embedsField[3].Value = balance
				}
				{
					embedsField[4].Name = "Links"
					embedsField[4].Value = generateLinks(brc.Ticker, brc.Creator)
				}

				for _, b := range embedsField {
					b.Inline = true
				}

				//value, ok := s.Handler.M.Get(brc.Ticker)
				if ok {
					var pctF float64
					pctStr, isString := value.(string)
					if isString {
						pctF, err = strconv.ParseFloat(pctStr, 64)
						if err != nil {
							logger.LogError(moduleName, err)
							continue
						}

						if (rawPercentage - pctF) >= s.PercentageIncreaseBetweenRefresh { // was at 3 before
							if err = s.Discord.SendNotification(disc, s.Discord.Webhook); err != nil {
								logger.LogError(moduleName, err)
							}
							if s.Verbose {
								logger.LogInfo(moduleName, fmt.Sprintf("🦅 increase spotted for %s | %.2f > %.2f", brc.Ticker, pctF, rawPercentage))
							}
						} else {
							logger.LogInfo(moduleName, fmt.Sprintf("〽️percentage increase is not sufficient for %s", brc.Ticker))
						}
					}
				} else {
					s.Bot.BotWebhook(buildWebhook(), "")

					if err = s.Discord.SendNotification(disc, s.Discord.Webhook); err != nil {
						logger.LogError(moduleName, err)
					}
					if s.Verbose {
						logger.LogInfo(moduleName, fmt.Sprintf("😇 new ticker found: %s", brc.Ticker))
					}
				}

				s.Handler.M.Set(brc.Ticker, fmt.Sprintf("%.2f", rawPercentage))
			}
		}
	*/
	return false
}

func buildWebhook() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{},
		},
	}
}
