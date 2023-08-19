package premint

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/logger"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net/url"
	"strings"
	"time"
)

const moduleName = "premint.xyz"

var maxRetriesReached = errors.New("maximum retries reached, aborting function")

func (p *Profile) Monitor(client *discord.Client, raffleTypes []RaffleType) {

	logger.LogStartup(moduleName)

	p.DiscordClient = client
	if err := p.login(); err != nil {
		logger.LogError(moduleName, err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			p.Monitor(client, raffleTypes)
			return
		}
	}()

	go func() {
		for {
			for _, raffleUrl := range raffleTypes {
				err := p.fetchRaffles(string(raffleUrl))
				if err != nil {
					logger.LogError(moduleName, err)
				}
			}
			time.Sleep(time.Duration(p.RetryDelay) * time.Millisecond)
		}
	}()
}

func (p *Profile) fetchRaffles(raffleUrl string) error {
	for {

		uri, err := url.Parse(raffleUrl)
		if err != nil {
			return err
		}

		req := &http.Request{
			Method: http.MethodGet,
			URL:    uri,
			Header: http.Header{
				"Authority":                 {"www.premint.xyz"},
				"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"Accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
				"Cache-Control":             {"max-age=0"},
				"Cookie":                    {p.getCookieHeader()},
				"Referer":                   {"https://www.premint.xyz/collectors/explore/"},
				"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\""},
				"Sec-Ch-Ua-Mobile":          {"?0"},
				"Sec-Fetch-Platform":        {"\"macOS\""},
				"Sec-Fetch-Dest":            {"document"},
				"Sec-Fetch-Mode":            {"navigate"},
				"Sec-Fetch-Site":            {"same-origin"},
				"Sec-Fetch-User":            {"?1"},
				"Upgrade-Insecure-Requests": {"1"},
				"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
			},
		}

		resp, err := p.Client.Do(req)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			continue
		}

		var raffles map[string]string
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			attr, _ := s.Attr("href")
			raffles[attr] = attr
		})

		filteredRafflesUrl := filter(raffles)
		for _, Url := range filteredRafflesUrl {
			if err = p.do(Url); err != nil {
				logger.LogError(moduleName, err)
			}
		}
		continue
	}
}

func (p *Profile) do(URL string) error {
	retries := 0
	task := &Webhook{}

	for {

		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: URL},
			Header: http.Header{
				"Authority":                 {"www.premint.xyz"},
				"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"Accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
				"Cache-Control":             {"max-age=0"},
				"Cookie":                    {p.getCookieHeader()},
				"Referer":                   {"https://www.premint.xyz/collectors/explore/"},
				"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\""},
				"Sec-Ch-Ua-Mobile":          {"?0"},
				"Sec-Fetch-Platform":        {"\"macOS\""},
				"Sec-Fetch-Dest":            {"document"},
				"Sec-Fetch-Mode":            {"navigate"},
				"Sec-Fetch-Site":            {"same-origin"},
				"Sec-Fetch-User":            {"?1"},
				"Upgrade-Insecure-Requests": {"1"},
				"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"},
			},
		}

		resp, err := p.Client.Do(req)
		if err != nil {
			retries++
			if retries >= 5 {
				return maxRetriesReached
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		if resp.StatusCode != 200 {
			if resp.StatusCode == 429 {
				return RateLimited
			}
			retries++
			if retries >= 5 {
				return maxRetriesReached
			}
			logger.LogError(moduleName, fmt.Errorf("[%s] Fetching Raffle %d/5 [%s]", resp.Status, retries, URL))
			continue
		}

		task.document, err = goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			continue
		}

		if err = task.checkClosed(); err != nil {
			return err
		}

		task.doAllTasks()

		p.handler.M.Set(task.Title, URL)

		if p.handler.M.Get(task.Title) == p.handler.MCopy.Get(task.Title) {
			return nil
		}

		time.Sleep(2 * time.Second)

		task.Slug = URL
		if err = p.DiscordClient.SendNotification(discord.Webhook{
			Username:  "Premint.xyz",
			AvatarUrl: "https://pbs.twimg.com/profile_images/1505785782002339840/mgeaHOqx_400x400.jpg",
			Embeds: []discord.Embed{
				{
					Title:       task.Title,
					Description: task.Desc,
					Url:         "https://www.premint.xyz" + task.Slug,
					Timestamp:   discord.GetTimestamp(),
					Color:       p.DiscordClient.Color,
					Footer: discord.EmbedFooter{
						Text:    p.DiscordClient.FooterText,
						IconUrl: p.DiscordClient.FooterImage,
					},
					Thumbnail: discord.EmbedThumbnail{
						Url: task.Image,
					},
					Fields: []discord.EmbedFields{
						{
							Name:   "Twitter Reqs.",
							Value:  task.Twitter.Total,
							Inline: false,
						},
						{
							Name:   "Discord Reqs.",
							Value:  task.Discord.Total,
							Inline: false,
						},
						{
							Name:   "Custom Reqs.",
							Value:  task.Custom.Total,
							Inline: false,
						},
					},
				},
			},
		}, discord.Premint); err != nil {
			logger.LogError(moduleName, err)
		}

		p.handler.M.ForEach(func(k string, v interface{}) {
			p.handler.MCopy.Set(k, v)
		})

		return nil
	}

}

