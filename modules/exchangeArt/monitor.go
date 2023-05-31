package exchangeArt

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gagliardetto/solana-go"
	"math/big"
	"net/http"
	"nft/discord"
	"nft/handler"
	"time"
)

// ExchangeArt's base API used to monitor new FCFS releases.
const baseURL = "https://api.exchange.art/v2/nfts/created?from=0&sort=listed-oldest&facetsType=collection&limit=10&profileId="

const retryDelay = 2500

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

var list = []string{
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

func Monitor(discordWebhook string, artists []string, monitor1Spl bool) {

	log.Info("ExchangeArt Monitor Started")

	h := handler.New()

	for {
		for _, artistURL := range artists {
			t := discord.ExchangeArtWebhook{}

			resp, err := http.Get(baseURL + artistURL)
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

			for i := 0; i < len(res.ContractGroups); i++ {

				// if secondary market (resale) it will keep looping as we only wanted FCFS releases
				// if primary sale happened we don't want to be notified neither
				if res.ContractGroups[i].Mint.Market == "secondary" || res.ContractGroups[i].Mint.MetadataAccount.PrimarySaleHappened {
					continue
				}

				if len(res.ContractGroups[0].AvailableContracts.EditionSales) > 0 {
					for j := 0; j < len(res.ContractGroups[0].AvailableContracts.EditionSales); j++ {
						if len(res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.SaleType) > 0 {
							SaleType := res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.SaleType
							switch SaleType {
							case "OpenEdition":
								t.ReleaseType = "Open Supply"
							case "LimitedEdition":
								t.ReleaseType = "Limited Supply"
							default:
								t.ReleaseType = "Unknown Release Type"
							}
						}
					}
				}

				/*
					currentTime := time.Now().Unix()
					for j := 0; j < len(jsonParser.ContractGroups[0].AvailableContracts.EditionSales); j++ {
						for k := 0; k < len(strconv.Itoa(jsonParser.ContractGroups[0].AvailableContracts.EditionSales[j].Data.Start)); k++ {
							if int(currentTime) > jsonParser.ContractGroups[0].AvailableContracts.EditionSales[j].Data.Start {
								break
							} else {
								t.LiveAt = "<t:" + strconv.Itoa(jsonParser.ContractGroups[0].AvailableContracts.EditionSales[j].Data.Start) + ">"
							}
						}
					}
				*/

				var price uint64
				for j := 0; j < len(res.ContractGroups[0].AvailableContracts.EditionSales); j++ {
					price = uint64(res.ContractGroups[0].AvailableContracts.EditionSales[j].Data.Price)
				}

				// converts Lamports to human readable amount in SOL.
				var lptsOnAccount = new(big.Float).SetUint64(price)
				var solBalance = new(big.Float).Quo(lptsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))
				t.Price = solBalance.Text('f', 10)

				t.Name = res.ContractGroups[0].Mint.Name
				t.Description = res.ContractGroups[0].Mint.Description
				t.Image = res.ContractGroups[0].Mint.Image
				t.Artist = res.ContractGroups[0].Mint.Brand.Name
				t.CMID = res.ContractGroups[0].Mint.ID
				t.Supply = fmt.Sprint(res.ContractGroups[0].Mint.MasterEditionAccount.MaxSupply)
				t.Minted = res.ContractGroups[0].Mint.MasterEditionAccount.CurrentSupply

				t.Edition = len(res.ContractGroups[0].AvailableContracts.EditionSales)
				if t.Edition == 0 { // isOnly 1 supply, we don't want that.
					if monitor1Spl {
						t.ToSend = true
						t.MintLink = "https://exchange.art/single/" + t.CMID
					} else {
						t.ToSend = false
					}
				} else {
					t.MintCap = res.ContractGroups[0].AvailableContracts.EditionSales[0].Data.WalletMintingCap
					t.MintLink = "https://exchange.art/editions/" + t.CMID
					t.ToSend = true
				}

			}

			h.M[t.Name] = t.Name
			if h.M[t.Name] == h.MCopy[t.Name] {
				h.MCopy[t.Name] = h.M[t.Name]
				continue
			}

			for k, v := range h.M {
				h.MCopy[k] = v
			}

			if t.ToSend {
				log.Info("Release Found", "artist", t.Artist)
				err = t.ExchangeArtNotification(discordWebhook)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
