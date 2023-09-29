package twitter

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/pkg/tls"
)

type Client struct {
	OAuthToken string
	CSRFToken  string
	Client     tls_client.HttpClient
}

func NewClient(OAuthToken, CSRFToken string) *Client {
	return &Client{
		CSRFToken: CSRFToken,
		Client:    tls.NewProxyLess(),
	}
}
