package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Push sends a Discord Embed
func Push(data []byte, webhookURL string) error {
	var resp *http.Response
	var err error

	for {

		resp, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		err = resp.Body.Close()
		if err != nil {
			return err
		}

		switch resp.StatusCode {
		case 204:
			return nil

		case 429:
			var timeout int

			timeout, err = strconv.Atoi(resp.Header.Get("retry-after"))
			if err == nil {
				time.Sleep(time.Duration(timeout) * time.Millisecond)
			} else {
				time.Sleep(5 * time.Second)
			}

		default:
			return fmt.Errorf("ERROR Push Discord Webhook invalid request %d [status exp. %s]", resp.StatusCode, "title")
		}
	}
}

func (w *ExchangeArtWebhook) ExchangeArtNotification(webhookURL string) error {
	content := Webhook{
		Username:  "ExchangeArt",
		AvatarUrl: weeaaImage,
		Embeds: []Embed{
			{
				Title:       w.Name,
				Description: w.Description,
				Thumbnail: EmbedThumbnail{
					Url: w.Image,
				},
				Color:     0xffffff,
				Timestamp: getTimestamp(),
				Footer: EmbedFooter{
					Text:    weeaaFooterText,
					IconUrl: weeaaImage,
				},
				Fields: []EmbedFields{{
					Name:   "Supply/Max(wallet)",
					Value:  fmt.Sprintf("`%s/%d`", w.Supply, w.MintCap),
					Inline: true,
				},
					{
						Name:   "Release Type",
						Value:  w.ReleaseType,
						Inline: true,
					},
					{
						Name:   "Artist",
						Value:  w.Artist,
						Inline: true,
					},
					{
						Name:   "Price",
						Value:  fmt.Sprintf("%s", w.Price[0:4]),
						Inline: true,
					},
					{
						Name:   "CandyMachine ID",
						Value:  fmt.Sprintf("`%s`", w.CMID),
						Inline: false,
					}},
			},
		},
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = Push(jsonData, webhookURL)
	return err
}

func getTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}
