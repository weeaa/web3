package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func NewClient() *Client {
	return &Client{}
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
			return fmt.Errorf("ERROR Push Discord Webhook invalid request %d [status exp. %s]", resp.StatusCode, "title")
		}
	}
}

func (c *Client) ExchangeArtNotification(content Webhook) error {

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = Push(jsonData, c.ExchangeArtWebhook)
	return err
}

func (c *Client) PremintNotification(content Webhook) error {

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = Push(jsonData, c.PremintWebhook)
	return err
}

func (c *Client) EtherscanNotification(content Webhook) error {

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = Push(jsonData, c.PremintWebhook)
	return err
}

func (c *Client) LaunchMyNFTNotification(content Webhook) error {

	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = Push(jsonData, c.PremintWebhook)
	return err
}

func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}

func (c *Client) CheckIfNil(webhook string) bool {
	return len(webhook) < 20
}
