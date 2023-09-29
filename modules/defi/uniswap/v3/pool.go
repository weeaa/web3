package v3

import (
	"errors"
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func (s *SwapInstance) getPoolAddress(token0, token1 common.Address, fee *big.Int) (common.Address, error) {
	poolAddr, err := s.uniswapFactory.GetPool(nil, token0, token1, fee)
	if err != nil {
		return common.Address{}, err
	}

	if poolAddr == (common.Address{}) {
		return common.Address{}, errors.New("pool does not exist")
	}

	return poolAddr, nil
}

func (s *SwapInstance) constructV3Pool(token0, token1 *coreEntities.Token, poolFee uint64) (*entities.Pool, error) {
	poolAddress, err := s.getPoolAddress(token0.Address, token1.Address, new(big.Int).SetUint64(poolFee))
	if err != nil {
		return nil, err
	}

	contractPool, err := contract.NewUniswapv3Pool(poolAddress, s.client)
	if err != nil {
		return nil, err
	}

	liquidity, err := contractPool.Liquidity(nil)
	if err != nil {
		return nil, err
	}

	slot0, err := contractPool.Slot0(nil)
	if err != nil {
		return nil, err
	}

	pooltick, err := contractPool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return nil, err
	}

	feeAmount := constants.FeeAmount(poolFee)
	ticks := []entities.Tick{
		{
			Index: entities.NearestUsableTick(sdkutils.MinTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
		{
			Index: entities.NearestUsableTick(sdkutils.MaxTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
	}

	// create tick data provider
	p, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	if err != nil {
		return nil, err
	}

	return entities.NewPool(token0, token1, constants.FeeAmount(poolFee),
		slot0.SqrtPriceX96, liquidity, int(slot0.Tick.Int64()), p)
}
