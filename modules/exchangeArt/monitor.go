package exchangeArt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func NewClient(discordClient *discord.Client, verbose, monitor1Supply bool, retryDelay int) *Settings {
	return &Settings{
		Discord:        discordClient,
		Handler:        handler.New(),
		Verbose:        verbose,
		Context:        context.Background(),
		Monitor1Supply: monitor1Supply,
		RetryDelay:     time.Duration(retryDelay),
	}
}

// Monitor monitors newest releases of an artist from ExchangeArt.
func (s *Settings) StartMonitor(artists []string) {
	logger.LogStartup(moduleName)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogInfo(moduleName, fmt.Sprintf("program panicked! [%v]", r))
				s.StartMonitor(artists)
				return
			}
		}()
		for !s.monitorArtists(artists) {
			select {
			case <-s.Context.Done():
				logger.LogShutDown(moduleName)
				return
			default:
				time.Sleep(s.RetryDelay * time.Millisecond)
				continue
			}
		}
	}()
}

func (s *Settings) monitorArtists(artists []string) bool {
	wg := sync.WaitGroup{}
	for _, artistURL := range artists {
		wg.Add(1)
		go func(artist string) {
			defer wg.Done()
			var ea discord.ExchangeArtWebhook

			req := &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Scheme: "https", Host: "api.exchange.art", Path: "/v2/bff/graphql"},
				Body:   io.NopCloser(strings.NewReader(q(artist))),
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}

			if resp.StatusCode != 200 {
				logger.LogError(moduleName, fmt.Errorf("artist: %s > invalid response status: %s", artist, resp.Status))
				time.Sleep(2500 * time.Millisecond)
				return
			}

			var res ResponseExchangeArt
			if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return
			}

			if err = resp.Body.Close(); err != nil {
				return
			}

			if res.ContractGroups[0].Mint.Market == "secondary" || res.ContractGroups[0].Mint.MetadataAccount.PrimarySaleHappened {
				return
			}

			if len(res.ContractGroups[0].AvailableContracts.EditionSales) > 0 {
				for j := 0; j < len(res.ContractGroups[0].AvailableContracts.EditionSales); j++ {
					if len(res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.SaleType) > 0 {
						SaleType := res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.SaleType
						switch SaleType {
						case "OpenEdition":
							ea.ReleaseType = "Open Supply"
						case "LimitedEdition":
							ea.ReleaseType = "Limited Supply"
						default:
							ea.ReleaseType = "Unknown Release Type"
						}
					}
				}
			}

			var price uint64
			for j := 0; j < len(res.ContractGroups[0].AvailableContracts.EditionSales); j++ {
				price = uint64(res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.Price)
			}

			// converts Lamports to human readable amount in SOL.
			ea.Price = new(big.Float).Quo(new(big.Float).SetUint64(price), new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL)).Text('f', 10)

			ea.Name = res.ContractGroups[0].Mint.Name
			ea.Description = res.ContractGroups[0].Mint.Description
			ea.Image = res.ContractGroups[0].Mint.Image
			ea.Artist = res.ContractGroups[0].Mint.Brand.Name
			ea.CMID = res.ContractGroups[0].Mint.ID
			ea.Supply = fmt.Sprint(res.ContractGroups[0].Mint.MasterEditionAccount.MaxSupply)
			ea.Minted = res.ContractGroups[0].Mint.MasterEditionAccount.CurrentSupply

			ea.ToSend = false
			ea.Edition = len(res.ContractGroups[0].AvailableContracts.EditionSales)
			if ea.Edition == 0 {
				if s.Monitor1Supply {
					ea.ToSend = true
					ea.MintLink = "https://exchange.art/single/" + ea.CMID
				}
			} else {
				ea.MintCap = res.ContractGroups[0].AvailableContracts.EditionSales[0].Data.WalletMintingCap
				ea.MintLink = "https://exchange.art/editions/" + ea.CMID
				ea.ToSend = true
			}

			s.Handler.M.Set(ea.Name, ea.Artist)
			if _, ok := s.Handler.MCopy.Get(ea.Name); ok {
				return
			}

			s.Handler.Copy()

			if ea.ToSend && s.Discord.Webhook != "" {
				if err = s.Discord.SendNotification(discord.Webhook{
					Username:  s.Discord.ProfileName,
					AvatarUrl: s.Discord.AvatarImage,
					Embeds: []discord.Embed{
						{
							Title:       ea.Name,
							Description: ea.Description,
							Thumbnail: discord.EmbedThumbnail{
								Url: ea.Image,
							},
							Color:     s.Discord.Color,
							Timestamp: discord.GetTimestamp(),
							Footer: discord.EmbedFooter{
								Text:    s.Discord.FooterText,
								IconUrl: s.Discord.FooterImage,
							},
							Fields: []discord.EmbedFields{
								{
									Name:   "Supply/Max(wallet)",
									Value:  fmt.Sprintf("`%s/%d`", ea.Supply, ea.MintCap),
									Inline: true,
								},
								{
									Name:   "Release Type",
									Value:  ea.ReleaseType,
									Inline: true,
								},
								{
									Name:   "Artist",
									Value:  ea.Artist,
									Inline: true,
								},
								{
									Name:   "Price",
									Value:  fmt.Sprintf("%s", ea.Price[0:4]),
									Inline: true,
								},
								{
									Name:   "CandyMachine ID",
									Value:  fmt.Sprintf("`%s`", ea.CMID),
									Inline: false,
								}},
						},
					},
				}, s.Discord.Webhook); err != nil {
					logger.LogError(moduleName, err)
				}
			}

			if s.Verbose {
				logger.LogInfo(moduleName, fmt.Sprintf("🧿 new release found: %s | %s", ea.Name, ea.Artist))
			}

		}(artistURL)
	}
	wg.Wait()
	return false
}

