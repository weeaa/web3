package lmnft

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

// monitorDrops monitors SOL Hype Mints that are minted above 60%
// & have a minimum of 100 mints.
func (s *Settings) monitorDrops(networks []Network) bool {

	resp, err := doRequest(networks)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.LogError(moduleName, fmt.Errorf("invalid response status: %s", resp.Status))
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var response resLaunchMyNFT
	if err = json.Unmarshal(body, &response); err != nil {
		return false
	}

	for i := 0; i < len(response.Results); i++ {
		drop := &Release{}

		switch response.Results[i].Hits[i].Document.Type {
		case string(Solana):
			var info resSolana

			info, err = scrapeInformation[resSolana](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.NewCandyMachineAccountID
			drop.Discord = info.Props.PageProps.Collection.Discord
			drop.Twitter = info.Props.PageProps.Collection.Twitter

		case string(Polygon):
			var info resPolygon

			info, err = scrapeInformation[resPolygon](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.PageProps.Collection.Address
			drop.Discord = info.PageProps.Collection.Discord
			drop.Twitter = info.PageProps.Collection.Twitter

		case string(Ethereum):
			var info resEthereum

			info, err = scrapeInformation[resEthereum](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.Address
			drop.Discord = info.Props.PageProps.Collection.Discord
			drop.Twitter = info.Props.PageProps.Collection.Twitter

		case string(Binance):
			var info resBinance

			info, err = scrapeInformation[resBinance](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.Address
			//drop.Twitter = info.Props.PageProps.Collection.Twitter

		case string(Aptos):
			var info resAptos

			info, err = scrapeInformation[resAptos](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.PageProps.Collection.Cm
			drop.Discord = info.PageProps.Collection.Discord
			drop.Twitter = info.PageProps.Collection.Twitter

		case string(Avalanche):
			var info resAvalanche

			info, err = scrapeInformation[resAvalanche](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.Address
			drop.Twitter = info.Props.PageProps.Collection.Twitter
		case string(Fantom):
			var info resFantom

			info, err = scrapeInformation[resFantom](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.Address
			drop.Twitter = info.Props.PageProps.Collection.Twitter
		case string(Stacks):
			var info resStacks

			info, err = scrapeInformation[resStacks](drop.Link)
			if err != nil {
				logger.LogError(moduleName, err)
				return false
			}

			drop.Contract = info.Props.PageProps.Collection.Address

		default:
			logger.LogError(moduleName, errors.New("unknown network"))
			continue
		}

		s.Handler.M.Set(drop.Name, drop.TotalMinted)

		if _, ok := s.Handler.MCopy.Get(drop.Name); ok {
			continue
		}

		drop.populateData(response, i)
		drop.handleSocials()
		if drop.Fraction >= 6 && drop.TotalMinted >= 100 && drop.Contract != "" {
			if err = s.sendDiscordNotification(drop); err != nil {
				logger.LogError(moduleName, err)
			}
		} else {
			logger.LogInfo(moduleName, fmt.Sprintf("🐙 collection progress too low: %s", drop.Name))
		}

		s.Handler.Copy()
	}

	return false
}

// scrapeInformation scrapes information from a collection page and
// decodes it into a struct passed as a generic parameter.
func scrapeInformation[T any](input string) (T, error) {
	var response T

	inputURL, err := url.Parse(input)
	if err != nil {
		return response, err
	}

	resp, err := http.DefaultClient.Do(&http.Request{Method: http.MethodGet, URL: inputURL})
	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		return response, fmt.Errorf("scrape info error: expected %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return response, err
	}

	if err = json.NewDecoder(strings.NewReader(doc.Find("script[id=__NEXT_DATA__]").Text())).Decode(&response); err != nil {
		return response, err
	}

	return response, nil
}

func doRequest(networks []Network) (*http.Response, error) {
	var filterByQuery string
	var buf bytes.Buffer

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
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: "s.launchmynft.io", Path: "/multi_search?x-typesense-api-key=UkN4Vnd3V2JMWWVIRlFNcTJ3dng4VGVtMGtvVGxBcmJJTTFFYS9MNXp1WT1Ha3dueyJmaWx0ZXJfYnkiOiJoaWRkZW46ZmFsc2UiLCJleGNsdWRlX2ZpZWxkcyI6ImhpZGRlbiIsInF1ZXJ5X2J5IjoiY29sbGVjdGlvbk5hbWUsb3duZXIiLCJsaW1pdF9oaXRzIjoyMDAsInNuaXBwZXRfdGhyZXNob2xkIjo1MH0%3D"},
		Body:   io.NopCloser(&buf),
	}

	return http.DefaultClient.Do(req)
}

func (s *Settings) sendDiscordNotification(drop *Release) error {
	return s.Discord.SendNotification(discord.Webhook{
		Username:  s.Discord.ProfileName,
		AvatarUrl: s.Discord.AvatarImage,
		Embeds: []discord.Embed{
			{
				Title:       drop.Name,
				Description: drop.Description,
				Url:         drop.Link,
				Timestamp:   discord.GetTimestamp(),
				Color:       s.Discord.Color,
				Footer: discord.EmbedFooter{
					Text:    s.Discord.FooterText,
					IconUrl: s.Discord.FooterImage,
				},
				Thumbnail: discord.EmbedThumbnail{
					Url: drop.Image,
				},
				Fields: []discord.EmbedFields{
					{
						Name:   "Price",
						Value:  fmt.Sprintf("%.2f %s", drop.Cost, drop.Network),
						Inline: true,
					},
					{
						Name:   "Supply",
						Value:  "`" + strconv.Itoa(drop.TotalMinted) + "/" + strconv.Itoa(drop.Supply) + fmt.Sprintf(" — %.2f%%`", drop.Fraction),
						Inline: true,
					},
					{
						Name:   "Contract",
						Value:  "`" + drop.Contract + "`",
						Inline: false,
					},
					{
						Name:   "Twitter",
						Value:  drop.Twitter,
						Inline: true,
					},
					{
						Name:   "Discord",
						Value:  drop.Discord,
						Inline: true,
					},
					{
						Name:   "Hyperspace",
						Value:  "[Secondary Market]" + "(" + drop.Secondary + ")",
						Inline: true,
					},
				},
			},
		},
	}, s.Discord.Webhook)
}

func (drop *Release) handleSocials() {
	if drop.Twitter == "" {
		drop.Twitter = "no account :("
	} else if drop.Discord == "" {
		drop.Discord = "no server :("
	} else {
		drop.Discord = "[Server](https://discord.gg/" + drop.Discord + ")"
		drop.Twitter = "[Account](https://twitter.com/" + drop.Twitter + ")"
	}
}

func (drop *Release) populateData(res resLaunchMyNFT, i int) {
	drop.Name = res.Results[i].Hits[i].Document.CollectionName
	drop.Description = res.Results[i].Hits[i].Document.Description
	drop.Link = "https://www.launchmynft.io/collections/" + res.Results[i].Hits[i].Document.Owner + "/" + res.Results[i].Hits[i].Document.ID
	drop.Image = res.Results[i].Hits[i].Document.CollectionCoverURL
	drop.Fraction = res.Results[i].Hits[i].Document.FractionMinted * 100
}
