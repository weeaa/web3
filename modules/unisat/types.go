package unisat

import (
	"context"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
)

const (
	moduleName                              = "Unisat BRC20"
	DefaultPercentageIncreaseBetweenRefresh = 3
)

type Settings struct {
	Discord                          *discord.Client
	Handler                          *handler.Handler
	Context                          context.Context
	Verbose                          bool
	RotateProxyOnBan                 bool
	Client                           tls_client.HttpClient
	ProxyList                        []string
	PercentageIncreaseBetweenRefresh float64
}

type ResTickers struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Height int `json:"height"`
		Total  int `json:"total"`
		Start  int `json:"start"`
		Detail []struct {
			Ticker                 string `json:"ticker"`
			HoldersCount           int    `json:"holdersCount"`
			HistoryCount           int    `json:"historyCount"`
			InscriptionNumber      int    `json:"inscriptionNumber"`
			InscriptionID          string `json:"inscriptionId"`
			Max                    string `json:"max"`
			Limit                  string `json:"limit"`
			Minted                 string `json:"minted"`
			TotalMinted            string `json:"totalMinted"`
			ConfirmedMinted        string `json:"confirmedMinted"`
			ConfirmedMinted1H      string `json:"confirmedMinted1h"`
			ConfirmedMinted24H     string `json:"confirmedMinted24h"`
			MintTimes              int    `json:"mintTimes"`
			Decimal                int    `json:"decimal"`
			Creator                string `json:"creator"`
			Txid                   string `json:"txid"`
			DeployHeight           int    `json:"deployHeight"`
			DeployBlocktime        int    `json:"deployBlocktime"`
			CompleteHeight         int    `json:"completeHeight"`
			CompleteBlocktime      int    `json:"completeBlocktime"`
			InscriptionNumberStart int    `json:"inscriptionNumberStart"`
			InscriptionNumberEnd   int    `json:"inscriptionNumberEnd"`
		} `json:"detail"`
	} `json:"data"`
}

type ResTickerInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Ticker                 string `json:"ticker"`
		HoldersCount           int    `json:"holdersCount"`
		HistoryCount           int    `json:"historyCount"`
		InscriptionNumber      int    `json:"inscriptionNumber"`
		InscriptionID          string `json:"inscriptionId"`
		Max                    string `json:"max"`
		Limit                  string `json:"limit"`
		Minted                 string `json:"minted"`
		TotalMinted            string `json:"totalMinted"`
		ConfirmedMinted        string `json:"confirmedMinted"`
		ConfirmedMinted1H      string `json:"confirmedMinted1h"`
		ConfirmedMinted24H     string `json:"confirmedMinted24h"`
		MintTimes              int    `json:"mintTimes"`
		Decimal                int    `json:"decimal"`
		Creator                string `json:"creator"`
		Txid                   string `json:"txid"`
		DeployHeight           int    `json:"deployHeight"`
		DeployBlocktime        int    `json:"deployBlocktime"`
		CompleteHeight         int    `json:"completeHeight"`
		CompleteBlocktime      int    `json:"completeBlocktime"`
		InscriptionNumberStart int    `json:"inscriptionNumberStart"`
		InscriptionNumberEnd   int    `json:"inscriptionNumberEnd"`
	} `json:"data"`
}

type ResHolders struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Height int `json:"height"`
		Total  int `json:"total"`
		Start  int `json:"start"`
		Detail []struct {
			Address             string `json:"address"`
			OverallBalance      string `json:"overallBalance"`
			TransferableBalance string `json:"transferableBalance"`
			AvailableBalance    string `json:"availableBalance"`
		} `json:"detail"`
	} `json:"data"`
}

type ResFees struct {
	FastestFee  string `json:"fastestFee"`
	HalfHourFee string `json:"halfHourFee"`
	HourFee     string `json:"hourFee"`
	Btcusd      string `json:"BTCUSD"`
}
