package abi_utils

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"net/url"
	"os"
	"strings"
)

type Chain string

var (
	Ethereum   Chain = "ethereum"
	AvalancheC Chain = "avalanche-c"
)

func ReadABI(filePath string) (abi.ABI, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return abi.ABI{}, err
	}

	return abi.JSON(strings.NewReader(string(file)))
}

func GenerateABI() abi.ABI {
	return abi.ABI{}
}

// GetABI returns the ABI of a contract utilizing Etherscan platforms.
func GetABI(chain Chain, apiKey string) (abi.ABI, error) {
	req := &http.Request{}

	switch chain {
	case Ethereum:
		req.URL = &url.URL{
			Scheme: "https",
			Host:   "",
			Path:   "",
		}
	case AvalancheC:
		req.URL = &url.URL{
			Scheme: "https",
			Host:   "",
			Path:   "",
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("client error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return abi.ABI{}, fmt.Errorf("getting abi invalid resp status: %s", resp.Status)
	}

	return abi.JSON(resp.Body)
}
