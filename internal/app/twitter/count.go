package twitter

import (
	"errors"
	"fmt"
	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/database"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/qianjunakasumi/project-shizuku/internal/kasumi"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/json"
)

type followers struct{}

func main(i string) (uint32, error) {

	res := kasumi.New(&kasumi.Request{
		Addr:   "api.twitter.com/graphql/4S2ihIKfF3xhp-ENxvUAfQ/UserByScreenName?variables=%7B%22screen_name%22%3A%22" + i + "%22%2C%22withHighlightedLabel%22%3Atrue%7D",
		Method: "GET",
		Header: [][]string{
			{"authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"},
			{"x-guest-token", token},
		},
	}).TwitterReq(configs.GetProxyAddr())
	if res == nil {
		return 0, errors.New("请求数据出错")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	api := &UserByScreenNameAPI{}
	err = json.JSON.Unmarshal(b, api)
	if err != nil {
		return 0, err
	}

	return api.Data.User.Legacy.FollowersCount, nil

}

func (f followers) OnCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	if c := qm.Call["idol"]; c != "_SHIZUKU默认检查专用" {
		qm.Type = sz.FuzzyGetIdol(c)
	}

	res, err := main(qm.Type.Twitter.Followers)
	if err != nil {
		return
	}

	// 提高精确度：输出千分位，阻止百分位非 9 数进一，删除最后一位小数
	CountString := strconv.FormatFloat(float64(res)/10000, 'f', 3, 64)

	rm = shizuku.NewText("> " + qm.Type.SeiyuuName + " 粉丝数：\n").
		AddText(strconv.FormatUint(uint64(res), 10)).
		AddText(" - 约 " + CountString[:len(CountString)-1] + " 万")

	return

}

type pushScheduleFollowersCount struct {
	yesterdayFollowersCount uint32 // 昨日粉丝数
	todayFollowersCount     uint32 // 今日粉丝数
	newFollowersCount       int32  // 新增粉丝数
}

func (p *pushScheduleFollowersCount) getDatabaseData(id string) error {

	var (
		t         = time.Now()
		yesterday = t.AddDate(0, 0, -1).Format("20060102")

		sql1 = fmt.Sprintf(`SELECT twitter_followers.* FROM twitter_followers WHERE twitter_followers.id LIKE '%v' ORDER BY twitter_followers.date DESC`, yesterday+id)
	)

	rows1, err := database.DB.Query(sql1)
	if err != nil {
		return err
	}

	yesterdayData := new(database.TwitterFollowers)
	if rows1.Next() {

		err = rows1.Scan(&yesterdayData.ID,
			&yesterdayData.Date,
			&yesterdayData.Account,
			&yesterdayData.FollowersCount,
		)

		if err != nil {
			return err
		}

	}

	p.yesterdayFollowersCount = yesterdayData.FollowersCount // 昨日粉丝数

	return nil

}

func (p *pushScheduleFollowersCount) calcTwitterFollowersData() {

	p.newFollowersCount = int32(p.todayFollowersCount) - int32(p.yesterdayFollowersCount) // 增加的粉丝数

}

func writeDB(a string, b string) {

	var (
		t    = time.Now()
		id   = t.Format("20060102") + a
		date = t.Format("2006-01-02")
	)

	insert, err := database.DB.Prepare(`INSERT INTO twitter_followers ( id, date, account, followersCount ) VALUES (?, ?, ?, ?)`)
	if err != nil {

		log.Error().Err(err).Msg("插入数据失败")
		return

	}
	defer insert.Close()

	_, err = insert.Exec(id, date, a, b)
	if err != nil {

		log.Error().Err(err).Msg("插入数据失败")
		return

	}

}

type test2 struct{}

func (t test2) OnTaskCall(sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	pushList := []struct {
		name   string
		target uint64
	}{
		{"前田", 1050964896},
		{"大西亜", 296973163},
		{"相良", 522730499},
	}

	for _, v := range pushList {

		profile := sz.FuzzyGetIdol(v.name)
		count, err := main(profile.Twitter.Followers)
		if err != nil {
			log.Error().Err(err).Msg("粉丝数推送错误")
			continue
		}

		follwersCount := uint64(count)
		data := new(pushScheduleFollowersCount)

		err = data.getDatabaseData(profile.ID)
		if err != nil {
			log.Error().Err(err).Msg("粉丝数推送错误")
			continue
		}

		data.todayFollowersCount = uint32(follwersCount)

		data.calcTwitterFollowersData()

		countStr := strconv.Itoa(int(count))

		go writeDB(profile.ID, countStr)

		rm = shizuku.NewText("> " + profile.SeiyuuName + " 粉丝数：\n").
			AddText("总数：" + countStr + "\n").
			AddText("变幅：" + strconv.FormatInt(int64(data.newFollowersCount), 10) + "\n")

		sz.Rina.SendGroupMsg(rm.To(v.target))

	}

	return nil, nil

}
