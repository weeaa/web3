package lmnft

import (
	"context"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
	"time"
)

const (
	moduleName        = "LaunchMyNFT"
	DefaultRetryDelay = 2500 * time.Millisecond
)

type Network string

const (
	Solana    Network = "Solana"
	Ethereum  Network = "Ethereum"
	Polygon   Network = "Polygon"
	Binance   Network = "Binance"
	Aptos     Network = "Aptos"
	Avalanche Network = "Avalanche"
	Fantom    Network = "Fantom"
	Stacks    Network = "Stacks"
	Sui       Network = "Sui"
)

type Settings struct {
	Discord *discord.Client
	Handler *handler.Handler
	Context context.Context
	Verbose bool
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

type Release struct {
	Name        string
	Description string
	Image       string
	Link        string
	Contract    string
	Supply      int
	TotalMinted int
	Cost        float64
	Fraction    float64
	Verified    bool
	Twitter     string
	Discord     string
	Secondary   string
	Network     Network
}

type resLaunchMyNFT struct {
	Results []struct {
		FacetCounts []struct {
			Counts []struct {
				Count       int    `json:"count"`
				Highlighted string `json:"highlighted"`
				Value       string `json:"value"`
			} `json:"counts"`
			FieldName string `json:"field_name"`
			Stats     struct {
				TotalValues int `json:"total_values"`
			} `json:"stats,omitempty"`
			Stats0 struct {
				Avg         float64 `json:"avg"`
				Max         float64 `json:"max"`
				Min         float64 `json:"min"`
				Sum         float64 `json:"sum"`
				TotalValues int     `json:"total_values"`
			} `json:"stats,omitempty"`
		} `json:"facet_counts"`
		Found int `json:"found"`
		Hits  []struct {
			Document struct {
				CollectionBannerURL string  `json:"collectionBannerUrl"`
				CollectionCoverURL  string  `json:"collectionCoverUrl"`
				CollectionName      string  `json:"collectionName"`
				Cost                float64 `json:"cost"`
				Deployed            int64   `json:"deployed"`
				Description         string  `json:"description"`
				FractionMinted      float64 `json:"fractionMinted"`
				ID                  string  `json:"id"`
				LastMintedAt        int64   `json:"lastMintedAt"`
				LaunchDate          int64   `json:"launchDate"`
				MaxSupply           int     `json:"maxSupply"`
				Owner               string  `json:"owner"`
				SoldOut             bool    `json:"soldOut"`
				TotalMints          int     `json:"totalMints"`
				TwitterVerified     bool    `json:"twitterVerified"`
				Type                string  `json:"type"`
			} `json:"document,omitempty"`
		} `json:"hits"`
		OutOf         int `json:"out_of"`
		Page          int `json:"page"`
		RequestParams struct {
			CollectionName string `json:"collection_name"`
			PerPage        int    `json:"per_page"`
			Q              string `json:"q"`
		} `json:"request_params"`
		SearchCutoff bool `json:"search_cutoff"`
		SearchTimeMs int  `json:"search_time_ms"`
	} `json:"results"`
}

// todo: add all possible responses
type resSolana struct {
	Props struct {
		PageProps struct {
			Collection struct {
				CollectionName           string      `json:"collectionName"`
				HeartCount               int         `json:"heartCount"`
				MintedLast30Mins         bool        `json:"mintedLast30mins"`
				Twitter                  string      `json:"twitter"`
				Version                  int         `json:"version"`
				Tx                       string      `json:"tx"`
				LaunchDate               int64       `json:"launchDate"`
				MetadataCID              string      `json:"metadataCID"`
				Hidden                   bool        `json:"hidden"`
				Cost                     string      `json:"cost"`
				TotalMints               int         `json:"totalMints"`
				HasSoldEnough            bool        `json:"hasSoldEnough"`
				LaunchLater              bool        `json:"launchLater"`
				NewCandyMachineAccountID string      `json:"newCandyMachineAccountId"`
				CollectionCoverURL       string      `json:"collectionCoverUrl"`
				LastMintedAt             int64       `json:"lastMintedAt"`
				MeSymbol                 string      `json:"meSymbol"`
				SoldOut                  bool        `json:"soldOut"`
				Description              string      `json:"description"`
				FractionMinted           float64     `json:"fractionMinted"`
				Discord                  string      `json:"discord"`
				HiddenMetadataCID        interface{} `json:"hiddenMetadataCID"`
				Deployed                 int64       `json:"deployed"`
				TwitterVerified          bool        `json:"twitterVerified"`
				MaxSupply                int         `json:"maxSupply"`
				Type                     string      `json:"type"`
				Owner                    string      `json:"owner"`
			} `json:"collection"`
		} `json:"pageProps"`
		NSsp bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildID      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}

type resEthereum struct {
	Props struct {
		PageProps struct {
			Collection struct {
				Owner         string      `json:"owner"`
				WhitelistCost interface{} `json:"whitelistCost"`
				Address       string      `json:"address"`
				Cost          string      `json:"cost"`
				Hidden        bool        `json:"hidden"`
				Abi           []struct {
					Inputs []struct {
						Indexed      bool   `json:"indexed,omitempty"`
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"inputs"`
					StateMutability string `json:"stateMutability,omitempty"`
					Type            string `json:"type"`
					Name            string `json:"name,omitempty"`
					Anonymous       bool   `json:"anonymous,omitempty"`
					Outputs         []struct {
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"outputs,omitempty"`
				} `json:"abi"`
				Deployed            int64       `json:"deployed"`
				Whitelist           interface{} `json:"whitelist"`
				LaunchLater         bool        `json:"launchLater"`
				Type                string      `json:"type"`
				SoldOut             bool        `json:"soldOut"`
				Version             int         `json:"version"`
				TransactionHash     string      `json:"transactionHash"`
				CollectionName      string      `json:"collectionName"`
				IsWhitelist         bool        `json:"isWhitelist"`
				ContractName        string      `json:"contractName"`
				MetadataCID         string      `json:"metadataCID"`
				HasWhitelistMinted  interface{} `json:"hasWhitelistMinted"`
				ChainId             int         `json:"chainId"`
				MaxMints            int         `json:"maxMints"`
				MaxSupply           int         `json:"maxSupply"`
				Twitter             string      `json:"twitter"`
				TwitterVerified     bool        `json:"twitterVerified"`
				CollectionBannerUrl string      `json:"collectionBannerUrl"`
				CollectionCoverUrl  string      `json:"collectionCoverUrl"`
				Description         string      `json:"description"`
				Discord             string      `json:"discord"`
				HasSoldEnough       bool        `json:"hasSoldEnough"`
				HeartCount          int         `json:"heartCount"`
				FractionMinted      float64     `json:"fractionMinted"`
				LastMintedAt        int64       `json:"lastMintedAt"`
				TotalMints          int         `json:"totalMints"`
				MintedLast30Mins    bool        `json:"mintedLast30mins"`
				StartTime           interface{} `json:"startTime"`
				LaunchDate          interface{} `json:"launchDate"`
			} `json:"collection"`
			DynamicCtx interface{} `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildId      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	DynamicIds   []int         `json:"dynamicIds"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}

type resAptos struct {
	PageProps struct {
		Collection struct {
			Owner               string `json:"owner"`
			Cost                string `json:"cost"`
			Hidden              bool   `json:"hidden"`
			Cm                  string `json:"cm"`
			MaxSupply           int    `json:"maxSupply"`
			LaunchDate          int64  `json:"launchDate"`
			Type                string `json:"type"`
			MetadataCID         string `json:"metadataCID"`
			CollectionName      string `json:"collectionName"`
			Deployed            int64  `json:"deployed"`
			Description         string `json:"description"`
			CollectionCoverURL  string `json:"collectionCoverUrl"`
			CollectionBannerURL string `json:"collectionBannerUrl"`
			Discord             string `json:"discord"`
			HasSoldEnough       bool   `json:"hasSoldEnough"`
			FractionMinted      int    `json:"fractionMinted"`
			SoldOut             bool   `json:"soldOut"`
			TotalMints          int    `json:"totalMints"`
			LastMintedAt        int64  `json:"lastMintedAt"`
			MintedLast30Mins    bool   `json:"mintedLast30mins"`
			Pos                 int    `json:"pos"`
			Featured            bool   `json:"featured"`
			Twitter             string `json:"twitter"`
			TwitterVerified     bool   `json:"twitterVerified"`
			HeartCount          int    `json:"heartCount"`
			StartTime           any    `json:"startTime"`
		} `json:"collection"`
		DynamicCtx any `json:"dynamicCtx"`
	} `json:"pageProps"`
	NSsp bool `json:"__N_SSP"`
}

type resPolygon struct {
	PageProps struct {
		Collection struct {
			WhitelistCost string `json:"whitelistCost"`
			Owner         string `json:"owner"`
			Address       string `json:"address"`
			Hidden        bool   `json:"hidden"`
			Abi           []struct {
				Inputs          []any  `json:"inputs"`
				StateMutability string `json:"stateMutability,omitempty"`
				Type            string `json:"type"`
				Name            string `json:"name,omitempty"`
				Anonymous       bool   `json:"anonymous,omitempty"`
				Outputs         []any  `json:"outputs,omitempty"`
			} `json:"abi"`
			Type                string  `json:"type"`
			Version             int     `json:"version"`
			TransactionHash     string  `json:"transactionHash"`
			CollectionName      string  `json:"collectionName"`
			ChainID             int     `json:"chainId"`
			ContractName        string  `json:"contractName"`
			MaxSupply           int     `json:"maxSupply"`
			HasWhitelistMinted  []any   `json:"hasWhitelistMinted"`
			Deployed            int64   `json:"deployed"`
			CollectionBannerURL string  `json:"collectionBannerUrl"`
			CollectionCoverURL  string  `json:"collectionCoverUrl"`
			LaunchDate          int64   `json:"launchDate"`
			LaunchLater         bool    `json:"launchLater"`
			MetadataCID         string  `json:"metadataCID"`
			SoldOut             bool    `json:"soldOut"`
			WlMaxMints          int     `json:"wlMaxMints"`
			HasSoldEnough       bool    `json:"hasSoldEnough"`
			IsWhitelist         bool    `json:"isWhitelist"`
			Whitelist           []any   `json:"whitelist"`
			MaxMints            int     `json:"maxMints"`
			Description         string  `json:"description"`
			Twitter             string  `json:"twitter"`
			TwitterVerified     bool    `json:"twitterVerified"`
			Discord             string  `json:"discord"`
			Cost                string  `json:"cost"`
			HeartCount          int     `json:"heartCount"`
			MintedLast30Mins    bool    `json:"mintedLast30mins"`
			FractionMinted      float64 `json:"fractionMinted"`
			LastMintedAt        int64   `json:"lastMintedAt"`
			TotalMints          int     `json:"totalMints"`
			StartTime           any     `json:"startTime"`
		} `json:"collection"`
		DynamicCtx any `json:"dynamicCtx"`
	} `json:"pageProps"`
	NSsp bool `json:"__N_SSP"`
}

type resAvalanche struct {
	Props struct {
		PageProps struct {
			Collection struct {
				ContractHash  string `json:"contractHash"`
				Type          string `json:"type"`
				Owner         string `json:"owner"`
				WhitelistCost any    `json:"whitelistCost"`
				Address       string `json:"address"`
				Hidden        bool   `json:"hidden"`
				Abi           []struct {
					Inputs          []any  `json:"inputs"`
					StateMutability string `json:"stateMutability,omitempty"`
					Type            string `json:"type"`
					Name            string `json:"name,omitempty"`
					Anonymous       bool   `json:"anonymous,omitempty"`
					Outputs         []any  `json:"outputs,omitempty"`
				} `json:"abi"`
				Deployed            int64   `json:"deployed"`
				Whitelist           any     `json:"whitelist"`
				LaunchLater         bool    `json:"launchLater"`
				Version             int     `json:"version"`
				TransactionHash     string  `json:"transactionHash"`
				CollectionName      string  `json:"collectionName"`
				WlMaxMints          any     `json:"wlMaxMints"`
				IsWhitelist         bool    `json:"isWhitelist"`
				ChainID             int     `json:"chainId"`
				ContractName        string  `json:"contractName"`
				MaxSupply           int     `json:"maxSupply"`
				MetadataCID         string  `json:"metadataCID"`
				MaxMints            string  `json:"maxMints"`
				EnforceRoyalties    bool    `json:"enforceRoyalties"`
				HasWhitelistMinted  any     `json:"hasWhitelistMinted"`
				CollectionBannerURL string  `json:"collectionBannerUrl"`
				CollectionCoverURL  string  `json:"collectionCoverUrl"`
				Twitter             string  `json:"twitter"`
				TwitterVerified     bool    `json:"twitterVerified"`
				HeartCount          int     `json:"heartCount"`
				SoldOut             bool    `json:"soldOut"`
				HasSoldEnough       bool    `json:"hasSoldEnough"`
				FractionMinted      float64 `json:"fractionMinted"`
				LastMintedAt        int64   `json:"lastMintedAt"`
				TotalMints          int     `json:"totalMints"`
				MintedLast30Mins    bool    `json:"mintedLast30mins"`
				Cost                string  `json:"cost"`
				StartTime           any     `json:"startTime"`
				LaunchDate          any     `json:"launchDate"`
			} `json:"collection"`
			DynamicCtx any `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSsp bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildID      string `json:"buildId"`
	IsFallback   bool   `json:"isFallback"`
	DynamicIds   []int  `json:"dynamicIds"`
	Gssp         bool   `json:"gssp"`
	ScriptLoader []any  `json:"scriptLoader"`
}

type resStacks struct {
	Props struct {
		PageProps struct {
			Collection struct {
				Owner           string      `json:"owner"`
				Address         string      `json:"address"`
				Cost            string      `json:"cost"`
				Hidden          bool        `json:"hidden"`
				ContractName    string      `json:"contractName"`
				MaxSupply       int         `json:"maxSupply"`
				LaunchLater     bool        `json:"launchLater"`
				MetadataCID     string      `json:"metadataCID"`
				Type            string      `json:"type"`
				TransactionHash string      `json:"transactionHash"`
				CollectionName  string      `json:"collectionName"`
				Deployed        int64       `json:"deployed"`
				FractionMinted  int         `json:"fractionMinted"`
				LastMintedAt    int64       `json:"lastMintedAt"`
				SoldOut         bool        `json:"soldOut"`
				TotalMints      int         `json:"totalMints"`
				StartTime       interface{} `json:"startTime"`
				LaunchDate      interface{} `json:"launchDate"`
			} `json:"collection"`
			DynamicCtx interface{} `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildId      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	DynamicIds   []int         `json:"dynamicIds"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}

type resSui struct {
	Props struct {
		PageProps struct {
			Collection struct {
				Owner               string      `json:"owner"`
				ImageCID            string      `json:"imageCID"`
				Creator             string      `json:"creator"`
				Cost                string      `json:"cost"`
				Hidden              bool        `json:"hidden"`
				Module              string      `json:"module"`
				Deployed            int64       `json:"deployed"`
				Cm                  string      `json:"cm"`
				Type                string      `json:"type"`
				Supply              int         `json:"supply"`
				CollectionName      string      `json:"collectionName"`
				MaxSupply           int         `json:"maxSupply"`
				NftDescription      string      `json:"nftDescription"`
				SoldOut             bool        `json:"soldOut"`
				Twitter             string      `json:"twitter"`
				TwitterVerified     bool        `json:"twitterVerified"`
				CollectionBannerUrl string      `json:"collectionBannerUrl"`
				CollectionCoverUrl  string      `json:"collectionCoverUrl"`
				Discord             string      `json:"discord"`
				Description         string      `json:"description"`
				HasSoldEnough       bool        `json:"hasSoldEnough"`
				FractionMinted      float64     `json:"fractionMinted"`
				LastMintedAt        int64       `json:"lastMintedAt"`
				TotalMints          int         `json:"totalMints"`
				MintedLast30Mins    bool        `json:"mintedLast30mins"`
				StartTime           interface{} `json:"startTime"`
				LaunchDate          interface{} `json:"launchDate"`
			} `json:"collection"`
			DynamicCtx interface{} `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildId      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	DynamicIds   []int         `json:"dynamicIds"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}

type resBinance struct {
	Props struct {
		PageProps struct {
			Collection struct {
				ContractHash string `json:"contractHash"`
				Type         string `json:"type"`
				Owner        string `json:"owner"`
				Address      string `json:"address"`
				Cost         string `json:"cost"`
				Hidden       bool   `json:"hidden"`
				RevealLater  bool   `json:"revealLater"`
				Abi          []struct {
					Inputs []struct {
						Indexed      bool   `json:"indexed,omitempty"`
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"inputs"`
					StateMutability string `json:"stateMutability,omitempty"`
					Type            string `json:"type"`
					Name            string `json:"name,omitempty"`
					Anonymous       bool   `json:"anonymous,omitempty"`
					Outputs         []struct {
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"outputs,omitempty"`
				} `json:"abi"`
				Deployed         int64       `json:"deployed"`
				HeartCount       int         `json:"heartCount"`
				Version          int         `json:"version"`
				TransactionHash  string      `json:"transactionHash"`
				CollectionName   string      `json:"collectionName"`
				Immutable        bool        `json:"immutable"`
				ChainId          int         `json:"chainId"`
				ContractName     string      `json:"contractName"`
				MaxSupply        int         `json:"maxSupply"`
				MetadataCID      string      `json:"metadataCID"`
				MaxMints         string      `json:"maxMints"`
				EnforceRoyalties bool        `json:"enforceRoyalties"`
				FractionMinted   int         `json:"fractionMinted"`
				LastMintedAt     int64       `json:"lastMintedAt"`
				SoldOut          bool        `json:"soldOut"`
				HasSoldEnough    bool        `json:"hasSoldEnough"`
				TotalMints       int         `json:"totalMints"`
				MintedLast30Mins bool        `json:"mintedLast30mins"`
				StartTime        interface{} `json:"startTime"`
				LaunchDate       interface{} `json:"launchDate"`
			} `json:"collection"`
			DynamicCtx interface{} `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildId      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	DynamicIds   []int         `json:"dynamicIds"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}

type resFantom struct {
	Props struct {
		PageProps struct {
			Collection struct {
				ContractHash string `json:"contractHash"`
				Type         string `json:"type"`
				Owner        string `json:"owner"`
				Address      string `json:"address"`
				Cost         string `json:"cost"`
				Hidden       bool   `json:"hidden"`
				RevealLater  bool   `json:"revealLater"`
				Abi          []struct {
					Inputs []struct {
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Indexed      bool   `json:"indexed,omitempty"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"inputs"`
					StateMutability string `json:"stateMutability,omitempty"`
					Type            string `json:"type"`
					Name            string `json:"name,omitempty"`
					Anonymous       bool   `json:"anonymous,omitempty"`
					Outputs         []struct {
						Name         string `json:"name"`
						InternalType string `json:"internalType"`
						Type         string `json:"type"`
						Components   []struct {
							Name         string `json:"name"`
							InternalType string `json:"internalType"`
							Type         string `json:"type"`
						} `json:"components,omitempty"`
					} `json:"outputs,omitempty"`
				} `json:"abi"`
				Deployed            int64       `json:"deployed"`
				HeartCount          int         `json:"heartCount"`
				Version             int         `json:"version"`
				TransactionHash     string      `json:"transactionHash"`
				CollectionName      string      `json:"collectionName"`
				Immutable           bool        `json:"immutable"`
				ChainId             int         `json:"chainId"`
				ContractName        string      `json:"contractName"`
				MaxSupply           int         `json:"maxSupply"`
				MetadataCID         string      `json:"metadataCID"`
				MaxMints            string      `json:"maxMints"`
				EnforceRoyalties    bool        `json:"enforceRoyalties"`
				Twitter             string      `json:"twitter"`
				TwitterVerified     bool        `json:"twitterVerified"`
				CollectionBannerUrl string      `json:"collectionBannerUrl"`
				CollectionCoverUrl  string      `json:"collectionCoverUrl"`
				Description         string      `json:"description"`
				StartTime           interface{} `json:"startTime"`
				LaunchDate          interface{} `json:"launchDate"`
				LastMintedAt        interface{} `json:"lastMintedAt"`
			} `json:"collection"`
			DynamicCtx interface{} `json:"dynamicCtx"`
		} `json:"pageProps"`
		NSSP bool `json:"__N_SSP"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Userid       string `json:"userid"`
		Collectionid string `json:"collectionid"`
	} `json:"query"`
	BuildId      string        `json:"buildId"`
	IsFallback   bool          `json:"isFallback"`
	DynamicIds   []int         `json:"dynamicIds"`
	Gssp         bool          `json:"gssp"`
	ScriptLoader []interface{} `json:"scriptLoader"`
}
