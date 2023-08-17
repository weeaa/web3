package discord

import (
	"github.com/PuerkitoBio/goquery"
)

type Site string

type Client struct {
	AvatarImage        string
	FooterImage        string
	FooterText         string
	Color              int
	BRC20MintsWebhook  string
	LaunchMyNFTWebhook string
	ExchangeArtWebhook string
	PremintWebhook     string
}

type ExchangeArtWebhook struct {
	Name        string
	Description string
	Image       string
	MintLink    string
	CMID        string
	Supply      string
	ReleaseType string
	Minted      int
	MintCap     int
	Artist      string
	Edition     interface{}
	EditionBool interface{}
	Price       string
	LiveAt      string
	ToSend      bool
}

type BRC20MintsWebhook struct {
	Name             string
	Supply           string
	HoldersCount     string
	MintLimit        string
	TotalMinted      string
	MintLink         string
	PercentageMinted string
	MintTimes        string
	Links            string
	Creator          string
	Block            string
	Holders          map[int]map[string]string
	Fees             string
	Timestamp        string
	BlockDeploy      string
}

type PremintWebhook struct {
	document *goquery.Document

	Title        string
	URL          string
	Image        string
	Desc         string
	Price        string
	BalanceFall  string
	ETHtoHold    string
	TimeClose    string
	WinnerAmount string
	Status       string
	StatusImg    string

	Twitter TwitterReqs
	Discord DiscordReqs
	Misc    MiscReqs
	Custom  Custom
}

type TwitterReqs struct {
	Total   string
	Account string
	Tweet   string
}

type DiscordReqs struct {
	Total  string
	Server string
	Role   string
}

type MiscReqs struct {
	Total          string
	Spots          string
	OverAllocating string
	RegOut         string
	LinkOut        string
}

type Custom struct {
	Total string
}

/*🍀 DISCORD TYPES 🍀*/
type Webhook struct {
	Content   string  `json:"content,omitempty"`
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Tts       bool    `json:"tts,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	Url         string         `json:"url,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedFields  `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedFields struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}
