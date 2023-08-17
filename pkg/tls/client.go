package tls

import (
	"fmt"
	tls_client "github.com/bogdanfinn/tls-client"
	"strings"
)

func New(proxyURL string) tls_client.HttpClient {
	ckJar := tls_client.NewCookieJar(nil)
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithTimeoutSeconds(tls_client.DefaultTimeoutSeconds),
		tls_client.WithCookieJar(ckJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithProxyUrl(newProxy(proxyURL)),
	)
	if err != nil {
		return nil
	}
	return client
}

func NewProxyLess() tls_client.HttpClient {
	ckJar := tls_client.NewCookieJar(nil)
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithTimeoutSeconds(tls_client.DefaultTimeoutSeconds),
		tls_client.WithTimeoutSeconds(10),
		tls_client.WithCookieJar(ckJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify())
	if err != nil {
		return nil
	}
	return client
}

func newProxy(unparsedProxy string) string {
	var proxy string
	var rawProxy []string
	rawProxy = strings.Split(unparsedProxy, ":")
	if len(rawProxy) > 2 {
		proxy = fmt.Sprintf("http://%s:%s@%s:%s", rawProxy[2], rawProxy[3], rawProxy[0], rawProxy[1])
		return proxy
	} else {
		proxy = fmt.Sprintf("http://%s:%s", rawProxy[0], rawProxy[1])
		return proxy
	}
}