func (t *Webhook) doAllTasks() {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		t.getProjectInfo()
		return nil
	})

	g.Go(func() error {
		t.getMiscInfo()
		return nil
	})

	g.Go(func() error {
		t.getDiscordInfo()
		return nil
	})

	g.Go(func() error {
		t.getCustomInfo()
		return nil
	})

	g.Go(func() error {
		t.getTwitterInfo()
		return nil
	})

	g.Wait()

	t.updateIfNil()
}

func (t *Webhook) getProjectInfo() {
	title := t.document.Find("title").Text()
	t.Title = strings.Replace(title, "| PREMINT", "", -1)
	t.Desc = t.document.Find("meta[name=description]").AttrOr("content", "")
	t.document.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.AttrOr("src", ""), "https://premint.imgix.net") {
			t.Image, _ = s.Attr("src")
			t.Image = strings.Replace(t.Image, "&amp;", "", -1)
		}
	})

	if strings.Contains(t.document.Text(), "This project will be overallocating") {
		t.Misc.OverAllocating = "> • This project will be overallocating.\n"
		t.Misc.Total += t.Misc.OverAllocating
	}
}

func (t *Webhook) getTwitterInfo() {
	if strings.Contains(t.document.Text(), "step-twitter") {
		t.document.Find("a[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "@") {
				slice := []string{s.Text()}
				for _, v := range slice {
					account := v
					res := strings.ReplaceAll(v, "@", "")
					t.Twitter.Account = "> • Must Follow [" + account + "](https://twitter.com/" + res + ")\n"
				}
				t.Twitter.Total += t.Twitter.Account
			}
		})
	} else {
		t.Twitter.Account = ""
	}

	if strings.Contains(t.document.Text(), "Must Like &amp;") {
		t.document.Find("div[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Find("a[href]").AttrOr("href", "empty"), "twitter.com/user/status") {
				t.Twitter.Tweet = "> • Must Like & Retweet this [Tweet]" + "(" + s.Find("a[href]").AttrOr("href", "empty") + ")\n"
				t.Twitter.Total += t.Twitter.Tweet
			}
		})
	} else if strings.Contains(t.document.Text(), "Must Like") {
		t.document.Find("div[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Find("a[href]").AttrOr("href", "empty"), "twitter.com/user/status") {
				t.Twitter.Tweet = "> • Must Like this [Tweet]" + "(" + s.Find("a[href]").AttrOr("href", "empty") + ")\n"
				t.Twitter.Total += t.Twitter.Tweet
			}
		})
	}
}

