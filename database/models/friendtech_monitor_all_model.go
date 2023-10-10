package models

type (
	FriendTechMonitorAll struct {
		BaseAddress     string `json:"base_address"`     //base_address text
		Status          string `json:"status"`           //status text
		Followers       string `json:"followers"`        //followers text
		TwitterUsername string `json:"twitter_username"` //twitter_username text
		TwitterName     string `json:"twitter_name"`     //twitter_name text
		TwitterURL      string `json:"twitter_url"`      //twitter_url text
		UserID          int    `json:"user_id"`          //user_id integer
	}
)
