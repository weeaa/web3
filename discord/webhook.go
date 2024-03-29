package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func NewClient(
	footerText,
	footerImage string,
	color int) *Client {
	return &Client{
		FooterImage: footerImage,
		FooterText:  footerText,
		Color:       color,
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

func (c *Client) SendNotification(content Webhook, webhook string) error {
	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return Push(jsonData, webhook)
}

func GetTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05-0700")
}

func IsWebhookLenValid(webhook string) bool {
	return len(webhook) < 20
}

func (c *Client) WebhookNotificationTest(webhook string) error {
	return c.SendNotification(Webhook{Username: "Test Webhook", Embeds: []Embed{{Title: "Test Successful 🦄", Color: 0x008000}}}, webhook)
}
