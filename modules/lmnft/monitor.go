package lmnft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const moduleName = "LaunchMyNFT"
const retryDelay = 2500 * time.Millisecond

// Monitor monitors SOL Hype Mints by default.
// If you want to monitor non SOL Mints, switch the "Solana" from the
// payload to the network you want to monitor.
func Monitor(client *discord.Client, networks []Network, delay time.Duration) {

	logger.LogStartup(moduleName)

	h := handler.New()
	t := &Webhook{}

	var buf bytes.Buffer
	var filterByQuery string

	defer func() {
		if r := recover(); r != nil {
			Monitor(client, networks, delay)
			return
		}
	}()

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
	}

	go func() {
		for {

			//payload := strings.NewReader("{\"searches\":[{\"query_by\":\"collectionName,owner\",\"per_page\":10,\"sort_by\":\"lastMintedAt:desc\",\"highlight_full_fields\":\"collectionName,owner\",\"collection\":\"collections\",\"q\":\"*\",\"facet_by\":\"soldOut,twitterVerified,type,cost\",\"filter_by\":\"soldOut:=[false] && twitterVerified:=[true] && type:=[`Solana`]\",\"max_facet_values\":10,\"page\":3},{\"query_by\":\"collectionName,owner\",\"per_page\":10,\"sort_by\":\"lastMintedAt:desc\",\"highlight_full_fields\":\"collectionName,owner\",\"collection\":\"collections\",\"q\":\"*\",\"facet_by\":\"soldOut\",\"filter_by\":\"twitterVerified:=[true] && type:=[`Solana`]\",\"max_facet_values\":10,\"page\":1},{\"query_by\":\"collectionName,owner\",\"per_page\":10,\"sort_by\":\"lastMintedAt:desc\",\"highlight_full_fields\":\"collectionName,owner\",\"collection\":\"collections\",\"q\":\"*\",\"facet_by\":\"twitterVerified\",\"filter_by\":\"soldOut:=[false] && type:=[`Solana`]\",\"max_facet_values\":10,\"page\":1},{\"query_by\":\"collectionName,owner\",\"per_page\":10,\"sort_by\":\"lastMintedAt:desc\",\"highlight_full_fields\":\"collectionName,owner\",\"collection\":\"collections\",\"q\":\"*\",\"facet_by\":\"type\",\"filter_by\":\"soldOut:=[false] && twitterVerified:=[true]\",\"max_facet_values\":10,\"page\":1}]}")

			req := &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Scheme: "https", Host: "s.launchmynft.io", Path: "/multi_search?x-typesense-api-key=UkN4Vnd3V2JMWWVIRlFNcTJ3dng4VGVtMGtvVGxBcmJJTTFFYS9MNXp1WT1Ha3dueyJmaWx0ZXJfYnkiOiJoaWRkZW46ZmFsc2UiLCJleGNsdWRlX2ZpZWxkcyI6ImhpZGRlbiIsInF1ZXJ5X2J5IjoiY29sbGVjdGlvbk5hbWUsb3duZXIiLCJsaW1pdF9oaXRzIjoyMDAsInNuaXBwZXRfdGhyZXNob2xkIjo1MH0%3D"},
				Body:   io.NopCloser(&buf),
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}

			if resp.StatusCode != 200 {
				logger.LogError(moduleName, fmt.Errorf("invalid response status: %s", resp.Status))
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}

			var res resLaunchMyNFT
			if err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&res); err != nil {
				continue
			}

			if err = resp.Body.Close(); err != nil {
				continue
			}

			for i := 0; i < len(res.Results); i++ {
				t.Name = res.Results[i].Hits[i].Document.CollectionName
				t.Fraction = res.Results[i].Hits[i].Document.FractionMinted * 100
				t.Cost = res.Results[i].Hits[i].Document.Cost
				t.Supply = res.Results[i].Hits[i].Document.MaxSupply
				t.TotalMinted = res.Results[i].Hits[i].Document.TotalMints
				t.Image = res.Results[i].Hits[i].Document.CollectionCoverURL
				t.Description = res.Results[i].Hits[i].Document.Description
				t.MintLink = "https://www.launchmynft.io/collections/" + res.Results[i].Hits[i].Document.Owner + "/" + res.Results[i].Hits[i].Document.ID

				switch res.Results[i].Hits[i].Document.Type {
				case string(Solana):
					t.CMID, t.Discord, t.Twitter = scrapeCMID(t.MintLink)
					t.Secondary = "https://hyperspace.xyz/collection/" + t.CMID
					if t.Twitter == "" {
						t.Twitter = "no account :("
					} else if t.Discord == "" {
						t.Discord = "no server :("
					} else {
						t.Discord = "[Server](https://discord.gg/" + t.Discord + ")"
						t.Twitter = "[Account](https://twitter.com/" + t.Twitter + ")"
					}
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

				h.M.Set(t.Name, t.TotalMinted)

				if h.M.Get(t.Name) == h.MCopy.Get(t.Name) {
					time.Sleep(40 * time.Second)
					continue
				}

				if t.Fraction >= 6 && t.TotalMinted >= 100 && t.CMID != "" {
					if err = client.SendNotification(discord.Webhook{
						Username:  "LaunchMyNFT",
						AvatarUrl: client.AvatarImage,
						Embeds: []discord.Embed{
							{
								Title:       t.Name,
								Description: t.Description,
								Url:         t.MintLink,
								Timestamp:   discord.GetTimestamp(),
								Color:       client.Color,
								Footer: discord.EmbedFooter{
									Text:    client.FooterText,
									IconUrl: client.FooterImage,
								},
								Thumbnail: discord.EmbedThumbnail{
									Url: t.Image,
								},
								Fields: []discord.EmbedFields{

									{
										Name:   "Price",
										Value:  fmt.Sprintf("%.2f SOL", t.Cost),
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
					}, discord.LaunchMyNFT); err != nil {
						logger.LogError(moduleName, err)
					}
				} else {
					//log.Warn("Collection 2 Low :(", "name", t.Name, "fraction", t.Fraction, "totalMinted", t.TotalMinted, "network", t.Network)
				}

				h.M.ForEach(func(k string, v interface{}) {
					h.MCopy.Set(k, v)
				})
			}
			time.Sleep(delay)
		}
	}()
}

// used for SOL
func scrapeCMID(input string) (string, string, string) {

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
