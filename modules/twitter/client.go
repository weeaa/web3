package twitter

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/pkg/tls"
)

type Client struct {
	OAuthToken string
	CSRFToken  string
	Client     tls_client.HttpClient
	Proxies    []string
}

// NewClient creates a client. If you are not using Nitter, be
// sure to provide values for CSRFToken and OAuthToken.
func NewClient(OAuthToken, CSRFToken string, proxies []string) *Client {
	return &Client{
		OAuthToken: OAuthToken,
		CSRFToken:  CSRFToken,
		Client:     tls.NewProxyLess(),
		Proxies:    proxies,
	}
}
