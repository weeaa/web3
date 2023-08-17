package etherscan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const moduleName = "Etherscan Verified Contracts"

func Monitor(client discord.Client) {

	logger.LogStartup(moduleName)

	h := handler.New()
	e := &Webhook{}

	go func() {
		for {
			req := &http.Request{
				Method: http.MethodGet,
				URL:    &url.URL{Scheme: "https://", Host: "etherscan.io", Path: "/contractsVerified"},
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

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Error("etherscan.Monitor: Client ERR When Performing Request", "error", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error("etherscan.Monitor: Unable to Read Body", "error", err)
				continue
			}

			resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Error("etherscan.Monitor: Invalid Response Code", "respStatus", resp.Status)
				continue
			}

			document, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
			if err != nil {
				log.Error("etherscan.Monitor: Unable to Read Body Document", "error", err)
				continue
			}

			e.Address = trimFirstRune(document.Find("td").First().Text())
			e.Link = "https://etherscan.io/address" + "/" + e.Address
			e.Name = document.Find("td").First().Next().Text()

			h.M.Set(e.Address, e.Address)

			if h.M.Get(e.Address) == h.MCopy.Get(e.Address) {
				time.Sleep(retryDelay * time.Millisecond)
				continue
			}

			if err = client.EtherscanNotification(discord.Webhook{
				Content:   "",
				Username:  "ETH Verified Contract",
				AvatarUrl: "https://media.discordapp.net/attachments/1025845688900259980/1026979346046525520/ARCANA_LOGO_ICON_BG_BLACK_-_AQUA.png?width=1318&height=1318",
				Embeds: []discord.Embed{
					{
						Title:       e.Name,
						Description: "",
						Url:         e.Link,
						Timestamp:   discord.GetTimestamp(),
						Color:       client.Color,
						Footer: discord.EmbedFooter{
							Text: "ETH Verified Contract — Arcana",
						},

						Fields: []discord.EmbedFields{
							{
								Name:   "Contract Address",
								Value:  "`" + e.Address + "`",
								Inline: true,
							},
							{
								Name:   "Write Contract | Code",
								Value:  "[Contract](https://etherscan.io/address/" + e.Address + "#writeContract) | [Contract](https://etherscan.io/address/" + e.Address + "#code)",
								Inline: true,
							},
						},
					},
				},
			}); err != nil {
				logger.LogError(moduleName, fmt.Errorf("unable to Send Discord Webhook: %w", err))
			}

			h.M.ForEach(func(k string, v interface{}) {
				h.MCopy.Set(k, v)
			})

			time.Sleep(retryDelay * time.Millisecond)
		}
	}()
}

func trimFirstRune(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return ""
}
