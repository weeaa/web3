package twitter

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"math"
	"net/url"
	"strings"
	"time"
)

// FetchNitter offers a reduced rate limit alternative (& free) to the Twitter API.
func (c *Client) FetchNitter(username string) (*NitterResponse, error) {

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "nitter.cz", Path: "/" + username},
		Header: http.Header{
			"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"accept-language":           {"en-US,en;q=0.9"},
			"cache-control":             {"max-age=0"},
			"connection":                {"keep-alive"},
			"dnt":                       {"1"},
			"sec-fetch-dest":            {"document"},
			"sec-fetch-mode":            {"navigate"},
			"sec-fetch-site":            {"none"},
			"sec-fetch-user":            {"?1"},
			"upgrade-insecure-requests": {"1"},
			"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
			"sec-ch-ua":                 {"\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\""},
			"sec-ch-ua-mobile":          {"?0"},
			"sec-ch-ua-platform":        {"\"macOS\""},
		},
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nitter client error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting twitter information: bad response status [%s]", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return &NitterResponse{
		Followers:  strings.ReplaceAll(document.Find("li[class=followers]").Find("span[class=profile-stat-num]").Text(), ",", ""),
		JoinDate:   document.Find("div[class=profile-joindate]").Find("span").AttrOr("title", ""),
		AccountAge: getAccountAgeNitter(document.Find("div[class=profile-joindate]").Find("span").AttrOr("title", "")),
	}, nil
}

func GetAccountAge(date string) string {
	dateFormat := "Mon Jan 02 15:04:05 -0700 2006"

	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d days", int(math.Abs(t.Sub(time.Now()).Hours()/24)))
}

func getAccountAgeNitter(date string) string {
	timeParsed, err := time.Parse("3:04 PM - 2 Jan 2006", date)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d days", int(math.Abs(time.Since(timeParsed).Hours()/24)))
}
