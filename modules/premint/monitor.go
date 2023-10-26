package premint

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bwmarrin/discordgo"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/discord/bot"
	"github.com/weeaa/nft/pkg/handler"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"golang.org/x/sync/errgroup"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"
)

func NewClient(bot *bot.Bot, raffleTypes []RaffleType, verbose bool, profile Profile) *Settings {
	return &Settings{
		RaffleTypes: raffleTypes,
		Handler:     handler.New(),
		Context:     context.Background(),
		Verbose:     verbose,
		Profile:     profile,
		Bot:         bot,
	}
}

func (s *Settings) StartMonitor() {
	logger.LogStartup(moduleName)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.monitorRaffles()
				return
			}
		}()
		for !s.monitorRaffles() {
			select {
			case <-s.Context.Done():
				logger.LogShutDown(moduleName)
				return
			default:
				time.Sleep(10 * time.Minute)
				continue
			}
		}
	}()
}

func NewRaffle() *Raffle {
	return &Raffle{Misc: MiscReqs{}, Discord: DiscordReqs{}, Twitter: TwitterReqs{}}
}

func (s *Settings) monitorRaffles() bool {
	wg := sync.WaitGroup{}
	for _, raffleUrl := range s.RaffleTypes {
		go func(rfType RaffleType) {
			wg.Add(1)
			defer wg.Done()

			raffles, err := s.parseRafflesURLs(string(rfType))
			if err != nil {
				logger.LogError(moduleName, err)
				return
			}

			filteredRaffles := s.filter(raffles)

			for _, raffle := range filteredRaffles {
				if _, ok := s.Handler.M.Get(raffle); !ok {
					go func(raffleURL string) {
						if err = s.do(raffleURL); err != nil {
							logger.LogError(moduleName, err)
						}
					}(raffle)
				}
			}
		}(raffleUrl)
	}
	wg.Wait()
	time.Sleep(s.RetryDelay)
	return false
}

func (s *Settings) Monitor(raffleTypes []RaffleType) {

	logger.LogStartup(moduleName)

	if err := s.Profile.login(); err != nil {
		logger.LogError(moduleName, err)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.Monitor(raffleTypes)
				return
			}
		}()
		for !s.monitorRaffles() {
			select {
			case <-s.Context.Done():
				return
			default:
				continue
			}
		}
	}()
}

// fetchRaffles parses available raffles on Premint.
func (s *Settings) parseRafflesURLs(raffleUrl string) ([]string, error) {
	raffleURL, err := url.Parse(raffleUrl)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    raffleURL,
		Header: http.Header{
			"Authority":                 {"www.premint.xyz"},
			"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			"Accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
			"Cache-Control":             {"max-age=0"},
			"Cookie":                    {s.Profile.getCookieHeader()},
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

	resp, err := s.Profile.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		tls.HandleRateLimit(s.Profile.Client, s.Profile.ProxyList, "")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	var raffles []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		raffles = append(raffles, s.AttrOr("href", ""))
	})

	return s.filter(raffles), nil
}

func (s *Settings) do(URL string) error {
	retries := 0
	raffle := NewRaffle()

	for i := 0; i < retries; i++ {
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: "www.premint.xyz", Path: URL},
			Header: http.Header{
				"Authority":                 {"www.premint.xyz"},
				"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"Accept-Language":           {"en-US,en;q=0.9,fr-FR;q=0.8,fr;q=0.7"},
				"Cache-Control":             {"max-age=0"},
				"Cookie":                    {s.Profile.getCookieHeader()},
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

		resp, err := s.Profile.Client.Do(req)
		if err != nil {
			retries++
			if retries >= 5 {
				return maxRetriesReached
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
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

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		raffle.document, err = goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			continue
		}

		if err = resp.Body.Close(); err != nil {
			continue
		}

		if raffle.isClosed() {
			return fmt.Errorf("raffle is closed")
		}

		raffle.doAllTasks()

		s.Handler.M.Set(raffle.Title, URL)
		raffle.Slug = URL
		s.Bot.BotWebhook(raffle.buildWebhook(), bot.PremintChannel)
		return nil
	}

	return fmt.Errorf("max error reached %s", URL)
}

// doAllTasks executes all the tasks used to fetch Raffle Information.
func (r *Raffle) doAllTasks() {
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		r.getProjectInfo()
		return nil
	})

	g.Go(func() error {
		r.getMiscInfo()
		return nil
	})

	g.Go(func() error {
		r.getDiscordInfo()
		return nil
	})

	g.Go(func() error {
		r.getCustomInfo()
		return nil
	})

	g.Go(func() error {
		r.getTwitterInfo()
		return nil
	})

	_ = g.Wait()

	r.updateIfNil()
}

func (r *Raffle) getProjectInfo() {
	title := r.document.Find("title").Text()
	r.Title = strings.Replace(title, "| PREMINT", "", -1)
	r.Desc = r.document.Find("meta[name=description]").AttrOr("content", "")
	r.document.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.AttrOr("src", ""), "https://premint.imgix.net") {
			r.Image, _ = s.Attr("src")
			r.Image = strings.Replace(r.Image, "&amp;", "", -1)
		}
	})

	if strings.Contains(r.document.Text(), "This project will be overallocating") {
		r.Misc.OverAllocating = "> • This project will be overallocating.\n"
		r.Misc.Total += r.Misc.OverAllocating
	}
}

