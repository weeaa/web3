package utils

import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io"
	"os"
	"strings"
)

type Chain string

var (
	Ethereum   Chain = "ethereum"
	AvalancheC Chain = "avalanche-c"
)

func ReadABI(path string) (abi.ABI, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return abi.ABI{}, err
	}

	return abi.JSON(strings.NewReader(string(file)))
}

func GenerateABI() abi.ABI {
	return abi.ABI{}
}

// GetABI returns the ABI of a contract.
func GetABI(chain Chain, apiKey string) (abi.ABI, error) {
	req := &http.Request{}

	switch chain {
	case Ethereum:
	case AvalancheC:
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("client error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return abi.ABI{}, fmt.Errorf("getting abi invalid resp status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

	}

	var response map[string]string
	if err = json.Unmarshal(body, &response); err != nil {

	}

	return abi.ABI{}, nil
}
