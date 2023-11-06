package fren_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/rs/zerolog/log"
	"github.com/weeaa/nft/modules/friendtech/constants"
	"io"
	"math/big"
	"net/url"
)

var (
	ErrEmptyBearerToken = errors.New("error bearer token is empty")
)

type Client struct {
	Bearer  string
	Client  tls_client.HttpClient
	Proxies []string
}

func NewFriendTechUserClient(bearer string, client tls_client.HttpClient, proxies []string) *Client {
	return &Client{Bearer: fmt.Sprintf("Bearer %s", bearer), Client: client, Proxies: proxies}
}

// AddWishList adds every user you want to your wishlist.
func (c *Client) AddWishList(address string) error {
	if c.Bearer == "" {
		return ErrEmptyBearerToken
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: constants.ProdBaseApi, Path: "/watchlist-users/" + address},
		Host:   constants.ProdBaseApi,
		Header: http.Header{
			"authority":          {"prod-api.kosetto.com"},
			"accept":             {"application/json"},
			"authorization":      {c.Bearer},
			"accept-language":    {"en-US,en;q=0.9"},
			"accept-encoding":    {"gzip, deflate, br"},
			"referer":            {"https://www.friend.tech/"},
			"origin":             {"https://www.friend.tech"},
			"connection":         {"keep-alive"},
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-fetch-site":     {"cross-site"},
			"content-type":       {"application/json"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-dest":     {"empty"},
			"user-agent":         {constants.IphoneUserAgent},
		},
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error adding to wishlist: %s", resp.Status)
	}

	log.Info().Str("watchlist add", address)

	return nil
}

// RedeemCodes fetches all the invite codes of a user.
func (c *Client) RedeemCodes() ([]string, error) {
	if c.Bearer == "" {
		return nil, ErrEmptyBearerToken
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: constants.ProdBaseApi, Path: "/invite-codes"},
		Host:   constants.ProdBaseApi,
		Header: http.Header{
			"authority":          {"prod-api.kosetto.com"},
			"accept":             {"application/json"},
			"authorization":      {c.Bearer},
			"accept-language":    {"en-US,en;q=0.9"},
			"accept-encoding":    {"gzip, deflate, br"},
			"referer":            {"https://www.friend.tech/"},
			"origin":             {"https://www.friend.tech"},
			"connection":         {"keep-alive"},
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-fetch-site":     {"cross-site"},
			"content-type":       {"application/json"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-dest":     {"empty"},
			"user-agent":         {constants.IphoneUserAgent},
		},
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error redeeming codes: bad resp status %s – %s", resp.Status, string(body))
	}

	type Response struct {
		InviteCodes []struct {
			Code   string `json:"code"`
			IsUsed bool   `json:"isUsed"`
		} `json:"inviteCodes"`
	}

	var response Response
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	codes := make([]string, len(response.InviteCodes))
	for _, invite := range response.InviteCodes {
		codes = append(codes, invite.Code)
	}

	return codes, nil
}

// GetUserInformation returns the basic information of a user registered on FriendTech.
func GetUserInformation(address string, client tls_client.HttpClient) (UserInformation, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: constants.ProdBaseApi, Path: "/users/" + address},
		Host:   constants.ProdBaseApi,
		Header: http.Header{
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"user-agent":         {constants.IphoneUserAgent},
			"referer":            {"https://www.friend.tech/"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-ch-ua-mobile":   {"?0"},
			"dnt":                {"1"},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return UserInformation{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInformation{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return UserInformation{}, fmt.Errorf("error fetching user %s: bad resp status %s", address, resp.Status)
	}

	var r UserInformation
	if err = json.Unmarshal(body, &r); err != nil {
		return r, err
	}

	return r, nil
}

// AssertImportance assigns a status.
func AssertImportance(t, v any, impType ImpType) Importance {
	const none = "none"
	switch impType {
	case Followers:

		switch {
		case "friend tech":
		}

		n, ok := t.(int)
		if !ok {
			return none
		}

		// if superior to 5k is shrimp
		if n >= thresholds[0] && n <= thresholds[1] {
			return Shrimp
		}

		// if superior to
		if n >= thresholds[1] && n <= thresholds[2] {
			return Fish
		}

		if n >= thresholds[2] {
			return Whale
		}

		return none
	case Balance:
		n, ok := t.(*big.Int)
		if !ok {
			return none
		}

		thresholds := []*big.Int{}

		if n.Int64() >= thresholds[0].Int64() && n.Int64() <= thresholds[1].Int64() {
			return Shrimp
		}

		// if superior to
		if n.Int64() >= thresholds[1].Int64() && n.Int64() <= thresholds[2].Int64() {
			return Fish
		}

		if n.Int64() >= thresholds[2].Int64() {
			return Whale
		}

		return none
	default:
		return none
	}
}

func AssertChannelID() {}
