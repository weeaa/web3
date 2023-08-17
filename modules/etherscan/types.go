package etherscan

// remain unchanged if you do not want to get 403.
const (
	retryDelay = 3000
)

type Webhook struct {
	Address string
	Name    string
	Link    string
}
