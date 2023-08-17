package exchangeArt

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gagliardetto/solana-go"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/logger"
	"math/big"
	"net/http"
	"time"
)

// ExchangeArt's base API used to monitor new FCFS releases.
const baseURL = "https://api.exchange.art/v2/nfts/created?from=0&sort=listed-oldest&facetsType=collection&limit=10&profileId="
const moduleName = "ExchangeArt"
const DefaultRetryDelay = 2500

// curated list of Artists we used to monitor.
var (
	hyblinxx     = "I2LwzWoHzdcibq3ngiFtumfqmJV2"
	john         = "2Su4KEW92kbhvsv6faONTEDtg9j1"
	adamApe      = "WQh4eWseb7ObpFVGFxTov08HlSr2"
	degenpoet    = "wnzwvrxRolf9qXSzWuC9rN4XbcB2"
	marculinopng = "sP4HMFiwaGXEH4yNLpEU2hdZjAk1"
	lauraEl      = "UEhYxRezIhNe2TNgXOZXGe2iOD02"
	RESIST       = "ZYclOZWnpLUo3IlHgmWlJKYBZtf2"
	purple       = "BEfp6MTyacR0PFXq20ohfbBgRqx1"
	solPlayBoy   = "kKcu8o7TdkWotU1rw6x5s2r0typ1"
	flagMonkez   = "0C2v55yVrjbG4fP07yjyZndDqUm2"
	scum         = "6JDKTmzmFQaGWecWty07xQWQutp1"
	zen0         = "Dt5i9PQN3ocqYfhmk3jDyBdZCw63"
	trevElViz    = "pGFZfmUNDGgSarDiI6MhBIOFymJ3"
)

var DefaultList = []string{
	hyblinxx,
	john,
	adamApe,
	degenpoet,
	marculinopng,
	lauraEl,
	RESIST,
	purple,
	solPlayBoy,
	flagMonkez,
	zen0,
	scum,
	trevElViz,
}

// Monitor monitors newest releases of an artist.
func Monitor(client discord.Client, artists []string, monitor1Spl bool, retryDelay time.Duration) {

	log.Info("ExchangeArt Monitor Started")

	h := handler.New()

	defer func() {
		if r := recover(); r != nil {
			Monitor(client, artists, monitor1Spl, retryDelay)
			return
		}
	}()

	for _, artistURL := range artists {
		go func(artistUrl string) {
			for {
				ea := discord.ExchangeArtWebhook{}

				resp, err := http.Get(baseURL + artistUrl)
				if err != nil {
					log.Errorf("exchangeArt.Monitor: ERROR Client [%w]", err)
					time.Sleep(retryDelay * time.Millisecond)
					continue
				}

				if resp.StatusCode != 200 {
					log.Error("exchangeArt.Monitor: Invalid Response Code", "respStatus", resp.Status)
					time.Sleep(retryDelay * time.Millisecond)
					continue
				}

				res := ResponseExchangeArt{}
				err = json.NewDecoder(resp.Body).Decode(&res)
				if err != nil {
					log.Error(err)
					continue
				}

				err = resp.Body.Close()
				if err != nil {
					log.Error(err)
					continue
				}

				if res.ContractGroups[0].Mint.Market == "secondary" || res.ContractGroups[0].Mint.MetadataAccount.PrimarySaleHappened {
					continue
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
				var lptsOnAccount = new(big.Float).SetUint64(price)
				var solBalance = new(big.Float).Quo(lptsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
				ea.Price = solBalance.Text('f', 10)

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
					if monitor1Spl {
						ea.ToSend = true
						ea.MintLink = "https://exchange.art/single/" + ea.CMID
					}
				} else {
					ea.MintCap = res.ContractGroups[0].AvailableContracts.EditionSales[0].Data.WalletMintingCap
					ea.MintLink = "https://exchange.art/editions/" + ea.CMID
					ea.ToSend = true
				}

				h.M.Set(ea.Name, ea.Name)
				if h.M.Get(ea.Name) == h.MCopy.Get(ea.Name) {
					continue
				}

				h.M.ForEach(func(k string, v interface{}) {
					h.MCopy.Set(k, v)
				})

				if ea.ToSend {
					log.Info("Release Found", "artist", ea.Artist, "collection", ea.Name)
					if err = client.ExchangeArtNotification(discord.Webhook{
						Username:  "ExchangeArt",
						AvatarUrl: client.AvatarImage,
						Embeds: []discord.Embed{
							{
								Title:       ea.Name,
								Description: ea.Description,
								Thumbnail: discord.EmbedThumbnail{
									Url: ea.Image,
								},
								Color:     client.Color,
								Timestamp: discord.GetTimestamp(),
								Footer: discord.EmbedFooter{
									Text:    client.FooterText,
									IconUrl: client.FooterImage,
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
					}); err != nil {
						logger.LogError(moduleName, err)
					}
				}
			}
		}(artistURL)
	}
}
