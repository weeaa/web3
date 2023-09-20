package etherscan

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewClient(discordClient *discord.Client, verbose bool) *Settings {
	return &Settings{
		Discord: discordClient,
		Handler: handler.New(),
		Verbose: verbose,
		Context: context.Background(),
	}
}

// StartMonitor monitors all newest ETH Verified Contracts audited by Etherscan.
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
		for !s.monitorVerifiedContracts() {
			select {
			case <-s.Context.Done():
				logger.LogShutDown(moduleName)
				return
			default:
				time.Sleep(time.Duration(s.RetryDelay) * time.Millisecond)
				continue
			}
		}
	}()
}

func (s *Settings) monitorVerifiedContracts() bool {
	resp, err := doRequest()
	if err != nil {
		return false
	}

	if resp.StatusCode != 200 {
		logger.LogError(moduleName, fmt.Errorf("invalid response status: %s", resp.Status))
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	if err = resp.Body.Close(); err != nil {
		return false
	}

	contract := ParseHTML(goquery.NewDocumentFromReader(strings.NewReader(string(body))))

	s.Handler.M.Set(contract.Address, contract.Name)
	if _, ok := s.Handler.MCopy.Get(contract.Address); ok {
		return false
	}

	s.Handler.Copy()

	if s.Discord.Webhook != "" {
		if err = s.sendDiscordNotification(contract); err != nil {
			logger.LogError(moduleName, err)
		}
	}

	if s.Verbose {
		logger.LogInfo(moduleName, fmt.Sprintf("🎈 new contract found: %s | %s", contract.Address, contract.Name))
	}

	return false
}

func doRequest() (*http.Response, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "etherscan.io", Path: "/contractsVerified"},
		Header: http.Header{
			"Authority":                 {"etherscan.io"},
			"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			"Accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
			"Cache-Control":             {"max-age=0"},
			"Referer":                   {"https://etherscan.io/contractsVerified"},
			"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\""},
			"Sec-Ch-Ua-Mobile":          {"?0"},
			"Sec-Ch-Ua-Platform":        {"\"Windows\""},
			"Sec-Fetch-Dest":            {"document"},
			"Sec-Fetch-Mode":            {"navigate"},
			"Sec-Fetch-Site":            {"same-origin"},
			"Sec-Fetch-User":            {"?1"},
			"Upgrade-Insecure-Requests": {"1"},
			"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"},
		},
	}
	return http.DefaultClient.Do(req)
}

func (s *Settings) sendDiscordNotification(contract Contract) error {
	return s.Discord.SendNotification(discord.Webhook{
		Username:  s.Discord.ProfileName,
		AvatarUrl: s.Discord.AvatarImage,
		Embeds: []discord.Embed{
			{
				Title:     contract.Name,
				Url:       "https://etherscan.io/address/" + contract.Address,
				Timestamp: discord.GetTimestamp(),
				Color:     s.Discord.Color,
				Footer: discord.EmbedFooter{
					Text:    s.Discord.FooterText,
					IconUrl: s.Discord.FooterImage,
				},

				Fields: []discord.EmbedFields{
					{
						Name:   "Contract Address",
						Value:  "`" + contract.Address + "`",
						Inline: true,
					},
					{
						Name:   "Write Contract | Code",
						Value:  "[Contract](https://etherscan.io/address/" + contract.Address + "#writeContract) | [Contract](https://etherscan.io/address/" + contract.Address + "#code)",
						Inline: true,
					},
				},
			},
		},
	}, s.Discord.Webhook)
}

func ParseHTML(document *goquery.Document, err error) Contract {
	if err != nil {
		return Contract{}
	}
	return Contract{
		Address: document.Find("td").First().Find("span").Find("a").AttrOr("title", ""),
		Name:    document.Find("td").First().Next().Text(),
	}
}
