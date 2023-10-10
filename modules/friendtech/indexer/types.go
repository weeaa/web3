package indexer

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/database/db"
)

type Indexer struct {
	DB            *db.DB
	UserCounter   int
	ProxyList     []string
	Client        tls_client.HttpClient
	RotateEachReq bool // Whether or not you want your proxies to rotate for each request.
}
