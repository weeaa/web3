package etherscan

const (
	moduleName = "Etherscan Verified Contract"
	retryDelay = 3000
)

type Contract struct {
	Address string
	Name    string
	Link    string
}
