package premint

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils/ethereum"
	"io"
	"net/url"
	"strings"
)

func NewProfile(publicAddress, privateKey, proxy string, retryDelay int) *Profile {
	var client tls_client.HttpClient

	if proxy == "" {
		client = tls.NewProxyLess()
	} else {
		client = tls.New(proxy)
	}

	return &Profile{
		publicAddress: publicAddress,
		privateKey:    privateKey,
		Client:        client,
		RetryDelay:    retryDelay,
		isLoggedIn:    false,
		Wallet:        ethereum.InitWallet(privateKey),
	}
}

func (p *Profile) login() error {
	retries := 0
	for i := 0; i < retries; i++ {

		// 1. Get cookies (session).
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: "/login"},
			Host:   "www.premint.xyz",
			Header: http.Header{
				"user-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
				"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"sec-fetch-site":            {"same-origin"},
				"sec-fetch-mode":            {"navigate"},
				"sec-fetch-user":            {"?1"},
				"sec-fetch-dest":            {"document"},
				"accept-language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
				"upgrade-insecure-requests": {"1"},
			},
		}

		resp, err := p.Client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusTooManyRequests {
				return ErrRateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "csrftoken" {
				p.csrfToken = cookie.Value
			}
			if cookie.Name == "session_id" {
				p.sessionID = cookie.Value
			}
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		// 2. Initiate login.
		params := url.Values{
			"username": {p.publicAddress},
		}

		req = &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: "/v1/signup_api/"},
			Body:   io.NopCloser(strings.NewReader(params.Encode())),
			Host:   "www.premint.xyz",
			Header: http.Header{
				"cookie":                    {p.getCookieHeader()},
				"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
				"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"sec-fetch-site":            {"same-origin"},
				"sec-fetch-mode":            {"cors"},
				"sec-fetch-dest":            {"empty"},
				"accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
				"upgrade-Insecure-Requests": {"1"},
				"x-csrftoken":               {p.csrfToken},
				"content-type":              {"application/x-www-form-urlencoded"},
				"referer":                   {"https://www.premint.xyz/login/?next=/coliseum/"},
				"accept-encoding":           {"gzip, deflate, br"},
				"connection":                {"keep-alive"},
				"cache-control":             {"max-age=0"},
			},
		}

		resp, err = p.Client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusTooManyRequests {
				return ErrRateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		if err = p.getNonce(); err != nil {
			return err
		}

		signature, err := sign("Welcome to PREMINT!\n\nSigning is the only way we can truly know \nthat you are the owner of the wallet you \nare connecting. Signing is a safe, gas-less \ntransaction that does not in any way give \nPREMINT permission to perform any \ntransactions with your wallet.\n\nWallet address:\n"+p.publicAddress+"\n\nNonce: "+p.nonce, p.Wallet.PrivateKey)
		if err != nil {
			return err
		}

		// 3. Finish login & refresh sessionID.
		req = &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: "/v1/login_api/"},
			Body:   io.NopCloser(strings.NewReader("web3provider=metamask&address=" + p.publicAddress + "&signature=" + signature)),
			Host:   "www.premint.xyz",
			Header: http.Header{
				"cookie":             {p.getCookieHeader()},
				"sec-ch-ua":          {"\"Chromium\";v=\"115\", \"Not/A)Brand\";v=\"99\""},
				"x-csrftoken":        {p.csrfToken},
				"sec-ch-ua-mobile":   {"?0"},
				"content-type":       {"application/x-www-form-urlencoded; charset=UTF-8"},
				"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
				"sec-ch-ua-platform": {"\"macOS\""},
				"accept":             {"*/*"},
				"origin":             {"https://www.premint.xyz"},
				"sec-fetch-site":     {"same-origin"},
				"sec-fetch-mode":     {"cors"},
				"sec-fetch-dest":     {"empty"},
				"referer":            {"https://www.premint.xyz/login_api/"},
				"accept-language":    {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
			},
		}

		resp, err = p.Client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusForbidden {
				return ErrRateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "session_id" {
				p.sessionID = cookie.Value
			}
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		p.isLoggedIn = true
		log.Info().Bool("premint account login", p.isLoggedIn)

		return nil
	}

	p.isLoggedIn = false
	return ErrMaxRetriesReached
}

func (p *Profile) getNonce() error {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: "/v1/login_api/"},
		Host:   "www.premint.xyz",
		Header: http.Header{
			"Cookie":             {p.getCookieHeader()},
			"sec-ch-ua":          {"\"Chromium\";v=\"115\", \"Not/A)Brand\";v=\"99\""},
			"x-csrftoken":        {p.csrfToken},
			"sec-ch-ua-mobile":   {"?0"},
			"content-type":       {"application/x-www-form-urlencoded; charset=UTF-8"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"accept":             {"*/*"},
			"origin":             {"https://www.premint.xyz"},
			"sec-fetch-site":     {"same-origin"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-dest":     {"empty"},
			"referer":            {"https://www.premint.xyz/login/"},
			"accept-language":    {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
		},
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			return ErrRateLimited
		}
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var nonceResponse map[string]any
	if err = json.Unmarshal(body, &nonceResponse); err != nil {
		return err
	}

	if nonceResponse["success"].(bool) {
		p.nonce = nonceResponse["data"].(string)
		return nil
	} else {
		return errors.New("unable to get nonce")
	}
}

func sign(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	signatureBytes[64] += 27
	return hexutil.Encode(signatureBytes), nil
}

func (p *Profile) getCookieHeader() string {
	return fmt.Sprintf("csrftoken=%s;session_id=%s", p.csrfToken, p.sessionID)
}