func (t *Webhook) getDiscordInfo() {
	if strings.Contains(t.document.Text(), "Join the") && strings.Contains(t.document.Text(), "and have the") {
		var first string
		var firstRole string
		t.document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.AttrOr("href", "empty"), "discord.gg") || strings.Contains(s.AttrOr("href", "empty"), "discord.com") {
				//ServerName := s.Text()
				ServerURL, _ := s.Attr("href")
				sl := []string{ServerURL}
				first = sl[0]

			}
		})
		t.Discord.Server = "> • Must Join the [Discord Server](" + first + ")\n"
		t.Discord.Total += t.Discord.Server

		if strings.Contains(t.document.Find("div[class]").Add("span[class]").Text(), "and have the") {

			t.document.Find("div[class]").Add("span[class]").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.AttrOr("class", ""), "c-base-1 strong-700") {
					firstRole = s.Text()
				}
			})

			t.Discord.Role = "> • Must have the `" + firstRole + "` Role\n"
			t.Discord.Total += t.Discord.Role
		}

	} else if strings.Contains(t.document.Find("div[class]").Add("span[class]").Text(), "Join the") {
		var first string
		t.document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.AttrOr("href", "empty"), "discord.gg") || strings.Contains(s.AttrOr("href", "empty"), "discord.com") {
				//ServerName := s.Text()
				ServerURL, ok := s.Attr("href")
				if !ok {

				}
				sl := []string{ServerURL}
				first = sl[0]
			}
		})
		t.Discord.Server = "> • Must Join the [Discord Server](" + first + ")\n"
		t.Discord.Total += t.Discord.Server
	}

}

func (t *Webhook) getCustomInfo() {
	if strings.Contains(t.document.Text(), "step-custom") {
		t.document.Find("label[class]").Each(func(i int, s *goquery.Selection) {
			t.Custom.Total = "> • " + s.Text()
		})
	}
}

func (t *Webhook) getMiscInfo() {
	if strings.Contains(t.document.Text(), "Have at least") {
		t.document.Find("i[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Have at least") { //todo: needs a fix
				//log.Println(s.Text())
			}
		})

		if strings.Contains(t.document.Text(), "disqualified if your balance") {
			t.BalanceFall = "> • Your entry will be disqualified if your balance falls below" + t.ETHtoHold
			t.Misc.Total += t.BalanceFall
		} else {

		}
	}
}

func (t *Webhook) updateIfNil() {
	if t.Twitter.Total == "" {
		t.Twitter.Total = "> • ❌"
	}
	if t.Discord.Total == "" {
		t.Discord.Total = "> • ❌"
	}
	if t.Misc.Total == "" {
		t.Misc.Total = "> • ❌"
	}
	if t.Custom.Total == "" {
		t.Custom.Total = "> • ❌"
	}
}

func filter(slice map[string]string) (filteredSlice []string) {
	for _, value := range slice {
		ok := checkKeywords(value)
		if ok {
			filteredSlice = append(filteredSlice, value)
		}
	}
	return filteredSlice
}

func checkKeywords(word string) bool {
	bannedKeywords := []string{
		"#",
		"https://www.premint.xyz/creators/",
		"https://www.premint.xyz/accounts/discord/login/?process=login&next=%2Fcollectors%2Fexplore%2Ftop%2F",
		"https://www.premint.xyz/accounts/twitter/login/?process=login&next=%2Fcollectors%2Fexplore%2Ftop%2F",
		"https://www.premint.xyz/collectors/",
		"https://www.premint.xyz/logout/",
		"https://www.premint.xyz/collectors/unfollow/",
		"https://www.premint.xyz/dashboard/",
		"https://www.premint.xyz/profile/",
		"https://www.premint.xyz/collectors/explore/new/",
		"https://www.premint.xyz/collectors/offers/",
		"https://www.premint.xyz/collectors/calendar/",
		"https://www.premint.xyz/explore/",
	}

	for _, keyword := range bannedKeywords {
		if strings.Contains(word, "") || word == keyword {
			return false
		}
	}

	return true
}

func (t *Webhook) checkClosed() error {
	if strings.Contains(t.document.Text(), "This list is no longer accepting entries") {
		log.Println("ERR Raffle Closed")
		return fmt.Errorf("raffle closed")
	}
	return nil
}
