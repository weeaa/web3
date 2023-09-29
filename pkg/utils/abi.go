package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"os"
	"strings"
)

func ReadABI(path string) (abi.ABI, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return abi.ABI{}, err
	}

	return abi.JSON(strings.NewReader(string(file)))
}

func GetABI() {

}
