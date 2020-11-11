package database

type MemeUploadAgreeer struct {
	QQID      uint64
	AgreeInfo string
}

type TweetPush struct {
	ID      uint8
	TweetID uint64
}

type TwitterFollowers struct {
	ID             uint32
	Date           string
	Account        uint8
	FollowersCount uint32
}
