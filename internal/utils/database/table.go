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
	ID                uint32
	Date              string
	Account           uint8
	FollowersCount    uint32
	NewFollowersCount uint32
	HuanbiRate        float32
	YdayHuanbiRate    float32
	DingjiRate        float32
	SyueDingjiRate    float32
}
