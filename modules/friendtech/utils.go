package friendtech

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/weeaa/nft/db"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/handler"
	"github.com/weeaa/nft/pkg/files"
	"github.com/weeaa/nft/pkg/tls"
	"github.com/weeaa/nft/pkg/utils"
	"io"
	"net/url"
)

func NewClient(discordClient *discord.Client, verbose bool, WSSNodeUrl, HTTPNodeUrl string, db *db.DB) (*Settings, error) {
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

	return &Settings{
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

func NewSniper(privateKey, HTTPNodeUrl string) (*Sniper, error) {
	httpClient, err := ethclient.Dial(HTTPNodeUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to http node: %w", err)
	}

	return &Sniper{
		PrivateKey: privateKey,
		Wallet:     utils.InitWallet(privateKey),
		HttpClient: tls.NewProxyLess(),
		Client:     httpClient,
	}, nil
}

func NewIndexer(db *db.DB, proxyFilePath string) (*Indexer, error) {
	f, err := files.ReadJSON[map[string]int]("latestUserID.json")
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	return &Indexer{

		UserCounter: f["id"],
		DB:          db,
		ProxyList:   proxyList,
		Client:      tls.New(tls.RandProxyFromList(proxyList)),
	}, nil
}

/*
func (s *Settings) retrieveUserFromDB(sender string) (db.User, error) {

}
*/

// todo finish
func Login() *Account {
	u := &Account{}

	u.Bearer = "Bearer " + ""
	return u
}

// AddWishList adds every user you want to your wishlist.
func (u *Account) AddWishList(address string) error {

	req := &http.Request{
		Method: http.MethodPost,
		URL:    &url.URL{Scheme: "https", Host: prodBaseApi, Path: "/watchlist-users/" + address},
		Host:   prodBaseApi,
		Header: http.Header{
			"content-type":   {"application/json"},
			"accept":         {"application/json"},
			"authorization":  {u.Bearer},
			"sec-fetch-site": {"cross-site"},
			//	"accept-language": {""},
			"accept-encoding": {"gzip, deflate, br"},
		},
	}

	resp, err := u.client.Do(req)
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
func (u *Account) RedeemCodes() ([]string, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: prodBaseApi, Path: "/invite-codes"},
		Host:   prodBaseApi,
		Header: http.Header{
			"accept":          {"application/json"},
			"authorization":   {u.Bearer},
			"sec-fetch-site":  {"cross-site"},
			"accept-language": {"fr-FR,fr;q=0.9"},
			"accept-encoding": {"gzip, deflate, br"},
			"user-agent":      {iphoneUA},
			"referer":         {"https://www.friend.tech/"},
			"origin":          {"https://www.friend.tech"},
			"connection":      {"keep-alive"},
			"content-type":    {"applicatio/json"},
			"sec-fetch-mode":  {"cors"},
			"sec-fetch-dest":  {"empty"},
			"if-none-match":   {"W/\"24b-37JGT5Va+8NXn5xuE1XV2aJZJHY\""},
		},
	}

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error redeeming codes %s – %s", resp.Status, string(body))
	}

	type Response struct {
		InviteCodes []struct {
			Code   string `json:"code"`
			IsUsed bool   `json:"isUsed"`
		} `json:"inviteCodes"`
	}

	var r Response
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	codes := make([]string, len(r.InviteCodes))
	for _, invite := range r.InviteCodes {
		codes = append(codes, invite.Code)
	}

	return codes, nil
}

// GetUserInformation returns the basic information of a user registered on FriendTech.
func (s *Sniper) GetUserInformation(address string) (UserInformation, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: prodBaseApi, Path: "/users/" + address},
		Host:   prodBaseApi,
		Header: http.Header{
			"sec-ch-ua":          {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"user-agent":         {iphoneUA},
			"referer":            {"https://www.friend.tech/"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-ch-ua-mobile":   {"?0"},
			"dnt":                {"1"},
		},
	}

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return UserInformation{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInformation{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return UserInformation{}, fmt.Errorf("error fetching user information: bad resp status %s", resp.Status)
	}

	var r UserInformation
	if err = json.Unmarshal(body, &r); err != nil {
		return r, err
	}

	return r, nil
}

func buildWebhook(rolePing string) discord.Webhook {
	return discord.Webhook{
		Content: rolePing,
		Embeds: []discord.Embed{
			{
				Title: "",
			},
		},
	}
}

// can be
func assertImportance(t any, impType ImpType) (Importance, error) {
	switch impType {
	case Followers:
		n, ok := t.(int)
		if !ok {
			return Shrimp, fmt.Errorf("unable to assert role, defaulting to shrimp")
		}

		thresholds := []int{
			10000,
			50000,
		}

		if n <= thresholds[0] {
			return Shrimp, nil
		}

		if n >= thresholds[0] && n <= thresholds[1] {
			return Fish, nil
		}

		if n >= thresholds[1] {
			return Whale, nil
		}

		return Shrimp, fmt.Errorf("unable to assert importance: %d | %s", n, string(impType))
	case Balance:

		return Shrimp, nil
	default:
		return Shrimp, fmt.Errorf("")
	}
}
