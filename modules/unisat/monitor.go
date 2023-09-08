package unisat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewClient(discordClient *discord.Client, verbose bool, client tls_client.HttpClient, proxyList []string, rotateOnProxyBan bool) *Settings {
	return &Settings{
		Discord:          discordClient,
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

	if resp.StatusCode != 200 {
		if ok := s.handleRateLimit(resp.StatusCode); !ok {
			return ok
		}
		logger.LogInfo(moduleName, fmt.Sprintf("unexpected response status: monitorDrops: %s", resp.Status))
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

			disc.Username = s.Discord.ProfileName

			embed := disc.Embeds[0]
			embedsField := embed.Fields
			rawHoldersData, ok := s.fetchHolders(brc.Ticker, supply)
			if !ok || rawHoldersData == nil {
				continue
			}

			holders, balance := prettyPrintHolders(rawHoldersData)

			embed.Title = brc.Ticker
			embed.Description = fmt.Sprintf("token deployed at: <t:%d> – block: `%d`", brc.DeployBlocktime, brc.DeployHeight)
			embed.Url = "https://unisat.io/brc20/" + brc.Ticker
			embed.Color = s.Discord.Color
			embed.Timestamp = discord.GetTimestamp()
			embed.Footer = discord.EmbedFooter{
				Text:    fmt.Sprintf("⛽️ %s sats/byte – %s", fees.FastestFee, s.Discord.FooterText),
				IconUrl: s.Discord.FooterImage,
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

			value, ok := s.Handler.M.Get(brc.Ticker)
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
	return false
}

// fetchHolders fetches 5 top holders of a BRC20 token on Unisat.
func (s *Settings) fetchHolders(ticker string, supply int) (map[int]map[string]string, bool) {
	var r ResHolders

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "unisat.io", Path: fmt.Sprintf("/brc20-api-v2/brc20/%s/holders?start=0&limit=5", ticker)},
		Header: http.Header{
			"Authority":          {"unisat.io"},
			"Accept":             {"application/json, text/plain, */*"},
			"Accept-Language":    {"en-US,en;q=0.9"},
			"Dnt":                {"1"},
			"Referer":            {"https://unisat.io/brc20/" + ticker},
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
		return nil, false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if ok := s.handleRateLimit(resp.StatusCode); !ok {
			return nil, ok
		}
		logger.LogInfo(moduleName, fmt.Sprintf("unexpected response status: monitorDrops: %s", resp.Status))
		return nil, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false
	}

	if err = json.Unmarshal(body, &r); err != nil {
		return nil, false
	}

	var holders = make(map[int]map[string]string)
	for i, holder := range r.Data.Detail {
		mInfo := make(map[string]string)
		mInfo["address"] = holder.Address
		mInfo["balance"] = holder.OverallBalance
		balance, _ := strconv.Atoi(holder.OverallBalance)
		mInfo["percentage"] = fmt.Sprintf("%.2f", calculatePercentage(balance, supply))
		holders[i] = mInfo
	}

	return holders, true
}

func (s *Settings) GetTickerInfo(ticker string) (ResTickerInfo, bool) {
	var r ResTickerInfo

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "unisat.io", Path: fmt.Sprintf("/brc20-api-v2/brc20/%s/info", ticker)},
		Header: http.Header{
			"Authority":          {"unisat.io"},
			"Accept":             {"application/json, text/plain, */*"},
			"Accept-Language":    {"en-US,en;q=0.9"},
			"Dnt":                {"1"},
			"Referer":            {"https://unisat.io/brc20/" + ticker},
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
		return r, false
	}

	if resp.StatusCode != 200 {
		if ok := s.handleRateLimit(resp.StatusCode); !ok {
			return r, false
		}
		logger.LogInfo(moduleName, fmt.Sprintf("unexpected response status: monitorDrops: %s", resp.Status))
		return r, false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return r, false
	}

	if err = json.Unmarshal(body, &r); err != nil {
		return r, false
	}

	return r, true
}

func calculatePercentage(n1, n2 int) float64 {
	return float64(n1) / float64(n2) * 100
}

func generateLinks(ticker, deployer string) string {
	return fmt.Sprintf("[Unisat](https://unisat.io/unisat/%s) – [Deployer](https://btcscan.org/address/%s) – [Twitter Search](https://twitter.com/search?q=$%s&f=live)", ticker, deployer, ticker)
}

// GetFees returns current BTC fees.
func GetFees() (ResFees, error) {
	var res ResFees

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "bitcoinfees.billfodl.com", Path: "/api/fees/"},
		Header: http.Header{
			"user-agent": {"golang :)"},
		},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return res, fmt.Errorf("unexpected response fetching BTC fees: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, &res)
	return res, err
}

// prettyPrintHolders pretty prints data passed as param.
func prettyPrintHolders(holders map[int]map[string]string) (string, string) {
	var addressesBuilder strings.Builder
	var balancesBuilder strings.Builder

	var keys []int
	for k := range holders {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, key := range keys {
		holder := holders[key]
		address := fmt.Sprintf("[%s](https://mempool.space/address/%s)", extractFirstAndLastFourLetters(holder["address"]), holder["address"])
		balance := fmt.Sprintf("%s (%s%%/spl)", holder["balance"], holder["percentage"])

		addressesBuilder.WriteString(address + "\n")
		balancesBuilder.WriteString(balance + "\n")
	}

	names := addressesBuilder.String()
	balances := balancesBuilder.String()

	return names, balances
}

func extractFirstAndLastFourLetters(input string) string {
	return fmt.Sprintf("%s...%s", input[:4], input[len(input)-4:])
}

func (s *Settings) handleRateLimit(respStatus int) bool {
	if respStatus == http.StatusTooManyRequests {
		if s.RotateProxyOnBan {
			if err := tls.RotateProxy(s.Client, s.ProxyList); err != nil {
				logger.LogError(moduleName, fmt.Errorf("unable to rotate proxy on client: %v", err))
				return false
			}
			logger.LogInfo(moduleName, "rotated proxy due to rate limit")
			return true
		}
		return false
	}
	return true
}
