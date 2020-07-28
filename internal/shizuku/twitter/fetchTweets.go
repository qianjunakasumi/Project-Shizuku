/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : fetchTweets.go
*   File Path    : internal/shizuku/twitter/
*   Author       : Qianjunakasumi
*   Description  : 获取并解析 Tweet
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     FetchTweets func(calls map[string]string) -- 公开函数
*     ScheduleFetchTweets func(call string)     -- 公开函数
*     conversationIdList map[string]string      -- 存储 Tweet Conversation ID
*
*   type fetchTwitter struct                               -- 保存 Tweet 相关信息和提供相关方法的容器
*     func (f *fetchTwitter) main(id string) error         -- 获取和保存输入的帐号的最新两条 Tweet 并按时间倒序排序生成其索引
*     func (f *fetchTwitter) writeTweet(which uint8) error -- 获取指定定位的 Tweet 并写入
*     func (f *fetchTwitter) writeHeader()                 -- 判断 Tweet 类型，套用对应模板，重新生成标准 Tweet 内容
*     func (f *fetchTwitter) tidyContent()                 -- 转换短链接为原始链接，删除 Tweet 最后可能存在的链接，去除所有
*                                                             http(s):// 前缀，转义 HTML 转义符
*     func (f *fetchTwitter) downloadImage()               -- 下载 Tweet 包含的图片的第一张图片缩略图
*     func (f *fetchTwitter) translateTweet()              -- 调用 百度翻译API 翻译 Tweet
*     func (f *fetchTwitter) writeFooter()                 -- 添加 Tweet 创建时间和被喜欢次数
*
*   func main2(twitter *fetchTwitter, message *messagechain.MessageChain)        -- fetchTwitter 的部分封装
*   func scheduleFetchTweets(call string) (*messagechain.MessageChain, error)    -- 处理来自 定时任务函数 的调用
*   func fetchTweet(calls map[string]string) (*messagechain.MessageChain, error) -- 处理来自 Uehara 的调用
*
*----------------------------------------------------------------------------------------------------------------------*
* Copyright:
*
*   Copyright (C) 2020-present QianjuNakasumi
*
*   E-mail qianjunakasumi@gmail.com
*
*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU Affero General Public License as published
*   by the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.
*
*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU Affero General Public License for more details.
*
*   You should have received a copy of the GNU Affero General Public License
*   along with this program.  If not, see https://github.com/qianjunakasumi/shizuku/blob/master/LICENSE.
*----------------------------------------------------------------------------------------------------------------------*/

package twitter

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"math/rand"
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
	"github.com/rs/zerolog/log"
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

var (
	FetchTweets         = fetchTweet              // 获取推文
	ScheduleFetchTweets = scheduleFetchTweets     // 定时推送
	conversationIdList  = make(map[string]string) // 推文ID列表
)

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
		f.wantTweetAddition = "\n转推了:\n" + f.tweetsList[quoteId].(map[string]interface{})["full_text"].(string)
	}

	// 回复
	if replyId, ok := f.wantTweetMap["in_reply_to_status_id_str"].(string); ok {
		f.wantTweetHeader = "回复了:" + "\n"
		f.wantTweetText += "\n给推文:\n" + f.tweetsList[replyId].(map[string]interface{})["full_text"].(string)
	}

}

// 去除后缀链接 | 替换为原始链接 | 去除 http(s):// | 转换HTML转义符
func (f *fetchTwitter) tidyContent() {
	tweetURLs, ok := (f.wantTweetMap["entities"].(map[string]interface{}))["urls"].([]interface{})
	if !ok {
		// 针对无链接但存在引用例如转推或图片等扩展内容链接下的URL删除
		reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
		if err != nil {
			return
		}
		f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
		f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")
	}

	for i := 0; i < len(tweetURLs); i++ {
		f.wantTweetText = strings.ReplaceAll(f.wantTweetText, (tweetURLs[i].(map[string]interface{}))["url"].(string), (tweetURLs[i].(map[string]interface{}))["expanded_url"].(string))
		f.wantTweetAddition = strings.ReplaceAll(f.wantTweetAddition, (tweetURLs[i].(map[string]interface{}))["url"].(string), (tweetURLs[i].(map[string]interface{}))["expanded_url"].(string))
	}

	reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	reg2, err := regexp.Compile(`https?://`)
	if err != nil {
		return
	}
	f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
	f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")
	f.wantTweetText = reg2.ReplaceAllString(f.wantTweetText, "")
	f.wantTweetAddition = reg2.ReplaceAllString(f.wantTweetAddition, "")

	f.wantTweetText = html.UnescapeString(f.wantTweetText)
	f.wantTweetAddition = html.UnescapeString(f.wantTweetAddition)
}

