package twitter

type AccountInformation struct {
	Address    string
	Name       string
	eth        string
	Followers  int
	SharePrice string
	Age        string
	Holders    string
	Twitter    Twitter
}

type Twitter struct {
	ProfilePicture string
	Name           string
	Username       string
	Description    string
	Age            string
}

type AccountResponse struct {
	Data struct {
		User struct {
			Result struct {
				Typename                   string `json:"__typename"`
				Id                         string `json:"id"`
				RestId                     string `json:"rest_id"`
				AffiliatesHighlightedLabel struct {
				} `json:"affiliates_highlighted_label"`
				HasGraduatedAccess bool   `json:"has_graduated_access"`
				IsBlueVerified     bool   `json:"is_blue_verified"`
				ProfileImageShape  string `json:"profile_image_shape"`
				Legacy             struct {
					CanDm               bool   `json:"can_dm"`
					CanMediaTag         bool   `json:"can_media_tag"`
					CreatedAt           string `json:"created_at"`
					DefaultProfile      bool   `json:"default_profile"`
					DefaultProfileImage bool   `json:"default_profile_image"`
					Description         string `json:"description"`
					Entities            struct {
						Description struct {
							Urls []interface{} `json:"urls"`
						} `json:"description"`
						Url struct {
							Urls []struct {
								DisplayUrl  string `json:"display_url"`
								ExpandedUrl string `json:"expanded_url"`
								Url         string `json:"url"`
								Indices     []int  `json:"indices"`
							} `json:"urls"`
						} `json:"url"`
					} `json:"entities"`
					FastFollowersCount      int           `json:"fast_followers_count"`
					FavouritesCount         int           `json:"favourites_count"`
					FollowersCount          int           `json:"followers_count"`
					FriendsCount            int           `json:"friends_count"`
					HasCustomTimelines      bool          `json:"has_custom_timelines"`
					IsTranslator            bool          `json:"is_translator"`
					ListedCount             int           `json:"listed_count"`
					Location                string        `json:"location"`
					MediaCount              int           `json:"media_count"`
					Name                    string        `json:"name"`
					NormalFollowersCount    int           `json:"normal_followers_count"`
					PinnedTweetIdsStr       []interface{} `json:"pinned_tweet_ids_str"`
					PossiblySensitive       bool          `json:"possibly_sensitive"`
					ProfileBannerUrl        string        `json:"profile_banner_url"`
					ProfileImageUrlHttps    string        `json:"profile_image_url_https"`
					ProfileInterstitialType string        `json:"profile_interstitial_type"`
					ScreenName              string        `json:"screen_name"`
					StatusesCount           int           `json:"statuses_count"`
					TranslatorType          string        `json:"translator_type"`
					Url                     string        `json:"url"`
					Verified                bool          `json:"verified"`
					WantRetweets            bool          `json:"want_retweets"`
					WithheldInCountries     []interface{} `json:"withheld_in_countries"`
				} `json:"legacy"`
				SmartBlockedBy        bool `json:"smart_blocked_by"`
				SmartBlocking         bool `json:"smart_blocking"`
				LegacyExtendedProfile struct {
				} `json:"legacy_extended_profile"`
				IsProfileTranslatable           bool `json:"is_profile_translatable"`
				HasHiddenSubscriptionsOnProfile bool `json:"has_hidden_subscriptions_on_profile"`
				VerificationInfo                struct {
				} `json:"verification_info"`
				HighlightsInfo struct {
					CanHighlightTweets bool   `json:"can_highlight_tweets"`
					HighlightedTweets  string `json:"highlighted_tweets"`
				} `json:"highlights_info"`
				BusinessAccount struct {
				} `json:"business_account"`
				CreatorSubscriptionsCount int `json:"creator_subscriptions_count"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type NitterResponse struct {
	JoinDate   string
	Followers  string
	AccountAge string
}
