package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	fren_utils "github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/modules/twitter"
	"github.com/weeaa/nft/pkg/tls"
	"io"
	"net/url"
	"os"
	"strconv"
)

const (
	baseHost = "localhost:992"
)

func AddUserToMonitor(baseAddress, by string) (map[string]any, error) {
	var buf bytes.Buffer

	nitterClient := twitter.NewClient("", "", []string{})
	client := tls.NewProxyLess()

	URL, _ := url.Parse("http://localhost:992/v1/user")

	userInfo, err := fren_utils.GetUserInformation(baseAddress, tls.NewProxyLess())
	if err != nil {
		return nil, err
	}

	nitter, err := nitterClient.FetchNitter(userInfo.TwitterUsername)
	if err != nil {
		return nil, err
	}

	followers, _ := strconv.Atoi(nitter.Followers)
	status := fren_utils.AssertImportance(followers, fren_utils.Followers)

	m := map[string]any{
		"base_address":     userInfo.Address,
		"status":           status,
		"twitter_username": userInfo.TwitterUsername,
		"twitter_name":     userInfo.TwitterName,
		"twitter_url":      "https://x.com/" + userInfo.TwitterUsername,
		"user_id":          userInfo.Id,
		"added_by":         by,
	}

	if err = json.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL:    URL,
		Body:   io.NopCloser(&buf),
		Header: http.Header{
			"authorization": {fmt.Sprintf("Basic %s", os.Getenv("BASIC_HASH"))},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error client: retry in some moments")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error req: invalid resp status")
	}

	// todo parse followers with nitter

	return map[string]any{
		"image":            userInfo.TwitterPfpUrl,
		"twitter_username": userInfo.TwitterUsername,
		"twitter_name":     userInfo.TwitterName,
		"followers":        nitter.Followers,
		"status":           status,
	}, nil
}
