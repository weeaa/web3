package exchangeArt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"math/big"
	"net/http"
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
func (s *Settings) StartMonitor(artists *[]string) {
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

func (s *Settings) monitorArtists(artists *[]string) bool {
	wg := sync.WaitGroup{}
	for _, artistURL := range *artists {
		wg.Add(1)
		go func(artist string) {
			defer wg.Done()
			var ea discord.ExchangeArtWebhook

			resp, err := http.Get(baseURL + artist)
			if err != nil {
				return
			}

			if resp.StatusCode != 200 {
				logger.LogError(moduleName, fmt.Errorf("invalid response status: %s", resp.Status))
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

			s.Handler.M.ForEach(func(k string, v interface{}) {
				s.Handler.MCopy.Set(k, v)
			})

			if ea.ToSend {
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
