package exchangeArt

import (
	"context"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
	"time"
)

const (
	moduleName = "ExchangeArt"
	baseURL    = "https://api.exchange.art/v2/nfts/created?from=0&sort=listed-oldest&facetsType=collection&limit=10&profileId=" // ExchangeArt's base API used to monitor new FCFS releases.
)

type Settings struct {
	Discord        *discord.Client
	Handler        *handler.Handler
	Context        context.Context
	Verbose        bool
	Monitor1Supply bool
	RetryDelay     time.Duration
}

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

// DefaultList is a curated list of Artists we used to monitor.
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

type ResponseExchangeArt struct {
	TotalGroups    int `json:"totalGroups"`
	ContractGroups []struct {
		Mint struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			MetadataAccount struct {
				UpdateAuthority      string `json:"updateAuthority"`
				Mint                 string `json:"mint"`
				Name                 string `json:"name"`
				Symbol               string `json:"symbol"`
				URI                  string `json:"uri"`
				SellerFeeBasisPoints int    `json:"sellerFeeBasisPoints"`
				IsMutable            bool   `json:"isMutable"`
				PrimarySaleHappened  bool   `json:"primarySaleHappened"`
				Creators             []struct {
					Share    int    `json:"share"`
					Address  string `json:"address"`
					Verified int    `json:"verified"`
				} `json:"creators"`
			} `json:"metadataAccount"`
			OffChainJSONMetadata struct {
				Name        string `json:"name"`
				Symbol      string `json:"symbol"`
				Description string `json:"description"`
				Image       string `json:"image"`
				Attributes  []struct {
					TraitType string `json:"trait_type"`
					Value     string `json:"value"`
				} `json:"attributes"`
				Files []struct {
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"files"`
			} `json:"offChainJsonMetadata"`
			Market string `json:"market"`
			Symbol string `json:"symbol"`
			Brand  struct {
				ID            string `json:"id"`
				Name          string `json:"name"`
				Description   string `json:"description"`
				BannerPath    string `json:"bannerPath"`
				ThumbnailPath string `json:"thumbnailPath"`
			} `json:"brand"`
			ArtistProfile struct {
				Md struct {
					DisplayName string `json:"displayName"`
					Slug        string `json:"slug"`
					BrandID     string `json:"brandId"`
					CreatorID   string `json:"creatorId"`
				} `json:"md"`
				Assets struct {
					Banner    string `json:"banner"`
					Thumbnail string `json:"thumbnail"`
				} `json:"assets"`
				Twitter struct {
					Handle       string `json:"handle"`
					ProfileImage string `json:"profileImage"`
				} `json:"twitter"`
				PrintServicesEnabled bool `json:"printServicesEnabled"`
			} `json:"artistProfile"`
			Collection struct {
				ID                string   `json:"id"`
				Name              string   `json:"name"`
				ThumbnailPath     string   `json:"thumbnailPath"`
				IsCertified       bool     `json:"isCertified"`
				IsCurated         bool     `json:"isCurated"`
				IsOneOfOne        bool     `json:"isOneOfOne"`
				IsNsfw            bool     `json:"isNsfw"`
				IsDiscounted      bool     `json:"isDiscounted"`
				Description       string   `json:"description"`
				PrimaryCategory   string   `json:"primaryCategory"`
				SecondaryCategory string   `json:"secondaryCategory"`
				TertiaryCategory  string   `json:"tertiaryCategory"`
				Tags              []string `json:"tags"`
				BannerPath        string   `json:"bannerPath"`
			} `json:"collection"`
			Description string `json:"description"`
			Image       string `json:"image"`
			Attributes  []struct {
				TraitType string `json:"trait_type"`
				Value     string `json:"value"`
			} `json:"attributes"`
			MasterEditionKey        string `json:"masterEditionKey"`
			MasterEditionAccountPDA string `json:"masterEditionAccountPDA"`
			IsMasterEdition         bool   `json:"isMasterEdition"`
			MasterEditionAccount    struct {
				CurrentSupply int `json:"currentSupply"`
				MaxSupply     int `json:"maxSupply"`
			} `json:"masterEditionAccount"`
			Stats struct {
				NumBuynows    int           `json:"numBuynows"`
				NumAuctions   int           `json:"numAuctions"`
				NumOffers     int           `json:"numOffers"`
				HighestOffers []interface{} `json:"highestOffers"`
				LastSale      struct {
					Currency  string      `json:"currency"`
					Amount    interface{} `json:"amount"`
					AmountUsd float64     `json:"amountUsd"`
				} `json:"lastSale"`
				LowestBuynows []struct {
					Currency string `json:"currency"`
					Amount   int    `json:"amount"`
				} `json:"lowestBuynows"`
			} `json:"stats"`
		} `json:"mint"`
		AvailableContracts struct {
			EditionSales []struct {
				Data struct {
					BlockTime           int    `json:"blockTime"`
					Start               int    `json:"start"`
					Price               int64  `json:"price"`
					Currency            string `json:"currency"`
					PricingType         string `json:"pricingType"`
					SaleType            string `json:"saleType"`
					EditionsMintedSoFar int    `json:"editionsMintedSoFar"`
					Version             int    `json:"version"`
					WalletMintingCap    int    `json:"walletMintingCap"`
					PreSaleWindow       int    `json:"preSaleWindow"`
					RoyaltyProtection   bool   `json:"royaltyProtection"`
				} `json:"data"`
				Keys struct {
					Initializer        string `json:"initializer"`
					MintKey            string `json:"mintKey"`
					DepositAccount     string `json:"depositAccount"`
					SaleAccount        string `json:"saleAccount"`
					AddressLookupTable string `json:"addressLookupTable"`
				} `json:"keys"`
				Type string `json:"type"`
			} `json:"editionSales"`
			Auctions []interface{} `json:"auctions"`
			Listings []interface{} `json:"listings"`
			Offers   []interface{} `json:"offers"`
		} `json:"availableContracts"`
	} `json:"contractGroups"`
	Facets struct {
		Md struct {
			Brands []struct {
				Name string `json:"name"`
				Nfts int    `json:"nfts"`
			} `json:"brands"`
			Collections []struct {
				Name string `json:"name"`
				Nfts int    `json:"nfts"`
			} `json:"collections"`
			Categories []struct {
				Name string `json:"name"`
				Nfts int    `json:"nfts"`
			} `json:"categories"`
			Tags []struct {
				Name string `json:"name"`
				Nfts int    `json:"nfts"`
			} `json:"tags"`
		} `json:"md"`
	} `json:"facets"`
}
