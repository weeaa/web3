package friendtech

import (
	"github.com/shopspring/decimal"
	"math/big"
)

func (s *Sniper) calculateSharePrice(maxUserInput float64, address string) float64 {
	info, err := s.GetUserInformation(address)

	decVal, err := decimal.NewFromString(info.DisplayPrice)
	if err != nil {

	}

	_ = decVal

	return 0
}

func stringToBigInt() *big.Int {
	return nil
}
