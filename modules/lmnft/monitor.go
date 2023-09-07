package lmnft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (s *Settings) StartMonitor(networks []Network) {
	logger.LogStartup(moduleName)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogInfo(moduleName, fmt.Sprintf("program panicked! [%v]", r))
				s.StartMonitor(networks)
				return
			}
		}()
		for !s.monitorDrops(networks) {
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

// monitorDrops monitors SOL Hype Mints by default.
// If you want to monitor non SOL Mints, switch the "Solana" from the
// payload to the network you want to monitor.
func (s *Settings) monitorDrops(networks []Network) bool {
	var buf bytes.Buffer
	var filterByQuery string

	if len(networks) > 0 {
		var filterByNetworks []string
		for _, network := range networks {
			filterByNetworks = append(filterByNetworks, string(network))
		}
		filterByQuery = fmt.Sprintf("&& type:=[`%s`]", strings.Join(filterByNetworks, "`,`"))
	} else {
		filterByQuery = fmt.Sprintf("type:=[`%s`]", networks[0])
	}

	payload := searchPayload{
		Searches: Searches{
			{
				QueryBy:             "collectionName,owner",
				PerPage:             10,
				SortBy:              "lastMintedAt:desc",
				HighlightFullFields: "collectionName,owner",
				Collection:          "collections",
				Q:                   "*",
				FacetBy:             "soldOut,twitterVerified,type,cost",
				FilterBy:            "soldOut:=[false] && twitterVerified:=[true] " + filterByQuery,
				MaxFacetValues:      10,
				Page:                3,
			},
			{
				QueryBy:             "collectionName,owner",
				PerPage:             10,
				SortBy:              "lastMintedAt:desc",
				HighlightFullFields: "collectionName,owner",
				Collection:          "collections",
				Q:                   "*",
				FacetBy:             "soldOut",
				FilterBy:            "twitterVerified:=[true] " + filterByQuery,
				MaxFacetValues:      10,
				Page:                1,
			},
			{
				QueryBy:             "collectionName,owner",
				PerPage:             10,
				SortBy:              "lastMintedAt:desc",
				HighlightFullFields: "collectionName,owner",
				Collection:          "collections",
				Q:                   "*",
				FacetBy:             "twitterVerified",
				FilterBy:            "soldOut:=[false] " + filterByQuery,
				MaxFacetValues:      10,
				Page:                1,
			},
			{
				QueryBy:             "collectionName,owner",
				PerPage:             10,
				SortBy:              "lastMintedAt:desc",
				HighlightFullFields: "collectionName,owner",
				Collection:          "collections",
				Q:                   "*",
				FacetBy:             "type",
				FilterBy:            "soldOut:=[false] && twitterVerified:=[true]",
				MaxFacetValues:      10,
				Page:                1,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		logger.LogError(moduleName, err)
		return false
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: "s.launchmynft.io", Path: "/multi_search?x-typesense-api-key=UkN4Vnd3V2JMWWVIRlFNcTJ3dng4VGVtMGtvVGxBcmJJTTFFYS9MNXp1WT1Ha3dueyJmaWx0ZXJfYnkiOiJoaWRkZW46ZmFsc2UiLCJleGNsdWRlX2ZpZWxkcyI6ImhpZGRlbiIsInF1ZXJ5X2J5IjoiY29sbGVjdGlvbk5hbWUsb3duZXIiLCJsaW1pdF9oaXRzIjoyMDAsInNuaXBwZXRfdGhyZXNob2xkIjo1MH0%3D"},
		Body:   io.NopCloser(&buf),
	}

	resp, err := http.DefaultClient.Do(req)
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

	var res resLaunchMyNFT
	var t Webhook
	if err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&res); err != nil {
		return false
	}

	if err = resp.Body.Close(); err != nil {
		return false
	}

	for i := 0; i < len(res.Results); i++ {
		var d discord.Webhook

		embed := d.Embeds[0]
		embedsField := embed.Fields

		embed.Title = res.Results[i].Hits[i].Document.CollectionName
		embed.Description = res.Results[i].Hits[i].Document.Description
		embed.Url = "https://www.launchmynft.io/collections/" + res.Results[i].Hits[i].Document.Owner + "/" + res.Results[i].Hits[i].Document.ID
		embed.Thumbnail.Url = res.Results[i].Hits[i].Document.CollectionCoverURL

		var contract, discordServer, twitterAccount, secondary string
		switch res.Results[i].Hits[i].Document.Type {
		case string(Solana):
			contract, discordServer, twitterAccount = scrapeInformation(t.MintLink)
			secondary = fmt.Sprintf("[Secondary Market](https://hyperspace.xyz/collection/%s", contract)
			if twitterAccount == "" {
				twitterAccount = "no account :("
			} else if discordServer == "" {
				discordServer = "no server :("
			} else {
				discordServer = "[Server](https://discord.gg/" + discordServer + ")"
				twitterAccount = "[Account](https://twitter.com/" + twitterAccount + ")"
			}
			//todo add support for other platforms
		case string(Polygon):
		case string(Ethereum):
		case string(Binance):
		case string(Aptos):
		case string(Avalanche):
		case string(Fantom):
		case string(Stacks):
		default:
			logger.LogError(moduleName, errors.New("unknown network"))
			continue
		}

		{
			embedsField[0].Name = "Price"
			embedsField[0].Value = fmt.Sprintf("%.2f %s", t.Cost, t.Network)
		}
		{
			embedsField[1].Name = "Supply"
			embedsField[1].Value = fmt.Sprintf("`%d/%d – %.2f%%`", res.Results[i].Hits[i].Document.TotalMints, res.Results[i].Hits[i].Document.MaxSupply, res.Results[i].Hits[i].Document.FractionMinted*100)
		}
		{
			embedsField[2].Name = fmt.Sprintf("%s Contract", res.Results[i].Hits[i].Document.Type)
			embedsField[2].Value = fmt.Sprintf("`%s`", contract)
		}
		{
			embedsField[3].Name = "Twitter"
			embedsField[3].Value = twitterAccount
		}
		{
			embedsField[4].Name = "Discord"
			embedsField[4].Name = discordServer
		}
		{
			embedsField[5].Name = "Secondary"
			embedsField[5].Value = secondary
		}

		t.Fraction = res.Results[i].Hits[i].Document.FractionMinted * 100

		s.Handler.M.Set(t.Name, t.TotalMinted)

		if _, ok := s.Handler.MCopy.Get(t.Name); ok {

			continue
		}

		if t.Fraction >= 6 && t.TotalMinted >= 100 && t.CMID != "" {

			for _, b := range embedsField {
				b.Inline = true
			}

			if err = s.Discord.SendNotification(discord.Webhook{
				Username:  s.Discord.ProfileName,
				AvatarUrl: s.Discord.AvatarImage,
				Embeds: []discord.Embed{
					{
						Title:       t.Name,
						Description: t.Description,
						Url:         t.MintLink,
						Timestamp:   discord.GetTimestamp(),
						Color:       s.Discord.Color,
						Footer: discord.EmbedFooter{
							Text:    s.Discord.FooterText,
							IconUrl: s.Discord.FooterImage,
						},
						Thumbnail: discord.EmbedThumbnail{
							Url: t.Image,
						},
						Fields: []discord.EmbedFields{
							{
								Name:   "Price",
								Value:  fmt.Sprintf("%.2f %s", t.Cost, t.Network),
								Inline: true,
							},
							{
								Name:   "Supply",
								Value:  "`" + strconv.Itoa(t.TotalMinted) + "/" + strconv.Itoa(t.Supply) + fmt.Sprintf(" — %.2f%%`", t.Fraction),
								Inline: true,
							},
							{
								Name:   "CandyMachine ID",
								Value:  "`" + t.CMID + "`",
								Inline: false,
							},
							{
								Name:   "Twitter",
								Value:  t.Twitter,
								Inline: true,
							},
							{
								Name:   "Discord",
								Value:  t.Discord,
								Inline: true,
							},
							{
								Name:   "Hyperspace",
								Value:  "[Secondary Market]" + "(" + t.Secondary + ")",
								Inline: true,
							},
						},
					},
				},
			}, s.Discord.Webhook); err != nil {
				logger.LogError(moduleName, err)
			}
		} else {
			logger.LogInfo(moduleName, fmt.Sprintf("🐙 collection too low: %s", res.Results[i].Hits[i].Document.CollectionName))
		}

		s.Handler.M.ForEach(func(k string, v any) {
			s.Handler.MCopy.Set(k, v)
		})
	}

	return false
}

func scrapeInformation(input string) (string, string, string) {

	defaultVal := "–"

	candymachineUrl, err := url.Parse(input)
	if err != nil {
		log.Errorf("lmnft.scrapeCMID: ERROR Parsing cmidURL [%w]", err)
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    candymachineUrl,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("lmnft.ScrapeCMID: ERROR Client [%w]", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

	}

	var res resSolana
	switch resp.StatusCode {
	case 200:
	case 400:
	default:
		return defaultVal, defaultVal, defaultVal
	}

	//var res resCMID
	doc, errDoc := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if errDoc != nil {
		log.Errorf("lmnft.ScrapeCMID: ERROR Decoding JSON Script [%w]", err)
		return defaultVal, defaultVal, defaultVal
	}

	candy := doc.Find("script[id=__NEXT_DATA__]").Text()

	err = json.NewDecoder(strings.NewReader(candy)).Decode(&res)
	if err != nil {
		log.Errorf("lmnft.ScrapeCMID: ERROR Decoding JSON Script [%w]", err)
		return defaultVal, defaultVal, defaultVal
	}

	return res.Props.PageProps.Collection.NewCandyMachineAccountID, res.Props.PageProps.Collection.Discord, res.Props.PageProps.Collection.Twitter
}

type searchPayload struct {
	Searches Searches `json:"searches"`
}

type Searches []struct {
	QueryBy             string `json:"query_by"`
	PerPage             int    `json:"per_page"`
	SortBy              string `json:"sort_by"`
	HighlightFullFields string `json:"highlight_full_fields"`
	Collection          string `json:"collection"`
	Q                   string `json:"q"`
	FacetBy             string `json:"facet_by"`
	FilterBy            string `json:"filter_by"`
	MaxFacetValues      int    `json:"max_facet_values"`
	Page                int    `json:"page"`
}
