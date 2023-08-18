package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func NewClient(exchangeArtWebhook, launchMyNFTWebhook, premintWebhook, etherscanWebhook, brc20Webhook, footerText, footerImage string, color int) *Client {
	return &Client{
		BRC20MintsWebhook:  brc20Webhook,
		ExchangeArtWebhook: exchangeArtWebhook,
		LaunchMyNFTWebhook: launchMyNFTWebhook,
		PremintWebhook:     premintWebhook,
		EtherscanWebhook:   etherscanWebhook,
		FooterImage:        footerImage,
		FooterText:         footerText,
		Color:              color,
	}
}

// Push sends a Discord Embed.
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
			return fmt.Errorf("push discord webhook invalid request: %s", resp.Status)
		}
	}
}

func (c *Client) SendNotification(content Webhook, module Module) error {
	var webhook string

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	switch module {
	case Premint:
		webhook = c.PremintWebhook
	case ExchangeArt:
		webhook = c.ExchangeArtWebhook
	case OpenSea:
		//todo add support
	case LaunchMyNFT:
		webhook = c.LaunchMyNFTWebhook
	case Etherscan:
		webhook = c.EtherscanWebhook
	}

	return Push(jsonData, webhook)
}

func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}

func (c *Client) CheckIfNil(webhook string) bool {
	return len(webhook) < 20
}