func doRequest(artist string) (*http.Response, error) {
	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: "api.exchange.art", Path: "/v2/bff/graphql"},
		Body:   io.NopCloser(strings.NewReader(q(artist))),
		Header: http.Header{
			"authority":          {"api.exchange.art"},
			"accept":             {"application/json, text/plain, */*"},
			"accept-language":    {"en-US,en;q=0.9"},
			"content-type":       {"application/json"},
			"origin":             {"https://exchange.art"},
			"referer":            {"https://exchange.art/"},
			"sec-ch-ua":          {"Not)A;Brand\";v=\"24\", \"Chromium\";v=\"116\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-site"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"},
		},
	}
	return http.DefaultClient.Do(req)
}

func q(artist string) string {
	return fmt.Sprintf("{\"query\":\"\n  \n  fragment SellingAgreementsFragment on SellingAgreements {\n    buyNow {\n      createdAt\n      createdBy\n      currency\n      price\n      priceUsd\n      governingProgram\n      stateAccount\n      escrowAccount\n      beginsAt\n      royaltyProtection\n    }\n    offers {\n      createdAt\n      createdBy\n      currency\n      price\n      priceUsd\n      governingProgram\n      stateAccount\n      escrowAccount\n    }\n    auctions {\n      createdAt\n      createdBy\n      currency\n      price\n      priceUsd\n      governingProgram\n      stateAccount\n      escrowAccount\n      endsAt\n      beginsAt\n      endingPhase\n      endingPhasePercentageFlip\n      extensionWindow\n      minimumIncrement\n      reservePrice\n      reservePriceUsd\n      highestBid\n      highestBidUsd\n      highestBidder\n      numberOfBids\n    }\n    editionSales {\n      createdAt\n      createdBy\n      currency\n      price\n      governingProgram\n      stateAccount\n      escrowAccount\n      preSaleWindow\n      pricingType\n      royaltyProtection\n      saleType\n      beginsAt\n      walletMintingCapacity\n      addressLookupTable\n    }\n  }\n\n  \n  \n  fragment SeriesFragment on Series {\n    id\n    description\n    isCertified\n    isCurated\n    isNsfw\n    isOneOfOne\n    name\n    primaryCategory\n    secondaryCategory\n    tags\n    tertiaryCategory\n    website\n    discord\n    twitter\n    thumbnailPath\n    bannerPath\n  }\n\n  fragment NftWithArtistProfileFragmentAndSeries on Nft {\n    id\n    blockchain\n    seriesIds\n    series {\n      ...SeriesFragment\n    }\n    masterEdition {\n      address\n      supply\n      maxSupply\n      permanentlyEnd\n    }\n    edition {\n      address\n      masterEditionId\n      parent\n      num\n    }\n    mintedAt\n    mintedOnExchange\n    artistProfileId\n    artistProfile {\n      md {\n        displayName\n        slug\n      }\n      assets {\n        thumbnail\n      }\n      twitter {\n        handle\n      }\n    }\n    curated\n    certified\n    discounted\n    nsfw\n    aiGenerated\n    royaltyProtected\n    category\n    metadata {\n      accountAddress\n      name\n      symbol\n      updateAuthority\n      primarySaleHappened\n      isMutable\n      creators {\n        address\n        royaltyBps\n      }\n    }\n    json {\n      uri\n      image\n      description\n      attributes {\n        value\n        traitType\n      }\n      files {\n        uri\n        type\n      }\n    }\n  }\n\n  \n  fragment StockReportFragment on StockReport {\n    totalBuyNowSellingAgreements\n    totalAuctionSellingAgreements\n    totalOfferSellingAgreements\n    lowestBuyNowPriceUsd\n    highestBuyNowPriceUsd\n    lowestAuctionPriceUsd\n    highestAuctionPriceUsd\n    lowestOfferPriceUsd\n    highestOfferPriceUsd\n  }\n\n  \n  fragment NftStatsFragment on NftStats {\n    lastSale {\n      currency\n      amount\n      amountUsd\n    }\n  }\n\n  query stockEntry($input: GetStockDto!) {\n    getStockResponse(input: $input) {\n      results {\n        nft {\n          ...NftWithArtistProfileFragmentAndSeries\n        }\n        sellingAgreements {\n          ...SellingAgreementsFragment\n        }\n        report {\n          ...StockReportFragment\n        }\n        nftStats {\n          ...NftStatsFragment\n        }\n      }\n      total\n    }\n  }\n\",\"variables\":{\"input\":{\"from\":0,\"sort\":\"newest_listed\",\"filters\":{\"currencies\":[],\"nftArtistProfileIds\":[\"VBCLsHr5bmNFuzq7fvRm9n5dlF12\"]},\"limit\":20}}}")
}
