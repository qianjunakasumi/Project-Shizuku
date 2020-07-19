/*
fetchTweets.go: 获取并解析Tweets
Copyright (C) 2020-present  QianjuNakasumi

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package twitter

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qianjunakasumi/shizuku/configs"
	"github.com/qianjunakasumi/shizuku/internal/uehara/messagechain"
	"github.com/qianjunakasumi/shizuku/pkg/networkware"
)

type fetchTwitter struct {
	tweetsList               map[string]interface{} // 所有推文列表
	tweetsListIndex          []string               // 所有推文列表的排序索引
	wantTweetMap             map[string]interface{} // 要获取的推文的对象
	wantTweetText            string                 // 要获取的推文的内容
	wantTweetAddition        string                 // 要获取的推文的附加内容
	wantTweetHeader          string                 // 要获取的推文的标头
	wantTweetFooter          string                 // 要获取的推文的后缀
	wantTweetImagePath       string                 // 要获取的推文的图片路径
	wantTweetTranslationText string                 // 要获取的推文的翻译内容
}

var FetchTweets = fetchTweet

// 获取推文 | 建立索引
func (f *fetchTwitter) main(id string) error {
	res, err := get("2/timeline/profile/" + id + ".json?include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&skip_status=1&cards_platform=Web-12&include_cards=1&include_composer_source=true&include_ext_alt_text=true&include_reply_count=1&tweet_mode=extended&include_entities=true&include_user_entities=true&include_ext_media_availability=true&send_error_codes=true&simple_quoted_tweet=true&include_tweet_replies=false&count=2&ext=mediaStats%2ChighlightedLabel%2CcameraMoment")
	if err != nil {
		return err
	}

	var ok bool
	if f.tweetsList, ok = res["globalObjects"].(map[string]interface{})["tweets"].(map[string]interface{}); !ok {
		return errors.New("解析推文时发生错误")
	}

	f.tweetsListIndex = make([]string, len(f.tweetsList))
	i := 0
	for k := range f.tweetsList {
		f.tweetsListIndex[i] = k
		i++
	}

	sort.Sort(sort.Reverse(sort.StringSlice(f.tweetsListIndex)))

	return nil
}

// 传入索引拉取推文对象和内容
func (f *fetchTwitter) writeTweet(which uint8) error {
	v, ok := f.tweetsList[f.tweetsListIndex[which]].(map[string]interface{})
	if !ok {
		return errors.New("ok is not true")
	}

	f.wantTweetMap = v
	f.wantTweetText = v["full_text"].(string)

	return nil
}

// 判断并写入要获取的推文的标头 | 修正内容
func (f *fetchTwitter) writeHeader() {

	// 转推
	if retweetId, ok := f.wantTweetMap["retweeted_status_id_str"].(string); ok {
		header := f.wantTweetText[2 : strings.Index(f.wantTweetText, ":")+2]
		f.wantTweetHeader = "转推了" + header + "\n"
		f.wantTweetText = f.tweetsList[retweetId].(map[string]interface{})["full_text"].(string)
	}

	// 带评论转推
	if quoteId, ok := f.wantTweetMap["quoted_status_id_str"].(string); ok {
		f.wantTweetHeader = "带评论:" + "\n"
		f.wantTweetAddition = "转推了:\n" + f.tweetsList[quoteId].(map[string]interface{})["full_text"].(string)
	}

	// 回复
	if replyId, ok := f.wantTweetMap["in_reply_to_status_id_str"].(string); ok {
		f.wantTweetHeader = "回复了:" + "\n"
		f.wantTweetText += "给推文:\n" + f.tweetsList[replyId].(map[string]interface{})["full_text"].(string)
	}

}

// 去除后缀链接 | 替换为原始链接 | 去除 http(s)://
func (f *fetchTwitter) urlHandle() error {
	tweetURLs, ok := (f.wantTweetMap["entities"].(map[string]interface{}))["urls"].([]interface{})
	if !ok {
		// 针对无链接但存在引用例如转推或图片等扩展内容链接下的URL删除
		reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
		if err != nil {
			return err
		}
		f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
		f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")

		return nil
	}

	for i := 0; i < len(tweetURLs); i++ {
		f.wantTweetText = strings.ReplaceAll(f.wantTweetText, (tweetURLs[i].(map[string]interface{}))["url"].(string), (tweetURLs[i].(map[string]interface{}))["expanded_url"].(string))
		f.wantTweetAddition = strings.ReplaceAll(f.wantTweetAddition, (tweetURLs[i].(map[string]interface{}))["url"].(string), (tweetURLs[i].(map[string]interface{}))["expanded_url"].(string))
	}

	reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	reg2, err := regexp.Compile(`https?://`)
	if err != nil {
		return err
	}
	f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
	f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")
	f.wantTweetText = reg2.ReplaceAllString(f.wantTweetText, "")
	f.wantTweetAddition = reg2.ReplaceAllString(f.wantTweetAddition, "")

	return nil
}

// 下载第一张图片缩略图
func (f *fetchTwitter) downloadImage() error {
	tweetMedia, ok := (f.wantTweetMap["entities"].(map[string]interface{}))["media"].([]interface{})
	if !ok {
		return nil
	}

	address := strings.Builder{}
	address.WriteString("assets/images/temp/twitter/tweets/")
	address.WriteString(time.Now().Format("200601") + "/")
	address.WriteString(f.wantTweetMap["conversation_id_str"].(string) + "/")
	path := address.String()
	address.WriteString(tweetMedia[0].(map[string]interface{})["id_str"].(string) + ".jpg")
	path2 := address.String()
	f.wantTweetImagePath = path2
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	req := new(networkware.Networkware)
	req.Address = tweetMedia[0].(map[string]interface{})["media_url_https"].(string) + "?format=jpg&name=small" // 缩略图
	req.Method = "GET"
	req.Header = [][]string{
		{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4128.3 Safari/537.36"},
	}
	req.Proxy = "http://127.0.0.1:10809"
	res, err := req.Send()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(path2)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, bytes.NewReader(cont))
	if err != nil {
		return err
	}

	return nil
}

// 翻译推文
func (f *fetchTwitter) translateTweet() error {
	content := f.wantTweetText
	reg, err := regexp.Compile(`[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	if err != nil {
		return err
	}
	content = reg.ReplaceAllString(content, "")

	query := url.Values{}
	query.Add("q", content)
	query.Add("from", "jp")
	query.Add("to", "zh")
	query.Add("appid", configs.Conf.TranslationAppID)
	query.Add("salt", "SHIZUKU")
	h := md5.New()
	_, err = io.WriteString(h, configs.Conf.TranslationAppID+content+"SHIZUKU"+configs.Conf.TranslationKey)
	if err != nil {
		return err
	}
	query.Add("sign", fmt.Sprintf("%x", h.Sum(nil)))

	res, err := http.Get("https://api.fanyi.baidu.com/api/trans/vip/translate?" + query.Encode())
	if err != nil {
		return err
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(robots, &result); err != nil {
		return err
	}
	transList, ok := result["trans_result"].([]interface{})
	if !ok {
		return errors.New("解析翻译时出问题")
	}
	for i := 0; i < len(transList); i++ {
		f.wantTweetTranslationText += transList[i].(map[string]interface{})["dst"].(string) + "\n"
	}

	return nil
}

func (f *fetchTwitter) writeFooter() {
	t, _ := time.Parse("Mon Jan 02 15:04:05 +0000 2006", f.wantTweetMap["created_at"].(string))
	beijing, _ := time.LoadLocation("Local")
	favoriteCount := strconv.FormatFloat(f.wantTweetMap["favorite_count"].(float64), 'f', 0, 64)

	f.wantTweetFooter = "发送时间：" + t.In(beijing).Format("01月02日 15时04分") + "\n被喜欢次数：" + favoriteCount
}

func fetchTweet(calls map[string]string) (*messagechain.MessageChain, error) {
	m := new(messagechain.MessageChain)
	profile := getProfile(calls["account"])

	m.AddText("> " + profile.name + " 的推文：\n")

	fetch := new(fetchTwitter)
	if err := fetch.main(profile.tweets); err != nil {
		return m, err
	}
	if err := fetch.writeTweet(0); err != nil {
		return m, err
	}
	fetch.writeHeader()
	if err := fetch.urlHandle(); err != nil {
		return m, err
	}
	if err := fetch.downloadImage(); err != nil {
		return m, err
	}
	if err := fetch.translateTweet(); err != nil {
		return m, err
	}
	fetch.writeFooter()

	m.AddText(fetch.wantTweetHeader + fetch.wantTweetText + "\n")
	if fetch.wantTweetImagePath != "" {
		m.AddImage(fetch.wantTweetImagePath)
		m.AddText("\n")
	}
	if fetch.wantTweetAddition != "" {
		m.AddText("\n" + fetch.wantTweetAddition + "\n")
	}
	m.AddText("\n翻译：\n" + fetch.wantTweetTranslationText + "\n")
	m.AddText(fetch.wantTweetFooter)

	return m, nil
}
