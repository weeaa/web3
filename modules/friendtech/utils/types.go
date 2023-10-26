package fren_utils

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type Importance string
type ImpType string

var (
	Whale  Importance = "high"
	Fish   Importance = "medium"
	Shrimp Importance = "low"
)

var (
	WhaleEmote  = "🐋"
	FishEmote   = "🐠"
	ShrimpEmote = "🦐"
)

var (
	Balance   ImpType = "balance"
	Followers ImpType = "followers"
)

type Account struct {
	Email    string
	Password string
	Bearer   string
	client   tls_client.HttpClient
}

type UserInformation struct {
	Id                         int    `json:"id"`
	Address                    string `json:"address"`
	TwitterUsername            string `json:"twitterUsername"`
	TwitterName                string `json:"twitterName"`
	TwitterPfpUrl              string `json:"twitterPfpUrl"`
	TwitterUserId              string `json:"twitterUserId"`
	LastOnline                 string `json:"lastOnline"`
	LastMessageTime            any    `json:"lastMessageTime"`
	HolderCount                int    `json:"holderCount"`
	HoldingCount               int    `json:"holdingCount"`
	WatchlistCount             int    `json:"watchlistCount"`
	ShareSupply                int    `json:"shareSupply"`
	DisplayPrice               string `json:"displayPrice"`
	LifetimeFeesCollectedInWei string `json:"lifetimeFeesCollectedInWei"`
	UserBio                    any    `json:"userBio"`
}
