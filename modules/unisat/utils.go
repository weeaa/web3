package unisat

import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// FetchHolders fetches 5 top holders of a BRC20 token on Unisat.
func (s *Settings) FetchHolders(ticker string, supply int) (map[int]map[string]string, error) {
	var r ResHolders

	ticker = "6f726469"
	//https://api.unisat.io/query-v4/brc20/66736174/holders?start=0&limit=20

	queries := url.Values{
		"start": {"0"},
		"limit": {"5"},
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "https", Host: host,
			Path:     fmt.Sprintf("%s%s/holders", path, ticker),
			RawQuery: queries.Encode(),
		},
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
			"x-appid":            {"1adcd79696032b1753f1812c9461cd36"},
			"x-sign":             {"2bb1317a6a02e93747ab5dba00bd7d95"},
			"x-ts":               {fmt.Sprint(time.Now().Unix())},
		},
	}

	log.Print(req.URL)
	log.Println("before req")
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println("after req")

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if ok := tls.HandleRateLimit(s.Client, s.ProxyList, ""); !ok {
			return nil, fmt.Errorf("")
		}
		logger.LogInfo(moduleName, fmt.Sprintf("unexpected response status: monitorDrops: %s", resp.Status))
		return nil, fmt.Errorf("")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println(string(body))

	if err = json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	holders := make(map[int]map[string]string)
	for i, holder := range r.Data.Detail {
		mInfo := make(map[string]string)
		mInfo["address"] = holder.Address
		mInfo["balance"] = holder.OverallBalance
		balance, _ := strconv.Atoi(holder.OverallBalance)
		mInfo["percentage"] = fmt.Sprintf("%.2f", calculatePercentage(balance, supply))
		holders[i] = mInfo
	}

	return holders, nil
}

// GetTickerInfo returns
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
		//if ok := s.handleRateLimit(resp.StatusCode); !ok {
		return r, false
		//}
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
		address := fmt.Sprintf("[%s](https://mempool.space/address/%s)", utils.FirstLastFour(holder["address"]), holder["address"])
		balance := fmt.Sprintf("%s (%s%%/spl)", holder["balance"], holder["percentage"])

		addressesBuilder.WriteString(address + "\n")
		balancesBuilder.WriteString(balance + "\n")
	}

	return addressesBuilder.String(), balancesBuilder.String()
}

// GetFees returns current BTC fees.
func GetFees() (ResFees, error) {
	var res ResFees

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "http", Host: "bitcoinfees.billfodl.com", Path: "/api/fees"},
		Header: http.Header{
			"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"},
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

func calculatePercentage(n1, n2 int) float64 {
	return float64(n1) / float64(n2) * 100
}

func generateLinks(ticker, deployer string) string {
	return fmt.Sprintf("[Unisat](https://unisat.io/unisat/%s) – [Deployer](https://btcscan.org/address/%s) – [Twitter Search](https://twitter.com/search?q=$%s&f=live)", ticker, deployer, ticker)
}
