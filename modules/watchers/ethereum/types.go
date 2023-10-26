package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/prometheus"
	"github.com/weeaa/nft/pkg/safemap"
	"math/big"
)

const moduleName = "Ethereum Wallet Watcher"

type Transactions struct {
	Addresses       []common.Address
	Client          *ethclient.Client
	LatestBlockChan chan *types.Block
	Params          safemap.SafeMap[common.Address, *TransactionsParams]
	LatestBlock     uint64
}

// TransactionsParams let you set parameters.
type TransactionsParams struct {
	ABI abi.ABI
}

type Settings struct {
	PromMetrics *prometheus.PromMetrics

	Bot            *bot.Bot
	DB             *db.DB
	Discord        *discord.Client
	Handler        *handler.Handler
	Verbose        bool
	Context        context.Context
	Client         *ethclient.Client
	MonitorParams  MonitorParams
	BalanceWatcher BalanceWatcher
	Transactions   Transactions
}

type (
	Contract string
	Account  string
)

type BalanceWatcher struct {
	Balances map[common.Address]*big.Int
}

type MonitorParams struct {
	BlacklistedTokens []string
}

// DefaultList is a list of valuable traders.
var DefaultList = []string{
	"0x9b26f57f9989C158C66b4A175C9dd5ae128A1F2B",
	"0x036d78c5e87E0aA07Bf61815d1efFe10C9FD5275",
	"0x27db134012676a0542c667c610920e269afe89b9",
	"0x6beEF2B2fE00FDDCa12A8CDA2D4B00435b0ba3b6",
	"0x982E09EBd5Bf6F4f9cCe5d0c84514fb96d91c5F9",
	"0x3f3E2f43f0aC69f30ec38d3E4FEC304bdF330E7A",
	"0x55A9C5180DCAFC98D99d3f3E4B248E9156B12Ac1",
	"0x3635B3d38B971ED37b17E6E1Ac685Af87bc8d930",
	"0x7AbCA3CBC8Aa182D10f742F72E2E8BC68c4a8839",
	"0xBdD95ABE8a7694CCD77143376b0fBea183E6a740",
	"0x71e7b94490837CCAF45F9f6C7c20a3e17bBEb7d3",
	"0x721931508DF2764fD4F70C53Da646Cb8aEd16acE",
	"0x8c40d627EE8a99D07FE9dBF041e11a3381c10697",
	"0xD0322cd77b6223F777b254E7f18FA55D74756B52",
	"0x6C8Ee01F1f8B62E987b3D18F6F28b22a0Ada755f",
	"0x54B174179Ae825Ed630Da40b625Bb3C883CD40ae",
	"0x29e01eC68521FA1c3bd685aA4aDa59FAe1e7C048",
	"0x8C18aA7d789417affA48f59616efBd3E9FFB80c5",
	"0xD9d1C2623fBB4377d9bf29075e610A9B8b4805b4",
	"0x9E29A34dFd3Cb99798E8D88515FEe01f2e4cD5a8",
	"0x9274E50E3922fBc7A3CE99f94EFc32D3BECa6c39",
	"0xb585b60De71E48032e8C19B90896984afc6a660d",
	"0x2329A3412BA3e49be47274a72383839BBf1cdE88",
	"0x6EEf09B526d883F98762a7005FABD2c800DfCA44",
	"0xA7B5cA022774BD02842932e4358DDCbea0CCaADe",
	"0x1BE3edd704be69A7f9E44b7Ad842dCa0757c1816",
	"0xf2E9db3c5D06015833Df31eD3C37172a2B34EE7F",
	"0xA45FC9c051738F135541F97faAE2631cc6167c7C",
	"0xD02d1718C2c62a5c152b27F86469B2bF2b436dC8",
	"0x2ea4815F47D685eD317e8bd243b89bCb26b369Fa",
	"0xE203eFc10f3B3063a34FD6599d754e7F25e2D841",
	"0x63748140C409b490952c37daE5a60715Bf915129",
	"0xB972C02761e51C9C502636c5DBF56635b41c1C26",
	"0x010D591520D0b462F4048Ddb5e591Ed1De3ef1Cb",
	"0xE36a124CaA7Ee0b75A96A934499CE68DaC6D9562",
	"0x83742faddde0b5b2b307ac46f24a1c118d332549",
	"0xeac666c37d94d25fab5977f52a8054427b759533",
	"0x73D4e8632BA37cc9bF9de5045e3E8834F78efa85",
	"0xd8226Dd110c7bA2bcD7A680d9EA5206BaC40F201",
	"0xafD1e0562c91A933f4B40154045cEe71939E95eA",
	"0xDaeD15EB94698CDd18cc2DaE0a5ACdad77E63ddf",
	"0xf75f7f4796874715bb3D2c9989861BCcEa3f305C",
	"0x6C5491665B5aAc18F8e197A26632381AF9732028",
	"0x4c1cd907ceaA5919CF7982679FcE88c58E423dcb",
	"0xf4BdC18c46f742d1f48B84c889371F080cFD709c",
	"0x26D7B4fe67f4601643304b5023b3CAF3A72E8504",
	"0xC2978441F46a76c60e0cd59E986498b75a40572D",
	"0x0B01F1310e7224DAfEd24C3B62d53CeC37d9fAf8",
	"0xC458e1a4eC03C5039fBF38221C54Be4e63731E2A",
	"0x8b3f4eb783270aefAAc9238ac1d165A433C8FbF3",
	"0xf2659a2b2b928a0555bf1596ebf2c30aa4b34a31",
	"0xbde1148eec7b6939f6d6ccf9aaa020f3c0bcc180",
	"0x935745c4539bf41017ae3b63d687a35f0272bc2b",
	"0x336f6beca25aed6bc4b4f9e6ef5b3eb415aeca67",
	"0x0e719677cb5679ff07858f58bfd6fe2a8234863c",
	"0x2d8aed38fc8efd32e3717353e524d1069def4855",
	"0x886478D3cf9581B624CB35b5446693Fc8A58B787",
	"0xD387A6E4e84a6C86bd90C158C6028A58CC8Ac459",
	"0x54BE3a794282C030b15E43aE2bB182E14c409C5e",
	"0xd6a984153aCB6c9E2d788f08C2465a1358BB89A7",
	"0x5ea9681C3Ab9B5739810F8b91aE65EC47de62119",
	"0x7d4823262Bd2c6e4fa78872f2587DDA2A65828Ed",
	"0x0F4BC970e348A061B69D05B7e2E5c13EB687E5e3",
	"0xA34D6cE0e9801562E55C90A3D0C7a1f8B68287Ff",
	"0xc85a9ddeDB6469ff715e8DC3C9616d9459Fa95Fb",
}
