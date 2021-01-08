package twitter

type UserByScreenNameAPI struct {
	Data struct {
		User struct {
			Legacy struct {
				FollowersCount uint32 `json:"followers_count"`
			} `json:"legacy"`
		} `json:"user"`
	} `json:"data"`
}

type tweetAPIContent struct {
	GlobalObjects struct {
		Tweets map[string]tweetContent `json:"tweets"`
		Users  map[string]struct {
			IDStr string `json:"id_str"`
			Name  string `json:"name"`
		} `json:"users"`
	} `json:"globalObjects"`
}

type tweetContent struct {
	ConversationIDStr string `json:"conversation_id_str"`
	CreatedAt         string `json:"created_at"`
	FullText          string `json:"full_text"`
	Entities          struct {
		URLs []struct {
			ExpandedURL string `json:"expanded_url"`
			URL         string `json:"url"`
		} `json:"urls"`
		Media []struct {
			IDStr         string `json:"id_str"`
			MediaURLHttps string `json:"media_url_https"`
		} `json:"media"`
	} `json:"entities"`

	RetweetedStatusIDStr string `json:"retweeted_status_id_str"`
	QuotedStatusIDStr    string `json:"quoted_status_id_str"`
	InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`

	InReplyToUserIDStr string `json:"in_reply_to_user_id_str"`
	UserIDStr          string `json:"user_id_str"`
	FavoriteCount      uint32 `json:"favorite_count"`
}
