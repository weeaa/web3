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
	"github.com/weeaa/nft/pkg/tls"
	"io"
	"net/url"
	"strings"
)

var RateLimited = errors.New("you are rate limited :( you got to wait till you're unbanned, which is approx 5+ minutes")

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
	}
}

func (p *Profile) login() error {
	for {

		//first req to get cookies
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https://", Host: "www.premint.xyz", Path: "/login"},
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

		if resp.StatusCode != 200 {
			if resp.StatusCode == 429 {
				return RateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		cookies := resp.Cookies()
		for _, c := range cookies {
			if c.Name == "csrftoken" {
				p.csrfToken = c.Value
			}
			if c.Name == "session_id" {
				p.sessionId = c.Value
			}
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		//second req is following the flow of the login 🌞
		params := url.Values{
			"username": {p.publicAddress},
		}

		req = &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Scheme: "https://", Host: "www.premint.xyz", Path: "/v1/signup_api/"},
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

		if resp.StatusCode != 200 {
			if resp.StatusCode == 429 {
				return RateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		privateKey, err := crypto.HexToECDSA(p.privateKey)
		if err != nil {
			return err
		}

		if err = p.getNonce(); err != nil {
			return err
		}

		signature, err := sign("Welcome to PREMINT!\n\nSigning is the only way we can truly know \nthat you are the owner of the wallet you \nare connecting. Signing is a safe, gas-less \ntransaction that does not in any way give \nPREMINT permission to perform any \ntransactions with your wallet.\n\nWallet address:\n"+p.publicAddress+"\n\nNonce: "+p.nonce, privateKey)
		if err != nil {
			return err
		}

		//third req completes login & update sessionId cookie
		req = &http.Request{
			Method: http.MethodPost,
			URL:    &url.URL{Scheme: "https://", Host: "www.premint.xyz", Path: "/v1/login_api/"},
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

		if resp.StatusCode != 200 {
			if resp.StatusCode == 429 {
				return RateLimited
			}
			return fmt.Errorf("invalid response status: %s", resp.Status)
		}

		cookieSess := resp.Cookies()
		for _, c := range cookieSess {
			if c.Name == "session_id" {
				p.sessionId = c.Value
			}
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}
	}
}

func (p *Profile) getNonce() error {

	type NonceResponse struct {
		Nonce   string `json:"data"`
		Success bool   `json:"success"`
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https://", Host: "www.premint.xyz", Path: "/v1/login_api/"},
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

	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			return RateLimited
		}
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var nr NonceResponse
	if err = json.Unmarshal(body, &nr); err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	if nr.Success {
		p.nonce = nr.Nonce
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
		return ":(", err
	}
	signatureBytes[64] += 27
	return hexutil.Encode(signatureBytes), nil
}

func (p *Profile) getCookieHeader() string {
	return fmt.Sprintf("csrftoken=%s;session_id=%s", p.csrfToken, p.sessionId)
}
