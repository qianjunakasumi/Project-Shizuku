package twitter

import (
	"github.com/qianjunakasumi/project-shizuku/internal/utils/json"
	"sort"
)

// tweetAPI Tweet API 解析器
type tweetAPI struct {
	data tweetAPIContent
}

// 新建 Tweet API 内容解析器
func newTweetAPIParser(d []byte) (con *tweetAPI, err error) {

	con = new(tweetAPI)
	err = json.JSON.Unmarshal(d, &con.data)
	if err != nil {
		return
	}

	return
}

// getOnlyMeTweetList 获取只有输入偶像的所有Tweet
func (t tweetAPI) getOnlyMeTweetList(i string) (l []string) {

	for k, v := range t.data.GlobalObjects.Tweets {
		if v.UserIDStr == i {
			l = append(l, k)
		}
	}

	return

}

// TODO 拉取更多信息，比如回复推文，转推等，通过字段获取

// getTweet 获取输入的偶像的第n条 Tweet
func (t tweetAPI) getTweet(i string, n uint8) *tweetContent {

	list := t.getOnlyMeTweetList(i)
	sort.Sort(sort.Reverse(sort.StringSlice(list)))

	tweet := t.data.GlobalObjects.Tweets[list[n]]

	return &tweet

}

// 获取文本
func (t tweetAPI) getText(n uint8) string { return t.getTweet("347849994", n).FullText }
