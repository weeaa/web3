package v3

import (
	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/pkg/utils"
	"math/big"
)

const (
	ContractV3Factory            = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
	ContractV3SwapRouterV1       = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	ContractV3SwapRouterV2       = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"
	ContractV3NFTPositionManager = "0xC36442b4a4522E871399CD717aBDD847Ab11FE88"
	ContractV3Quoter             = "0xb27308f9F90D607463bb33eA1BeBb41C27CE5AB6"
)

var (
	maxSlippage = coreEntities.NewPercent(big.NewInt(5), big.NewInt(1000))
	minSlippage = coreEntities.NewPercent(big.NewInt(25), big.NewInt(100))
)

type SwapInstance struct {
	Wallet         utils.Wallet
	client         *ethclient.Client
	uniswapFactory *contract.Uniswapv3Factory
	Token0         Token
	Token1         Token
	IsMEVEnabled   bool
}

type Token struct {
	Name    string
	Symbol  string
	Address string
	Hex     common.Hash
	Price   string
	Status  string
}
