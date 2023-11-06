package sniper

import (
	"github.com/shopspring/decimal"
	"math/big"
)

func functionStr(functionName Function) string {
	return string(functionName)
}

func isSharePriceOk(maxUserInput *big.Float, sharePrice string) bool {

	sharePriceDec, err := decimal.NewFromString(sharePrice)
	if err != nil {
		return false
	}

	maxUserInputDec, err := decimal.NewFromString(maxUserInput.String())
	if err != nil {
		return false
	}

	return sharePriceDec.GreaterThan(maxUserInputDec)
}

func calculateMaxFee(input *big.Int) {

}