// 下载第一张图片缩略图
func (f *fetchTwitter) downloadImage() {
	tweetMedia, ok := (f.wantTweetMap["entities"].(map[string]interface{}))["media"].([]interface{})
	if !ok {
		return
	}

	address := strings.Builder{}
	address.WriteString("assets/images/temp/twitter/tweets/")
	address.WriteString(time.Now().Format("200601") + "/")
	address.WriteString(f.wantTweetMap["conversation_id_str"].(string) + "/")
	path := address.String()
	address.WriteString(tweetMedia[0].(map[string]interface{})["id_str"].(string) + ".webp")
	path2 := address.String()
	_, err := os.Stat(path)
	if err == nil {
		f.wantTweetImagePath = path2
		return
	}

	req := new(networkware.Networkware)
	req.Address = tweetMedia[0].(map[string]interface{})["media_url_https"].(string) + "?format=webp&name=small" // 缩略图
	req.Method = "GET"
	req.Header = [][]string{
		{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4128.3 Safari/537.36"},
	}
	req.Proxy = "http://127.0.0.1:10809"
	res, err := req.Send()
	if err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "downloadImage").
			Msg("请求图片时出错")
		return
	}
	defer res.Body.Close()

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "downloadImage").
			Msg("创建缓存文件夹时出错")
		return
	}
	file, err := os.Create(path2)
	if err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "downloadImage").
			Msg("创建缓存图片时出错")
		return
	}
	defer file.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "downloadImage").
			Msg("保存缓存图片时出错")
		return
	}

	f.wantTweetImagePath = path2
}

// 翻译推文
func (f *fetchTwitter) translateTweet() {
	content := f.wantTweetText
	reg, err := regexp.Compile(`[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	if err != nil {
		return
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
		return
	}
	query.Add("sign", fmt.Sprintf("%x", h.Sum(nil)))

	res, err := http.Get("https://api.fanyi.baidu.com/api/trans/vip/translate?" + query.Encode())
	if err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "translateTweet").
			Msg("请求翻译API时出错")
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	result := make(map[string]interface{})
	if err := json.Unmarshal(robots, &result); err != nil {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "translateTweet").
			Msg("解析JSON时出错")
		return
	}
	transList, ok := result["trans_result"].([]interface{})
	if !ok {
		log.Warn().
			Str("包名", "twitter").
			Str("方法", "translateTweet").
			Msg("解析翻译内容时出错")
		return
	}
	for i := 0; i < len(transList); i++ {
		f.wantTweetTranslationText += transList[i].(map[string]interface{})["dst"].(string) + "\n"
	}
}

// 写入脚注信息
func (f *fetchTwitter) writeFooter() {
	t, _ := time.Parse("Mon Jan 02 15:04:05 +0000 2006", f.wantTweetMap["created_at"].(string))
	beijing, _ := time.LoadLocation("Local")
	favoriteCount := strconv.FormatFloat(f.wantTweetMap["favorite_count"].(float64), 'f', 0, 64)

	f.wantTweetFooter = "发送时间：" + t.In(beijing).Format("01月02日 15时04分") + "\n被喜欢次数：" + favoriteCount
}

func main2(twitter *fetchTwitter, message *messagechain.MessageChain) {
	twitter.writeHeader()
	twitter.tidyContent()
	twitter.downloadImage()
	twitter.translateTweet()
	twitter.writeFooter()

	message.AddText(twitter.wantTweetHeader + twitter.wantTweetText + "\n")
	if twitter.wantTweetImagePath != "" {
		message.AddImage(twitter.wantTweetImagePath)
		message.AddText("\n")
	}
	if twitter.wantTweetAddition != "" {
		message.AddText("\n" + twitter.wantTweetAddition + "\n")
	}
	message.AddText("\n翻译：\n" + twitter.wantTweetTranslationText + "\n")
	message.AddText(twitter.wantTweetFooter)
}

func scheduleFetchTweets(call string) (*messagechain.MessageChain, error) {
	m := new(messagechain.MessageChain)
	profile := getProfile(call)

	x := float64(time.Now().Hour())
	y := profile.push(x)
	r := rand.Intn(100)
	if r > int(y) {
		m.Cancel = true
		return m, nil
	}

	fetch := new(fetchTwitter)
	if err := fetch.main(profile.tweets); err != nil {
		return m, err
	}
	if err := fetch.writeTweet(0); err != nil {
		return m, err
	}
	if fetch.wantTweetMap["conversation_id_str"].(string) == conversationIdList[profile.name] {
		m.Cancel = true
		return m, nil
	}

	conversationIdList[profile.name] = fetch.wantTweetMap["conversation_id_str"].(string)
	m.AddText("推文推送服务 > " + profile.name + " 的推文：\n")
	main2(fetch, m)

	return m, nil
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
	main2(fetch, m)

	return m, nil
}
