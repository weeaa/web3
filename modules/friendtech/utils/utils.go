package fren_utils

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/modules/friendtech"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"log"
	"net/url"
)

func NewClient(discordClient *discord.Client, verbose bool, WSSNodeUrl, HTTPNodeUrl string, db *db.DB) (*friendtech.Settings, error) {
	ftAbi, err := utils.ReadABI("./abi/friendtech.json")
	if err != nil {
		return nil, fmt.Errorf("error reading abi: %w", err)
	}

	wssClient, err := ethclient.Dial(WSSNodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to wss node: %w", err)
	}

	httpClient, err := ethclient.Dial(HTTPNodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to http node: %w", err)
	}

	return &friendtech.Settings{
		WSSClient:  wssClient,
		HTTPClient: httpClient,
		Discord:    discordClient,
		Verbose:    verbose,
		Context:    context.Background(),
		Handler:    handler.New(),
		ABI:        ftAbi,
		DB:         db,
	}, nil
}

// AddWishList adds every user you want to your wishlist.
func AddWishList(address, bearer string, client tls_client.HttpClient) error {

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: ProdBaseApi, Path: "/watchlist-users/" + address},
		Host:   ProdBaseApi,
		Header: http.Header{
			"authority":          {"prod-api.kosetto.com"},
			"accept":             {"application/json"},
			"authorization":      {bearer},
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
			"user-agent":         {IphoneUserAgent},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error adding to wishlist: %s", resp.Status)
	}

	// à noter: pas de body resp, donc si 200 = OK

	return nil
}

// RedeemCodes fetches all the invite codes of a user.
func RedeemCodes(bearer string, client tls_client.HttpClient) ([]string, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: ProdBaseApi, Path: "/invite-codes"},
		Host:   ProdBaseApi,
		Header: http.Header{
			"authority":          {"prod-api.kosetto.com"},
			"accept":             {"application/json"},
			"authorization":      {bearer},
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
			"user-agent":         {IphoneUserAgent},
		},
	}

	resp, err := client.Do(req)
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
		URL:    &url.URL{Scheme: "https", Host: ProdBaseApi, Path: "/users/" + address},
		Host:   ProdBaseApi,
		Header: http.Header{
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"user-agent":         {IphoneUserAgent},
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

func AssertImportance(t any, impType ImpType) Importance {
	switch impType {
	case Followers:
		n, ok := t.(int)
		if !ok {
			log.Println("is not an int")
			return "none"
		}

		thresholds := []int{
			5000,
			10000,
			50000,
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

		return "none"
	case Balance:

		return Shrimp
	default:
		return Shrimp
	}
}
