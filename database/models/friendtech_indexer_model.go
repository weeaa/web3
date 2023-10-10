package models

type (
	FriendTechIndexer struct {
		UserID          string `json:"user_id"`
		TwitterUsername string `json:"twitter_username"`
		BaseAddress     string `json:"base_address"`
	}
)
