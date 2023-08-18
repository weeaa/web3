package opensea

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/params"
	"github.com/foundVanting/opensea-stream-go/opensea"
	"github.com/foundVanting/opensea-stream-go/types"
	"github.com/weeaa/nft/pkg/logger"
	"io"
	"math/big"
	"net/http"
	"net/url"
)

func NewClient(key string) *Client {
	client := opensea.NewStreamClient(types.MAINNET, key, nil, func(err error) {
		logger.LogError(moduleName, err)
		return
	})
	if err := client.Connect(); err != nil {
		logger.LogError(moduleName, err)
		return nil
	}
	return &Client{ApiKey: key, StreamClient: client}
}

func (c *Client) GetFloor(collectionSlug string) (float64, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https://", Host: "api.opensea.io", Path: fmt.Sprintf("/api/v1/collection/%s/stats", collectionSlug)},
		Header: c.getHeaders(),
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

func weiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func (c *Client) getHeaders() http.Header {
	return http.Header{
		"accept":    {"application/json"},
		"x-api-key": {c.ApiKey},
	}
}