func (r *Raffle) getTwitterInfo() {
	if strings.Contains(r.document.Text(), "step-twitter") {
		r.document.Find("a[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "@") {
				slice := []string{s.Text()}
				for _, v := range slice {
					account := v
					res := strings.ReplaceAll(v, "@", "")
					r.Twitter.Account = "> • Must Follow [" + account + "](https://twitter.com/" + res + ")\n"
				}
				r.Twitter.Total += r.Twitter.Account
			}
		})
	} else {
		r.Twitter.Account = ""
	}

	if strings.Contains(r.document.Text(), "Must Like &amp;") {
		r.document.Find("div[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Find("a[href]").AttrOr("href", "empty"), "twitter.com/user/status") {
				r.Twitter.Tweet = "> • Must Like & Retweet this [Tweet]" + "(" + s.Find("a[href]").AttrOr("href", "empty") + ")\n"
				r.Twitter.Total += r.Twitter.Tweet
			}
		})
	} else if strings.Contains(r.document.Text(), "Must Like") {
		r.document.Find("div[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Find("a[href]").AttrOr("href", "empty"), "twitter.com/user/status") {
				r.Twitter.Tweet = "> • Must Like this [Tweet]" + "(" + s.Find("a[href]").AttrOr("href", "empty") + ")\n"
				r.Twitter.Total += r.Twitter.Tweet
			}
		})
	}
}

func (r *Raffle) getDiscordInfo() {
	if strings.Contains(r.document.Text(), "Join the") && strings.Contains(r.document.Text(), "and have the") {
		var first string
		var firstRole string
		r.document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.AttrOr("href", "empty"), "discord.gg") || strings.Contains(s.AttrOr("href", "empty"), "discord.com") {
				//ServerName := s.Text()
				ServerURL, _ := s.Attr("href")
				sl := []string{ServerURL}
				first = sl[0]

			}
		})
		r.Discord.Server = "> • Must Join the [Discord Server](" + first + ")\n"
		r.Discord.Total += r.Discord.Server

		if strings.Contains(r.document.Find("div[class]").Add("span[class]").Text(), "and have the") {

			r.document.Find("div[class]").Add("span[class]").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.AttrOr("class", ""), "c-base-1 strong-700") {
					firstRole = s.Text()
				}
			})

			r.Discord.Role = "> • Must have the `" + firstRole + "` Role\n"
			r.Discord.Total += r.Discord.Role
		}

	} else if strings.Contains(r.document.Find("div[class]").Add("span[class]").Text(), "Join the") {
		var first string
		r.document.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.AttrOr("href", "empty"), "discord.gg") || strings.Contains(s.AttrOr("href", "empty"), "discord.com") {
				//ServerName := s.Text()
				ServerURL, ok := s.Attr("href")
				if !ok {

				}
				sl := []string{ServerURL}
				first = sl[0]
			}
		})
		r.Discord.Server = "> • Must Join the [Discord Server](" + first + ")\n"
		r.Discord.Total += r.Discord.Server
	}

}

func (r *Raffle) getCustomInfo() {
	if strings.Contains(r.document.Text(), "step-custom") {
		r.document.Find("label[class]").Each(func(i int, s *goquery.Selection) {
			r.Custom.Total = "> • " + s.Text()
		})
	}
}

func (r *Raffle) getMiscInfo() {
	if strings.Contains(r.document.Text(), "Have at least") {
		r.document.Find("i[class]").Each(func(i int, s *goquery.Selection) {
			if strings.Contains(s.Text(), "Have at least") { //todo: needs a fix
				//log.Println(s.Text())
			}
		})

		if strings.Contains(r.document.Text(), "disqualified if your balance") {
			r.BalanceFall = "> • Your entry will be disqualified if your balance falls below" + r.ETHtoHold
			r.Misc.Total += r.BalanceFall
		} else {

		}
	}
}

// updateIfNil verifies whether the content is empty. If it happens to be empty,
// it is replaced with a cross-mark symbol.
func (r *Raffle) updateIfNil() {
	if r.Twitter.Total == "" {
		r.Twitter.Total = "> • ❌"
	}
	if r.Discord.Total == "" {
		r.Discord.Total = "> • ❌"
	}
	if r.Misc.Total == "" {
		r.Misc.Total = "> • ❌"
	}
	if r.Custom.Total == "" {
		r.Custom.Total = "> • ❌"
	}
}

func (s *Settings) filter(slice []string) (filteredSlice []string) {
	for _, val := range slice {
		ok := s.verifyKeyword(val)
		if ok {
			filteredSlice = append(filteredSlice, val)
		}
	}
	return filteredSlice
}

// verifyKeyword verifies if the account is logged in by using a list of blacklisted URLs.
func (s *Settings) verifyKeyword(word string) bool {
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
		if word == keyword {
			s.Profile.isLoggedIn = false
			return false
		}
	}

	return true
}

func (r *Raffle) buildWebhook() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       r.Title,
				Description: r.Desc,
				URL:         "https://www.premint.xyz" + r.Slug,
				Timestamp:   discord.GetTimestamp(),
				Color:       bot.Purple,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf(""),
					IconURL: bot.Image,
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: r.Image,
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Twitter Reqs.",
						Value:  r.Twitter.Total,
						Inline: false,
					},
					{
						Name:   "Discord Reqs.",
						Value:  r.Discord.Total,
						Inline: false,
					},
					{
						Name:   "Custom Reqs.",
						Value:  r.Custom.Total,
						Inline: false,
					},
				},
			},
		},
	}
}

// isClosed verifies if the raffle does accept new entries.
func (r *Raffle) isClosed() bool {
	return strings.Contains(r.document.Text(), "This list is no longer accepting entries")
}
