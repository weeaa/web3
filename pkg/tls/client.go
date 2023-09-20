package tls

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bogdanfinn/tls-client"
	"math/rand"
	"os"
	"strings"
)

//var TestProxy = NewProxy(os.Getenv("TEST_PROXY"))

// New instantiates a TLS client and associates it with a user-defined proxy configuration.
func New(proxyURL string) tls_client.HttpClient {
	//ckJar := tls_client.NewCookieJar(nil)
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithTimeoutSeconds(tls_client.DefaultTimeoutSeconds),
		//	tls_client.WithCookieJar(ckJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithProxyUrl(NewProxy(proxyURL)),
	)
	if err != nil {
		return nil
	}
	return client
}

// NewProxyLess instantiates a TLS client configured to operate on the localhost IP address.
func NewProxyLess() tls_client.HttpClient {
	//ckJar := tls_client.NewCookieJar(nil)
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithTimeoutSeconds(tls_client.DefaultTimeoutSeconds),
		tls_client.WithTimeoutSeconds(10),
		//	tls_client.WithCookieJar(ckJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify())
	if err != nil {
		return nil
	}
	return client
}

func RotateProxy(client tls_client.HttpClient, proxyList []string) error {
	return client.SetProxy(NewProxy(RandProxyFromList(proxyList)))
}

// NewProxy parses a proxy in the correct format.
func NewProxy(unparsedProxy string) string {
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

// RandProxyFromList returns a random proxy stored in the list.
func RandProxyFromList(list []string) string {
	return list[rand.Intn(len(list))]
}

// ReadProxyFile reads a .txt file that contains a proxy on each new line & returns the proxies in a []string.
func ReadProxyFile(path string) (proxies []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		proxies = append(proxies, fileScanner.Text())
	}

	for i := range proxies {
		r := rand.Intn(i + 1)
		proxies[i], proxies[r] = proxies[r], proxies[i]
	}

	if len(proxies) == 0 {
		return nil, errors.New("empty proxy list")
	}

	return proxies, nil
}
