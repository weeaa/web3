package opensea

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/foundVanting/opensea-stream-go/opensea"
	"github.com/foundVanting/opensea-stream-go/types"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"math/big"
	"net/http"
	"net/url"
)

func NewClient(discordClient *discord.Client, verbose bool, openseaApiKey string, openseaFloorPCt float64) *Settings {
	client := opensea.NewStreamClient(types.MAINNET, openseaApiKey, nil, func(err error) {
		logger.LogError("OpenSea", err)
		return
	})
	if err := client.Connect(); err != nil {
		logger.LogError("OpenSea", err)
		return nil
	}

	return &Settings{
		OpenSeaFloorPct: openseaFloorPCt,
		Discord:         discordClient,
		Verbose:         verbose,
		Handler:         handler.New(),
		Context:         context.Background(),
		OpenSeaClient: &Client{
			ApiKey:       openseaApiKey,
			StreamClient: client,
		},
	}
}

func (s *Settings) GetFloor(collectionSlug string) (float64, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "api.opensea.io", Path: fmt.Sprintf("/api/v1/collection/%s/stats", collectionSlug)},
		Header: s.getHeaders(),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != 200 {
		return -1, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var cd CollectionData
	if err = json.Unmarshal(body, &cd); err != nil {
		return -1, err
	}

	if err = resp.Body.Close(); err != nil {
		return -1, err
	}

	return cd.Stats.FloorPrice, nil
}

func checkIfFloorBelowX(floor, pct float64) float64 {
	return floor - (floor * pct / 100)
}

//func (s *Settings) checkIfFloorBelowX(floor, pct float64) float64 { return floor - (floor / pct * 2) }

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).SetPrec(236).SetMode(big.ToNearestEven).Quo(new(big.Float).SetPrec(236).SetMode(big.ToNearestEven).SetInt(wei), big.NewFloat(params.Ether))
}

func (s *Settings) getHeaders() http.Header {
	return http.Header{
		"accept":    {"application/json"},
		"x-api-key": {s.OpenSeaClient.ApiKey},
	}
}
