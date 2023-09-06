package unisat

import (
	"encoding/json"
	"fmt"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Monitor monitors latest BRC20 Mints that are deducted as 'hype' mints.
// @params:
func Monitor(client *discord.Client, blackListedCoins *[]string) {

	h := handler.New()

	go func() {
		for {
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

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}

			if resp.StatusCode == 429 {
				logger.LogError(moduleName, rateLimited)
				time.Sleep(30 * time.Second)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			var r resInfo
			if err = json.Unmarshal(body, &r); err != nil {
				continue
			}

			if err = resp.Body.Close(); err != nil {
				continue
			}

			for _, brc := range r.Data.Detail {
				for _, coin := range blackListedCoins {
					if brc.Ticker != coin {
						var b discord.BRC20MintsWebhook

						supply, _ := strconv.Atoi(brc.Max)
						minted, _ := strconv.Atoi(brc.TotalMinted)
						rawPercentage := calculatePercentage(minted, supply)

						b.PercentageMinted = fmt.Sprintf("%.2f", rawPercentage)

						if brc.MintTimes >= 500 && rawPercentage >= 40 {
							var fees resFees

							b.Name = brc.Ticker
							b.Supply = brc.Max
							b.HoldersCount = fmt.Sprint(brc.HoldersCount)
							b.MintTimes = fmt.Sprint(brc.MintTimes)
							b.TotalMinted = brc.TotalMinted
							b.Creator, b.Block = getTickerInfo(brc.Ticker)
							b.Links = generateLinks(brc.Ticker, b.Creator)
							b.Holders = fetchHolders(brc.Ticker, supply)
							b.MintLink = "https://unisat.io/brc20/" + b.Name
							b.Timestamp = fmt.Sprint(brc.DeployBlocktime)
							b.BlockDeploy = fmt.Sprint(brc.DeployHeight)
							if fees, err = getFees(); err != nil {
								logger.LogError(moduleName, err)
							}

							b.Fees = fees.FastestFee

							value, ok := h.M.Get(b.Name)
							if ok {
								var pctF float64
								pctStr, isString := value.(string)
								if isString {
									pctF, err = strconv.ParseFloat(pctStr, 64)
									if err != nil {
										logger.LogError(moduleName, err)
									}

									if (rawPercentage - pctF) >= 3 {
										if err = client.SendNotification(discord.Webhook{}, moduleName); err != nil {
											logger.LogError(moduleName, err)
										}
									} else {
										//logger.Warn("Percentage increase not sufficient")
									}
								}

							} else {
								err = t.BRCHotMintNotification(webhookURL)
								if err != nil {
									log.Error("BRC NOTI ERROR", err)
								}
								log.Warn("NOT FIRST TIME")
							}

							h.M[t.Name] = t.PercentageMinted
							if h.M.Get() == h.MCopy[t.Name] {
								log.Info("is SAME")
								h.MCopy[t.Name] = h.M[t.Name]
								continue
							}

							h.M.ForEach(func(k string, v any) {
								h.MCopy.Set(k, v)
							})
						}
					}
					if brc.Ticker == "good" || brc.Ticker == "D3VL" {
						log.Printf("BLACKLISTED COIN")
					} else {
						//var doProceed bool

					}
				}
			}
			time.Sleep(700 * time.Millisecond)
		}
	}()
}

func fetchHolders(ticker string, supply int) map[int]map[string]string {

	holdersURL, err := url.Parse(fmt.Sprintf("https://unisat.io/brc20-api-v2/brc20/%s/holders?start=0&limit=5", ticker))
	if err != nil {
		log.Error(err)
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    holdersURL,
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("CLIENT", err)
	}

	fmt.Println("holders info", resp.Status)

	var res resHolders
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Error("DECODING", err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.LogError(moduleName, err)
	}

	var holders = make(map[int]map[string]string)
	for i, holder := range res.Data.Detail {
		mInfo := make(map[string]string)
		mInfo["address"] = holder.Address
		mInfo["balance"] = holder.OverallBalance
		balance, _ := strconv.Atoi(holder.OverallBalance)
		mInfo["percentage"] = fmt.Sprintf("%.2f", calculatePercentage(balance, supply))
		holders[i] = mInfo
	}

	return holders
}

func getTickerInfo(ticker string) (string, string) {

	tickerInfoURL, err := url.Parse(fmt.Sprintf("https://unisat.io/brc20-api-v2/brc20/%s/info", ticker))
	if err != nil {
		log.Error(err)
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    tickerInfoURL,
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("CLIENT TICKER", err)
	}

	fmt.Println("ticker info", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var res resTickerInfo
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("UNMARSHAL TICKER INFO", err)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Error(err)
	}

	return res.Data.Creator, fmt.Sprint(res.Data.DeployHeight)
}

func calculatePercentage(n1, n2 int) float64 {
	return float64(n1) / float64(n2) * 100
}

func generateLinks(ticker, deployer string) string {
	return fmt.Sprintf("[Unisat](https://unisat.io/unisat/%s) – [Deployer](https://btcscan.org/address/%s) – [Twitter Search](https://twitter.com/search?q=$%s&f=live)", ticker, deployer, ticker)
}

func getFees() (resFees, error) {

	var res resFees

	feesURL, err := url.Parse("https://bitcoinfees.billfodl.com/api/fees/")
	if err != nil {
		return res, err
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    feesURL,
		Header: http.Header{},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(body, &res)

	return res, err
}
