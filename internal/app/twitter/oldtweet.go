package twitter

import (
	"errors"
	"html"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/kasumi"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/database"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/json"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/networkware"

	"github.com/rs/zerolog/log"
)

type fetchTwitter struct {
	orgin              *tweetAPIContent        // 原始
	tweetsList         map[string]tweetContent // 所有推文列表
	tweetsListIndex    []string                // 所有推文列表的排序索引
	wantTweetMap       tweetContent            // 要获取的推文的对象
	wantTweetText      string                  // 要获取的推文的内容
	wantTweetAddition  string                  // 要获取的推文的附加内容
	wantTweetHeader    string                  // 要获取的推文的标头
	wantTweetFooter    string                  // 要获取的推文的后缀
	wantTweetImagePath string                  // 要获取的推文的图片路径
}

// 获取推文 | 建立索引
func (f *fetchTwitter) main(id string, seq uint64) error {

	res := kasumi.New(&kasumi.Request{
		Addr:   "api.twitter.com/2/timeline/profile/" + id + ".json?include_profile_interstitial_type=1&include_blocking=1&include_blocked_by=1&include_followed_by=1&include_want_retweets=1&include_mute_edge=1&include_can_dm=1&include_can_media_tag=1&skip_status=1&cards_platform=Web-12&include_cards=1&include_composer_source=true&include_ext_alt_text=true&include_reply_count=1&tweet_mode=extended&include_entities=true&include_user_entities=true&include_ext_media_availability=true&send_error_codes=true&simple_quoted_tweet=true&include_tweet_replies=false&count=" + strconv.FormatUint(seq+1, 10) + "&ext=mediaStats%2ChighlightedLabel%2CcameraMoment",
		Method: "GET",
		Header: [][]string{
			{"authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"},
			{"x-guest-token", token},
		},
	}).TwitterReq(configs.GetProxyAddr())
	if res == nil {
		log.Error().Msg("请求推文失败：空指针 res *http.Response")
		return errors.New("请求推文失败：空指针 res *http.Response")
	}

	defer res.Body.Close()

	resss, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	data := &tweetAPIContent{}
	err = json.JSON.Unmarshal(resss, data)
	if err != nil {
		return err
	}

	f.orgin = data
	f.tweetsList = data.GlobalObjects.Tweets

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
func (f *fetchTwitter) writeTweet(which uint64) error {

	v := f.tweetsList[f.tweetsListIndex[which]]

	f.wantTweetMap = v
	f.wantTweetText = v.FullText

	return nil

}

// 判断并写入要获取的推文的标头 | 修正内容
func (f *fetchTwitter) writeHeader() {

	// 转推
	if retweetId := f.wantTweetMap.RetweetedStatusIDStr; retweetId != "" {
		header := f.wantTweetText[2 : strings.Index(f.wantTweetText, ":")+2]
		f.wantTweetHeader = "转推了" + header + "\n"
		f.wantTweetText = f.tweetsList[retweetId].FullText
	}

	// 引用推文
	if quoteId := f.wantTweetMap.QuotedStatusIDStr; quoteId != "" {
		f.wantTweetHeader = "并引用「" + f.orgin.GlobalObjects.Users[f.tweetsList[quoteId].UserIDStr].Name + "」的推文说：\n"
	}

	// 回复
	if replyId := f.wantTweetMap.InReplyToStatusIDStr; replyId != "" {
		f.wantTweetHeader = "回复了:" + "\n"

		a, ok := f.tweetsList[replyId]
		if !ok {

			f.wantTweetText += "\n给推文:\n被回复的推文已经飞往火星...太过久远啦"
			return

		}

		b := a.FullText

		f.wantTweetText += "\n给推文:\n" + b
	}

}

// 去除后缀链接 | 替换为原始链接 | 去除 http(s):// | 转换HTML转义符
func (f *fetchTwitter) tidyContent() {

	tweetURLs := f.wantTweetMap.Entities.URLs
	if len(tweetURLs) == 0 {
		// 针对无链接但存在引用例如转推或图片等扩展内容链接下的URL删除
		reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
		if err != nil {
			return
		}
		f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
		//f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")
	}

	for i := 0; i < len(tweetURLs); i++ {
		f.wantTweetText = strings.ReplaceAll(f.wantTweetText, tweetURLs[i].URL, tweetURLs[i].ExpandedURL)
		//f.wantTweetAddition = strings.ReplaceAll(f.wantTweetAddition, tweetURLs[i].URL, tweetURLs[i].ExpandedURL)
	}

	reg, err := regexp.Compile(`https://t.co/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	if err != nil {
		return
	}
	reg2, err := regexp.Compile(`https?://`)
	if err != nil {
		return
	}
	f.wantTweetText = reg.ReplaceAllString(f.wantTweetText, "")
	//f.wantTweetAddition = reg.ReplaceAllString(f.wantTweetAddition, "")
	f.wantTweetText = reg2.ReplaceAllString(f.wantTweetText, "")
	//f.wantTweetAddition = reg2.ReplaceAllString(f.wantTweetAddition, "")

	f.wantTweetText = html.UnescapeString(f.wantTweetText)
	//f.wantTweetAddition = html.UnescapeString(f.wantTweetAddition)

}

// 下载第一张图片缩略图
func (f *fetchTwitter) downloadImage() {
	tweetMedia := f.wantTweetMap.Entities.Media
	if len(tweetMedia) == 0 {
		return
	}

	address := strings.Builder{}
	address.WriteString("assets/images/temp/twitter/tweets/")
	address.WriteString(time.Now().Format("200601") + "/")
	address.WriteString(f.wantTweetMap.ConversationIDStr + "/")
	path := address.String()
	address.WriteString(tweetMedia[0].IDStr + ".webp")
	path2 := address.String()
	_, err := os.Stat(path)
	if err == nil {
		f.wantTweetImagePath = path2
		return
	}

	req := new(networkware.Networkware)
	req.Address = tweetMedia[0].MediaURLHttps + "?format=webp&name=small" // 缩略图
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

// 写入脚注信息
func (f *fetchTwitter) writeFooter() {

	var (
		t, _       = time.Parse("Mon Jan 02 15:04:05 +0000 2006", f.wantTweetMap.CreatedAt)
		beijing, _ = time.LoadLocation("Local")
	)

	f.wantTweetFooter = "🕒" + t.In(beijing).Format("15时04分")

}

func main2(twitter *fetchTwitter, message *shizuku.Message) {

	twitter.writeHeader()
	twitter.tidyContent()
	twitter.writeFooter()

	message.AddText(twitter.wantTweetHeader + twitter.wantTweetText + "\n")
	if twitter.wantTweetImagePath != "" {
		message.AddImage(twitter.wantTweetImagePath)
		message.AddText("\n")
	}
	/*
		if twitter.wantTweetAddition != "" {
			message.AddText("\n" + twitter.wantTweetAddition + "\n")
		}
	*/
	message.AddText(twitter.wantTweetFooter)

}

func fetchLastTweetID(id string) (string, error) {

	rows, err := database.DB.Query(`SELECT * FROM tweet_push WHERE tweet_push.id = ?`, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("无法找到记录")
	}

	data := new(database.TweetPush)

	err = rows.Scan(&data.ID, &data.TweetID)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(data.TweetID, 10), nil

}

func writeTweetID(id string, tweetID string) error {

	update, err := database.DB.Prepare(`UPDATE tweet_push SET tweet_push.tweetID = ? WHERE tweet_push.id = ?`)
	if err != nil {
		return err
	}
	defer update.Close()

	_, err = update.Exec(tweetID, id)
	if err != nil {

		return err

	}

	return nil

}

func scheduleFetchTweets(call string, sz *shizuku.SHIZUKU) (*shizuku.Message, error) {

	var (
		m       = shizuku.NewMsg()
		profile = sz.FuzzyGetIdol(call)
		x       = float64(time.Now().Hour())
		y       = profile.Twitter.Push(x)
		r       = rand.Intn(100)
	)

	if r > int(y) {
		return nil, nil
	}

	fetch := new(fetchTwitter)
	if err := fetch.main(profile.Twitter.Tweets, 1); err != nil {
		return nil, err
	}
	if err := fetch.writeTweet(0); err != nil {
		return nil, err
	}

	conversationID, err := fetchLastTweetID(profile.ID)
	if err != nil {
		return nil, err
	}

	if fetch.wantTweetMap.ConversationIDStr == conversationID {
		return nil, nil
	}

	err = writeTweetID(profile.ID, fetch.wantTweetMap.ConversationIDStr)
	if err != nil {
		return nil, err
	}

	m.AddText("「" + profile.SeiyuuName + "」说：\n")
	main2(fetch, m)

	return m, nil
}

type tweet struct{}

func (t tweet) OnCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	if c := qm.Call["idol"]; c != "_SHIZUKU默认检查专用" {
		qm.Type = sz.FuzzyGetIdol(c)
	}

	rm = shizuku.NewText("> " + qm.Type.SeiyuuName + " 的推文：\n")

	fetch := new(fetchTwitter)
	if err := fetch.main(qm.Type.Twitter.Tweets, 1); err != nil {
		return nil, err
	}
	if err := fetch.writeTweet(0); err != nil {
		return nil, err
	}
	main2(fetch, rm)

	return
}
