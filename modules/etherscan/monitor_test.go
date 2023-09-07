package etherscan

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

// todo finish unit tests
type MockDiscord struct{}

func (m *MockDiscord) SendNotification(webhook discord.Webhook) error {
	return nil
}

func TestMonitorVerifiedContracts(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "valid response status & body",
			status: http.StatusOK,
		},
		{
			name:   "invalid response status",
			status: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.status)
			w.Write([]byte(html))
		}))

		if test.status == http.StatusTooManyRequests {
			assert.Error(t, fmt.Errorf("too many requests %d", test.status))
		}

		settings := &Settings{
			Handler: handler.New(),
			Context: context.Background(),
		}

		if ok := settings.monitorVerifiedContracts(); ok {
			assert.Error(t, fmt.Errorf(""))
		}

		contract, ok := settings.Handler.M.Get("BEV") // BEV Contract
		if ok {

		}

		if contract == "0xcf73b0d42c2c97219ce5895f311a7aa0fc930c98" {

		}

		server.Close()

		assert.NoError(t, nil)
	}

}

var html = `
<!doctype html>
<html id="html" lang="en">
    <head>
        <title>Ethereum Verified Contracts | Etherscan
</title>
        <style>
            .autolink-safari a {
                pointer-events: none !important;
                text-decoration: none !important;
                color: inherit !important;
            }
        </style>
        <style>
            .deciuybgh6j {
                -ms-flex-pack: center;
                justify-content: center;
            }

            .acil5hxi3zq {
                border-top: 1px solid #e7eaf3 !important;
            }

            body.dark-mode .acil5hxi3zq {
                border-color: #18365b !important;
            }

            @media (min-width: 992px) {
                .bdtffum2sfr {
                    display: block !important;
                }
            }

            @media (min-width: 992px) {
                .c78zi5egt7w {
                    display: inline-block !important;
                }
            }
        </style>
    </head>
    <body id="body" class="d-flex flex-column min-vh-100">
        <section id="masterNoticeBar"></section>
        <section id="masterTopBar" class="sticky-top bg-white border-bottom py-2 d-print-none">
            <div class="container-xxl d-flex align-items-center justify-content-between">
                <div id="ethPrice" class="d-none d-md-flex align-items-center gap-4 w-100 fs-sm">
                    <div class="text-muted">
                        <span class="text-muted">ETH Price:</span>
                        <a href="/chart/etherprice">$1,627.06</a>
                        <span data-bs-toggle="tooltip" data-bs-placement="top" title="Changes in the last 24 hours">
                            <span class="text-danger">(-0.59%)</span>
                        </span>
                    </div>
                    <div class="text-muted d-none d-lg-block">
                        <i class="fad fa-gas-pump me-1"></i>
                        Gas:
                        <span id="spanGasTooltip" data-bs-toggle="tooltip" data-bs-html="true" title="Base Fee: 23 Gwei<br>Priority Fee: 0 Gwei">
                            <a href="/gastracker">
                                <span class="gasPricePlaceHolder">23</span>
                                Gwei
                            </a>
                        </span>
                    </div>
                </div>
                <div class="d-flex justify-content-end align-items-center gap-2 w-100">
                    <div id="frmMaster" class="search-panel-container flex-grow-1 position-relative" style="width: 30rem;">
                        <form action="/search" method="GET" autocomplete="off" spellcheck="false" class="auto-search-max-height">
                            <span class="d-flex align-items-center position-absolute top-0 bottom-0" title="Search" style="left: 0.75rem;">
                                <i class="fa-regular fa-search fs-5 text-secondary"></i>
                            </span>
                            <input type="text" class="form-control form-control-lg bg-light bg-focus-white pe-20" autocomplete="off" spellcheck="false" id="search-panel" name="q" placeholder="Search by Address / Txn Hash / Block / Token / Domain Name" aria-describedby="button-header-search" onkeyup="handleSearchText(this);" maxlength="68" style="padding-left: 2.375rem;"/>
                            <a href="javascript:;" class="clear-icon d-none align-items-center position-absolute top-0 bottom-0 my-auto d-flex align-items-center" style="right: 3.375rem; cursor:pointer;">
                                <i class="fa-regular fa-xmark fs-5 text-secondary"></i>
                            </a>
                            <a href="javascript:;" class="search-icon d-none btn btn-sm btn-white my-1.5 align-items-center position-absolute top-0 bottom-0 d-flex align-items-center" style="right: 0.75rem; cursor:pointer;">
                                <i class="fa-regular fa-arrow-turn-down-left text-secondary"></i>
                            </a>
                            <input type="hidden" value id="hdnSearchText"/>
                            <input type="hidden" value id="hdnSearchLabel"/>
                            <input id="hdnIsTestNet" value="False" type="hidden"/>
                            <span class="shortcut-slash-icon align-items-center position-absolute top-0 bottom-0 align-items-center d-none d-sm-flex" title="Search" style="right: 0.5rem;">
                                <kbd class="bg-dark bg-opacity-25 fw-semibold px-2 py-0.5 text-white">/</kbd>
                            </span>
                        </form>
                    </div>
                    <div id="divThemeSetting" class="dropdown d-none d-lg-block">
                        <button class="btn btn-lg btn-white text-primary content-center" type="button" id="dropdownMenuTopbarSettings" data-bs-auto-close="outside" data-bs-toggle="dropdown" aria-expanded="false" style="width: 2.375rem; height: 2.375rem;">
                            <span class="theme-btn-main">
                                <i class="far fa-sun-bright theme-icon-active" data-href="#fa-sun-bright"></i>
                            </span>
                        </button>
                        <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="dropdownMenuTopbarSettings">
                            <li>
                                <button type="button" class="dropdown-item theme-btn active" data-bs-theme-value="light" onclick="setThemeMode('light');">
                                    <i class="far fa-sun-bright fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-sun-bright"></i>
                                    Light

                                </button>
                            </li>
                            <li>
                                <button type="button" class="dropdown-item theme-btn" data-bs-theme-value="dim" onclick="setThemeMode('dim');">
                                    <i class="far fa-moon-stars fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-moon-stars"></i>
                                    Dim

                                </button>
                            </li>
                            <li>
                                <button type="button" class="dropdown-item theme-btn" data-bs-theme-value="dark" onclick="setThemeMode('dark');">
                                    <i class="far fa-moon fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-moon"></i>
                                    Dark

                                </button>
                            </li>
                            <li>
                                <hr class="dropdown-divider">
                            </li>
                            <li>
                                <a class="dropdown-item" href="/settings">
                                    <i class="far fa-gear fa-fw dropdown-item-icon me-1"></i>
                                    Site Settings <i class="far fa-arrow-right ms-1"></i>
                                </a>
                            </li>
                        </ul>
                    </div>
                    <div id="divTestNet" class="dropdown d-none d-lg-block">
                        <button class="btn btn-lg btn-white content-center" type="button" id="dropdownTopbarNetworks" data-bs-toggle="dropdown" aria-expanded="false" style="width: 2.375rem; height: 2.375rem;">
                            <img width="10" data-img-theme="light" src="/images/svg/brands/ethereum-original.svg" alt="Ethereum Logo">
                            <img width="10" data-img-theme="darkmode" src="/images/svg/brands/ethereum-original-light.svg" alt="Ethereum Logo">
                        </button>
                        <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="dropdownTopbarNetworks">
                            <li>
                                <a href="https://etherscan.io" id="LI_Mainnet" class="dropdown-item active active">Ethereum Mainnet
</a>
                            </li>
                            <li>
                                <a href="https://cn.etherscan.com/?lang=zh-CN" id="LI_Mainnet_CN" class="dropdown-item">
                                    Ethereum Mainnet <span class="badge border bg-light text-dark ms-1">CN</span>
                                </a>
                            </li>
                            <li>
                                <a href="https://beaconscan.com/" id="LI2" class="dropdown-item">
                                    Beaconscan <span class="badge border bg-light text-dark ms-1">ETH2</span>
                                </a>
                            </li>
                            <li>
                                <hr class="dropdown-divider">
                            </li>
                            <li>
                                <a href="https://goerli.etherscan.io" id="LI58" class="dropdown-item">Goerli Testnet
</a>
                            </li>
                            <li>
                                <a href="https://sepolia.etherscan.io" id="LI9" class="dropdown-item">Sepolia Testnet
</a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </section>
        <header id="masterHeader" class="header border-bottom d-print-none">
            <nav class="navbar navbar-expand-lg navbar-light py-3 py-lg-0">
                <div class="container-xxl position-relative">
                    <a class="navbar-brand" href="/" target="_parent" aria-label="Etherscan">
                        <img width="150" data-img-theme="light" src="/assets/svg/logos/logo-etherscan.svg?v=0.0.5" alt="Etherscan Logo">
                        <img width="150" data-img-theme="darkmode" src="/assets/svg/logos/logo-etherscan-light.svg?v=0.0.5" alt="Etherscan Logo">
                    </a>
                    <div class="d-flex align-items-center gap-4">
                        <a class="link-dark d-block d-lg-none" href="/login">
                            <i class="far fa-user-circle me-1"></i>
                            Sign In

                        </a>
                        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                            <span class="navbar-toggler-icon"></span>
                        </button>
                    </div>
                    <div class="collapse navbar-collapse justify-content-end" id="navbarSupportedContent">
                        <ul class="navbar-nav gap-1 gap-lg-0 pt-4 pt-lg-0">
                            <li class="nav-item">
                                <a href="/" id="LI_default" class="nav-link" aria-current="page">Home
</a>
                            </li>
                            <li class="nav-item dropdown">
                                <a href="javascript:;" id="LI_blockchain" class="nav-link dropdown-toggle active" role="button" data-bs-toggle="dropdown" aria-expanded="false">Blockchain</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="/txs" id="LI12" class="dropdown-item">Transactions
</a>
                                    </li>
                                    <li>
                                        <a href="/txsPending" id="LI16" class="dropdown-item">Pending Transactions
</a>
                                    </li>
                                    <li>
                                        <a href="/txsInternal" id="LI14" class="dropdown-item">Contract Internal Transactions
</a>
                                    </li>
                                    <li>
                                        <a href="/txsBeaconDeposit" id="LI22" class="dropdown-item">Beacon Deposits
</a>
                                    </li>
                                    <li>
                                        <a href="/txsBeaconWithdrawal" id="LI_BeaconWithdrawals" class="dropdown-item">Beacon Withdrawals
</a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="/blocks" id="LI_blocks" class="dropdown-item">View Blocks
</a>
                                    </li>
                                    <li>
                                        <a href="/blocks_forked" id="LI_blocks2" class="dropdown-item">Forked Blocks (Reorgs)
</a>
                                    </li>
                                    <li>
                                        <a href="/uncles" id="LI8" class="dropdown-item">Uncles
</a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="/accounts" id="LI_accountall" class="dropdown-item">Top Accounts
</a>
                                    </li>
                                    <li>
                                        <a href="/contractsVerified" id="LI_contract_verified" class="dropdown-item active">Verified Contracts
</a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown">
                                <a href="javascript:;" id="LI_tokens" class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown" aria-expanded="false">Tokens</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="/tokens" id="LI21" class="dropdown-item">
                                            Top Tokens <span class="small text-muted">(ERC-20)</span>
                                        </a>
                                    </li>
                                    <li>
                                        <a href="/tokentxns" id="LI1" class="dropdown-item">
                                            Token Transfers <span class="small text-muted">(ERC-20)</span>
                                        </a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown">
                                <a href="javascript:;" id="LI_Nfts" class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown" aria-expanded="false">NFTs</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="/nft-top-contracts" id="LI63" class="dropdown-item">Top NFTs
</a>
                                    </li>
                                    <li>
                                        <a href="/nft-top-mints" id="LI67" class="dropdown-item">Top Mints
</a>
                                    </li>
                                    <li>
                                        <a href="/nft-trades" id="LI64" class="dropdown-item">Latest Trades
</a>
                                    </li>
                                    <li>
                                        <a href="/nft-transfers" id="LI65" class="dropdown-item">Latest Transfers
</a>
                                    </li>
                                    <li>
                                        <a href="/nft-latest-mints" id="LI66" class="dropdown-item">Latest Mints
</a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown">
                                <a href="javascript:;" id="LI_resources" class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown" aria-expanded="false">Resources</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="/charts" id="LI_charts2" class="dropdown-item">Charts And Stats
</a>
                                    </li>
                                    <li>
                                        <a href="/topstat" id="LI_topstat" class="dropdown-item">Top Statistics
</a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="/directory" id="LI62" class="dropdown-item">Directory
</a>
                                    </li>
                                    <li>
                                        <a href="https://info.etherscan.com/newsletters/" id="LI60" class="dropdown-item">Newsletter
</a>
                                    </li>
                                    <li>
                                        <a href="https://info.etherscan.com/" id="LI61" class="dropdown-item">Knowledge Base
</a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown">
                                <a href="#" id="li_developers" class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown" aria-expanded="false">Developers</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="/apis" id="LI5" class="dropdown-item">API Plans
</a>
                                    </li>
                                    <li>
                                        <a href="https://docs.etherscan.io/" id="LI6" class="dropdown-item" target="_blank">API Documentation
</a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="/code-reader" id="A1" class="dropdown-item">
                                            Code Reader <span class="badge border bg-light text-muted">Beta</span>
                                        </a>
                                    </li>
                                    <li>
                                        <a href="/verifyContract" id="LI17" class="dropdown-item">Verify Contract
</a>
                                    </li>
                                    <li>
                                        <a href="/find-similar-contracts" id="LI55" class="dropdown-item">Similar Contract Search
</a>
                                    </li>
                                    <li>
                                        <a href="/searchcontract" id="LI53" class="dropdown-item">Smart Contract Search
</a>
                                    </li>
                                    <li>
                                        <a href="/contractdiffchecker" id="LI54" class="dropdown-item">Contract Diff Checker
</a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="/vyper" id="LI27" class="dropdown-item">Vyper Online Compiler
</a>
                                    </li>
                                    <li>
                                        <a href="/opcode-tool" id="LI24" class="dropdown-item">Bytecode to Opcode
</a>
                                    </li>
                                    <li>
                                        <a href="/pushTx" id="LI10" class="dropdown-item">Broadcast Transaction
</a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown position-initial">
                                <a href="javascript:;" id="LI_services2" class="nav-link dropdown-toggle" role="button" data-bs-toggle="dropdown" aria-expanded="false">More</a>
                                <div class="dropdown-menu dropdown-menu-border dropdown-menu-mega">
                                    <div class="row">
                                        <div class="col-lg order-last order-lg-first">
                                            <div class="d-flex flex-column bg-light h-100 rounded-3 p-5">
                                                <div>
                                                    <h6>Tools &amp;Services</h6>
                                                    <p>Discover more of Etherscan's tools and services in one place.</p>
                                                </div>
                                                <div class="mt-auto">
                                                    <p class="text-muted mb-2">Sponsored</p>
                                                    <a target="_blank" href="https://chat.blockscan.com">
                                                        <img width="100" data-img-theme="light" src="/images/svg/blockscan-logo-dark.svg?v=0.0.4" alt/>
                                                        <img width="100" data-img-theme="dark" src="/images/svg/blockscan-logo-light.svg?v=0.0.4" alt/>
                                                    </a>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="col-sm py-5">
                                            <h6 class="px-3 mb-3">Tools</h6>
                                            <ul class="list-unstyled">
                                                <li>
                                                    <a href="/unitconverter" id="LI50" class="dropdown-item">
                                                        <i class="far fa-arrows-rotate dropdown-item-icon fa-fw me-1"></i>
                                                        Unit Converter

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/exportData" id="LI51" class="dropdown-item">
                                                        <i class="far fa-download fa-fw me-1"></i>
                                                        CSV Export

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/balancecheck-tool" id="LI52" class="dropdown-item">
                                                        <i class="far fa-file-invoice-dollar fa-fw me-1"></i>
                                                        Account Balance Checker

                                                    </a>
                                                </li>
                                            </ul>
                                        </div>
                                        <div class="col-sm py-5">
                                            <h6 class="px-3 mb-3">Explore</h6>
                                            <ul class="list-unstyled">
                                                <li>
                                                    <a href="/gastracker" id="LI19" class="dropdown-item">
                                                        <i class="far fa-gas-pump dropdown-item-icon fa-fw me-1"></i>
                                                        Gas Tracker

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/dex" id="LI4" class="dropdown-item">
                                                        <i class="far fa-arrow-right-arrow-left dropdown-item-icon fa-fw me-1"></i>
                                                        DEX Tracker

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/nodetracker" id="LI46" class="dropdown-item">
                                                        <i class="far fa-server dropdown-item-icon fa-fw me-1"></i>
                                                        Node Tracker

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/labelcloud" id="LI41" class="dropdown-item">
                                                        <i class="far fa-signs-post dropdown-item-icon fa-fw me-1"></i>
                                                        Label Cloud

                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/name-lookup" id="LI26" class="dropdown-item">
                                                        <i class="far fa-magnifying-glass-chart dropdown-item-icon fa-fw me-1"></i>
                                                        Domain Name Lookup

                                                    </a>
                                                </li>
                                            </ul>
                                        </div>
                                        <div class="col-sm py-5">
                                            <h6 class="px-3 mb-3">Services</h6>
                                            <ul class="list-unstyled">
                                                <li>
                                                    <a href="/tokenapprovalchecker" id="LI49" class="dropdown-item">
                                                        <i class="far fa-shield-keyhole dropdown-item-icon fa-fw me-1"></i>
                                                        Token Approvals <span class="badge border bg-light text-muted">Beta</span>
                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/verifiedSignatures" id="LI29" class="dropdown-item">
                                                        <i class="far fa-signature-lock dropdown-item-icon fa-fw me-1"></i>
                                                        Verified Signature

                                                    </a>
                                                </li>
                                                <li>
                                                    <a class="dropdown-item" href="/idm">
                                                        <i class="far fa-message-lines dropdown-item-icon fa-fw me-1"></i>
                                                        Input Data Messages (IDM) <span class="badge border bg-light text-muted">Beta</span>
                                                    </a>
                                                </li>
                                                <li>
                                                    <a href="/advanced-filter" id="LI31" class="dropdown-item">
                                                        <i class="far fa-filters dropdown-item-icon fa-fw me-1"></i>
                                                        Advanced Filter <span class="badge border bg-light text-muted">Beta</span>
                                                    </a>
                                                </li>
                                                <li>
                                                    <a class="dropdown-item" href="https://chat.blockscan.com" target="_blank">
                                                        <i class="far fa-messages dropdown-item-icon fa-fw me-1"></i>
                                                        Blockscan Chat <i class="far fa-arrow-up-right-from-square text-muted ms-1"></i>
                                                        <span class="badge border bg-light text-muted">Beta</span>
                                                    </a>
                                                </li>
                                            </ul>
                                        </div>
                                    </div>
                                </div>
                            </li>
                            <li class="nav-item dropdown d-block d-lg-none">
                                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">Explorers</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <a href="https://etherscan.io" id="LI_Mainnet_1" class="dropdown-item active active">Ethereum Mainnet
</a>
                                    </li>
                                    <li>
                                        <a href="https://cn.etherscan.com/?lang=zh-CN" id="LI_Mainnet_CN_1" class="dropdown-item">
                                            Ethereum Mainnet <span class="badge border bg-light text-dark ms-1">CN</span>
                                        </a>
                                    </li>
                                    <li>
                                        <a href="https://beaconscan.com/" id="LI77" class="dropdown-item">
                                            Beaconscan <span class="badge border bg-light text-dark ms-1">ETH2</span>
                                        </a>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a href="https://goerli.etherscan.io" id="LI78" class="dropdown-item">Goerli Testnet
</a>
                                    </li>
                                    <li>
                                        <a href="https://sepolia.etherscan.io" id="LI79" class="dropdown-item">Sepolia Testnet
</a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item dropdown d-block d-lg-none">
                                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">Appearance &amp;Settings</a>
                                <ul class="dropdown-menu dropdown-menu-border" style="min-width: 14rem;">
                                    <li>
                                        <button type="button" class="dropdown-item theme-btn active" data-bs-theme-value="light" onclick="setThemeMode('light');">
                                            <i class="far fa-sun-bright fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-sun-bright"></i>
                                            Light

                                        </button>
                                    </li>
                                    <li>
                                        <button type="button" class="dropdown-item theme-btn" data-bs-theme-value="dim" onclick="setThemeMode('dim');">
                                            <i class="far fa-moon-stars fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-moon-stars"></i>
                                            Dim

                                        </button>
                                    </li>
                                    <li>
                                        <button type="button" class="dropdown-item theme-btn" data-bs-theme-value="dark" onclick="setThemeMode('dark');">
                                            <i class="far fa-moon-stars fa-fw dropdown-item-icon theme-icon me-1" data-href="#fa-moon"></i>
                                            Dark

                                        </button>
                                    </li>
                                    <li>
                                        <hr class="dropdown-divider">
                                    </li>
                                    <li>
                                        <a class="dropdown-item" href="/settings">
                                            <i class="far fa-gear fa-fw dropdown-item-icon me-1"></i>
                                            Site Settings <i class="far fa-arrow-right ms-1"></i>
                                        </a>
                                    </li>
                                </ul>
                            </li>
                            <li class="nav-item align-self-center d-none d-lg-block">
                                <span class="text-secondary">|</span>
                            </li>
                            <li class="nav-item d-none d-lg-block">
                                <a class="nav-link" href="/login">
                                    <i class="far fa-user-circle me-1"></i>
                                    Sign In

                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
        </header>
        <main id="content" class="main-content" role="main">
            <input type="hidden" name="hdnAgeText" id="hdnAgeText" value="Age"/>
            <input type="hidden" name="hdnDateTimeText" id="hdnDateTimeText" value="Date Time (UTC)"/>
            <input type="hidden" name="hdnAgeTitle" id="hdnAgeTitle" value="Click to show Age Format"/>
            <input type="hidden" name="hdnDateTimeTitle" id="hdnDateTimeTitle" value="Click to show Datetime Format"/>
            <input type="hidden" name="hdnTxnText" id="hdnTxnText" value="Txn Fee"/>
            <input type="hidden" name="hdnGasPriceText" id="hdnGasPriceText" value="Gas Price"/>
            <input type="hidden" name="hdnTxnFeeTitle" id="hdnTxnFeeTitle" value="(Gas Price * Gas Used by Txns) in Ether"/>
            <input type="hidden" name="hdnGasPriceTitle" id="hdnGasPriceTitle" value="Gas Price in Gwei"/>
            <style>
                /*.tooltip-inner {
            max-width: 290px;
        }*/
            </style>
            <section class="container-xxl">
                <div class="d-flex flex-wrap justify-content-between align-items-center border-bottom gap-3 py-5">
                    <div class="d-flex flex-column gap-1">
                        <h1 class="h5 mb-0">Verified Contracts</h1>
                    </div>
                </div>
            </section>
            <div class="container-xxl">
                <div class="py-3">
                    <div class="d-flex text-muted" style="line-height: 2.2;">
                        <span id="ContentPlaceHolder1_lblAdResult2">
                            <ins data-revive-zoneid="2" data-revive-id="6452186c83cd256052c3c100f524ed97"></ins>
                            <script async src="//kta.etherscan.com/www/d/asyncjses.php?v=0.01"></script>
                        </span>
                        &nbsp;
                    </div>
                </div>
            </div>
            <span id="ContentPlaceHolder1_lblAdResult"></span>
            <section class="container-xxl pt-5 pb-12">
                <div id="ContentPlaceHolder1_mainrow" class="card">
                    <div class="card-header d-flex flex-column flex-sm-row justify-content-between gap-3">
                        <ul class="nav nav-pills text-nowrap align-items-center pb-3 gap-2" role="Filters">
                            <li class="nav-item snap-align-start">
                                <form>
                                    <button class="btn btn-sm btn-white dropdown-toggle" type="button" id="dropdownCategories" data-bs-toggle="dropdown" aria-expanded="false">
                                        Filter by |
                                        <div class="content-center" style="width: 16px; height: 16px;">
                                            <i class="far fa-grid-2 text-dark"></i>
                                        </div>
                                        Latest 500 Verified Contracts
                                    </button>
                                    <div id="selectTypeButton" class="dropdown-menu" aria-labelledby="dropdownCategories" style="min-width: 200px;">
                                        <div class="overflow-y-auto" style="max-height: 22rem;">
                                            <a class="dropdown-item d-flex align-items-center" href="/contractsVerified">
                                                <div class="content-center me-1.5" style="width: 16px; height: 16px;">
                                                    <i class="far fa-grid-2 text-dark"></i>
                                                </div>
                                                Latest 500 Verified Contracts
                                            </a>
                                            <li>
                                                <hr class="dropdown-divider">
                                            </li>
                                            <a class="dropdown-item d-flex align-items-center" href="/contractsVerified?filter=solc">
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/solidity.svg" data-img-theme="light" alt/>
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/solidity-light.svg" data-img-theme="dim" alt/>
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/solidity-light.svg" data-img-theme="dark" alt/>Solidity Compiler
                                            </a>
                                            <a class="dropdown-item d-flex align-items-center" href="/contractsVerified?filter=vyper">
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/vyper.svg" data-img-theme="light" alt/>
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/vyper-light.svg" data-img-theme="dim" alt/>
                                                <img class="rounded-circle me-1.5" width="16" src="./images/brands/vyper-light.svg" data-img-theme="dark" alt/>Vyper Compiler
                                            </a>
                                            <li>
                                                <hr class="dropdown-divider">
                                            </li>
                                            <a class="dropdown-item d-flex align-items-center" href="/contractsVerified?filter=opensourcelicense">
                                                <img class="filter-grayscale rounded-circle me-1.5" width="16" src="./images/brands/osi.png" alt/>Open Source License
                                            </a>
                                            <a class="dropdown-item d-flex align-items-center" href="/contractsVerified?filter=audit">
                                                <div class="content-center me-1.5" style="width: 16px; height: 16px;">
                                                    <i class="far fa-shield-check text-dark"></i>
                                                </div>
                                                Contract Security Audit
                                            </a>
                                        </div>
                                    </div>
                                </form>
                            </li>
                        </ul>
                        <div class="d-flex flex-wrap align-items-center gap-2">
                            <div class="dropdown order-md-2">
                                <button class="btn btn-sm btn-secondary js-dropdowns-input-focus" type="button" id="dropdownSearchFilter" data-bs-toggle="dropdown" aria-expanded="false">
                                    <i class="far fa-search"></i>
                                </button>
                                <div class="dropdown-menu" aria-labelledby="dropdownSearchFilter" style="min-width: 10rem">
                                    <div class="input-group" style="min-width: 20rem;">
                                        <form action="/searchcontractlist" method="get" class="js-focus-state input-group input-group-sm w-100" autocomplete="off">
                                            <input id="q" name="q" type="search" maxlength="60" value class="js-input-focus form-control py-2" placeholder="Search Contract Name"/>
                                            <input name="a" type="hidden" value="all"/>
                                            <button type="submit" class="btn btn-primary" data-bs-toggle="tooltip" title="Search Contract Name">Find</button>
                                        </form>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="card-body d-flex flex-wrap justify-content-between gap-3">
                        <div class="text-muted">
                            <span class="text-dark">&nbsp;Showing the last 500 verified contracts source code</span>
                        </div>
                        <nav aria-label="Table navigation">
                            <ul class="pagination pagination-sm mb-0">
                                <li class="page-item disabled">
                                    <span class="page-link">First</span>
                                </li>
                                <li class="page-item disabled">
                                    <span class="page-link px-3">
                                        <i class="fa fa-chevron-left small"></i>
                                        <span class="sr-only">Previous</span>
                                    </span>
                                </li>
                                <li Class="page-item disabled">
                                    <span Class="page-link">Page 1 of 20</span>
                                </li>
                                <li class="page-item" data-bs-toggle="tooltip" title="Go to Next">
                                    <a class="page-link px-3" href="/contractsVerified/2" aria-label="Next">
                                        <span aria-hidden="true">
                                            <i class="fa fa-chevron-right small"></i>
                                        </span>
                                        <span class="sr-only">Next</span>
                                    </a>
                                </li>
                                <li class="page-item" data-bs-toggle="tooltip" title="Go to Last">
                                    <a class="page-link" href="/contractsVerified/20" aria-label="Last">
                                        <span aria-hidden="true">Last</span>
                                        <span class="sr-only">Last</span>
                                    </a>
                                </li>
                            </ul>
                        </nav>
                    </div>
                    <div class="table-responsive">
                        <table class="table table-hover table-align-middle mb-0">
                            <thead class="text-nowrap">
                                <tr>
                                    <th scope="col">Address</th>
                                    <th scope="col">Contract Name</th>
                                    <th scope="col">Compiler</th>
                                    <th scope="col">Version</th>
                                    <th scope="col">Balance</th>
                                    <th scope="col" width="20">Txns</th>
                                    <th scope="col">Setting</th>
                                    <th scope="col">Verified</th>
                                    <th scope="col">
                                        Audit <i class="far fa-question-circle text-muted ms-1" data-bs-toggle="tooltip" data-bs-placement="top" data-bs-trigger="hover" title="Smart Contracts Audit and Security" data-bs-content="Smart Contracts Audit and Security"></i>
                                    </th>
                                    <th scope="col">
                                        License
                                        <a href="/contract-license-types" target="_blank">
                                            <i class="far fa-question-circle ms-1" data-bs-placement="top" data-bs-trigger="hover" data-bs-toggle="tooltip" title="Contract Source Code License Type, click for more info"></i>
                                        </a>
                                    </th>
                                    <th scope="col">Similar Contract</th>
                                </tr>
                            </thead>
                            <tbody class="align-middle text-nowrap">
                                <tr>
                                    <td>
                                        <span class="d-flex align-items-center gap-1">
                                            <a class="me-1" data-bs-trigger="hover" data-bs-toggle="tooltip" title="0xcf73b0d42c2c97219ce5895f311a7aa0fc930c98" href="/address/0xcf73b0d42c2c97219ce5895f311a7aa0fc930c98#code">
                                                <span class="d-flex align-items-center">
                                                    <i class="far fa-file-alt text-secondary me-1"></i>
                                                    0xcF73b0...fc930c98
                                                </span>
                                            </a>
                                            <a class="js-clipboard link-secondary " href="javascript:;" data-clipboard-text="0xcF73b0D42c2C97219Ce5895f311a7aa0fc930c98" data-bs-toggle="tooltip" data-bs-trigger="hover" title="Copy Address" data-hs-clipboard-options="{ &quot;type&quot;: &quot;tooltip&quot;, &quot;successText&quot;: &quot;Copied!&quot;, &quot;classChangeTarget&quot;: &quot;#linkIcon_f_tx_1&quot;, &quot;defaultClass&quot;: &quot;fa-copy&quot;, &quot;successClass&quot;: &quot;fa-check&quot; }">
                                                <i id="linkIcon_f_tx_1" class="far fa-copy fa-fw "></i>
                                            </a>
                                        </span>
                                    </td>
                                    <td>BEV</td>
                                    <td>Solidity</td>
                                    <td>
                                        <span>0.8.21</span>
                                    </td>
                                    <td>0 ETH</td>
                                    <td>1</td>
                                    <td>
                                        <div class="d-flex gap-1">
                                            <span class="badge bg-light border text-muted rounded-circle content-center" style="width: 1.5rem; height: 1.5rem;" data-bs-toggle="tooltip" data-bs-placement="top" title="Optimization Enabled">
                                                <i class="fas fa-bolt"></i>
                                            </span>
                                        </div>
                                    </td>
                                    <td>9/4/2023</td>
                                    <td>-</td>
                                    <td>None</td>
                                    <td>
                                        <a class="btn btn-sm btn-secondary" href="/find-similar-contracts?a=0xcf73b0d42c2c97219ce5895f311a7aa0fc930c98&m=low">
                                            <i class="far fa-search me-0.5"></i>
                                            Search
                                        </a>
                                    </td>
                                </tr>
`
