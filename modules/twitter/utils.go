package twitter

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FetchNitter is a non rate limit alternative to fetch a Twitter user's information.
func FetchNitter(username string) (NitterResponse, error) {
	var nitter NitterResponse

	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{Scheme: "https", Host: "nitter.cz", Path: fmt.Sprintf("/%s", username)},
		Header: http.Header{
			"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/119.0"},
			"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
			"accept-language":           {"fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3"},
			"accept-encoding":           {"gzip, deflate, br"},
			"upgrade-insecure-requests": {"1"},
			"sec-fetch-dest":            {"document"},
			"connection":                {"keep-alive"},
			"sec-fetch-mode":            {"navigate"},
			"sec-fetch-site":            {"none"},
			"sec-fetch-user":            {"?1"},
		},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nitter, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nitter, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nitter, err
	}

	nitter.Followers = doc.Find("li[class=followers]").Find("span[class=profile-stat-num]").Text()
	nitter.JoinDate = doc.Find("div[class=profile-joindate]").Find("span").AttrOr("title", "")
	nitter.AccountAge = GetAccountAge(nitter.JoinDate)

	return nitter, nil
}

func GetAccountAge(date string) string {
	dateFormat := "Mon Jan 02 15:04:05 -0700 2006"

	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d days", int(math.Abs(t.Sub(time.Now()).Hours()/24)))
}
