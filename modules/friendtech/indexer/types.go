package indexer

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/modules/twitter"
	"time"
)

const DefaultDelay = 3 * time.Second

type Indexer struct {
	DB           *db.DB
	ProxyList    []string
	Client       tls_client.HttpClient
	NitterClient *twitter.Client

	// The UserID from where you want to count. Note that it starts at 11.
	userCounter uint

	// The sleeping time between each request. Suggested 3 seconds.
	Delay time.Duration

	// Whether you want your proxies to rotate for each request.
	RotateEachReq bool

	// Whether you want to print out errors.
	Verbose bool
}
