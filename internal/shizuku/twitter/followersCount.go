/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : followersCount.go
*   File Path    : internal/shizuku/twitter/
*   Author       : Qianjunakasumi
*   Description  : 获取并解析 Twitter 粉丝数
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     FetchFollowersCount    -- 公开函数
*     ScheduleFollowersCount -- 公开函数
*
*   func main(id string) (string, error) -- 获取 Twitter 对象粉丝数
*
*   type pushScheduleFollowersCount struct -- 保存 Twitter 粉丝数计算数据和提供相关方法的容器
*     func (p *pushScheduleFollowersCount) getDatabaseData(id string) error -- 获取数据库保存的计算数据
*     func (p *pushScheduleFollowersCount) calcTwitterFollowersData()       -- 计算输出数据
*
*   func writeDB(a string, b string, c int32, d float32, e float32, f float32, g float32) -- 写入数据至数据库
*   func fetchFollowersCount(calls map[string]string) (*message.Chain, error)             -- 处理来自 Uehara 的调用
*   func scheduleFollowersCount(name string) (*message.Chain, error)                      -- 处理来自 定时任务函数 的调用
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
*   along with this program.  If not, see https://github.com/qianjunakasumi/project-shizuku/blob/master/LICENSE.
*----------------------------------------------------------------------------------------------------------------------*/

package twitter

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/uehara/message"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/database"

	_ "github.com/go-sql-driver/mysql" // 连接数据库需要的包
	"github.com/rs/zerolog/log"
)

type pushScheduleFollowersCount struct {
	yesterdayFollowersCount uint32  // 昨日粉丝数
	yuechuFollowersCount    uint32  // 月初粉丝数
	yesterdayHuanbiRate     float32 // 昨日环比率
	shangyueDingjiRate      float32 // 上月定基率
	todayFollowersCount     uint32  // 今日粉丝数
	newFollowersCount       int32   // 新增粉丝数
	toYesterdayHuanbi       float32 // 较昨日环比
	toYesterdayHuanbiRate   float32 // 较昨日环比率
	toYuechuDingji          float32 // 较本月月初定基
	toShangyueDingjiRate    float32 // 较上月定基率同比
}

var (
	FetchFollowersCount    = fetchFollowersCount    // 获取粉丝数
	ScheduleFollowersCount = scheduleFollowersCount // 定时粉丝数推送
)

func main(id string) (string, error) {

	res, err := get("graphql/-xfUfZsnR_zqjFd-IfrN5A/UserByScreenName?variables=%7B%22screen_name%22%3A%22" + id + "%22%2C%22withHighlightedLabel%22%3Atrue%7D")
	if err != nil {

		return "", err

	}

	count, ok := (((res["data"].(map[string]interface{}))["user"].(map[string]interface{}))["legacy"].(map[string]interface{}))["followers_count"].(float64)
	if !ok {

		log.Error().
			Str("内容", fmt.Sprintf("%v", res)).
			Msg("解析错误")

		return "", errors.New("错误")

	}

	return strconv.FormatFloat(count, 'f', 0, 64), nil

}

func fetchFollowersCount(calls map[string]string, info *message.Info) (*message.Chain, error) {

	var (
		m          = new(message.Chain)
		profile    = configs.AutoMatch(calls["account"], m, info.GroupName)
		count, err = main(profile.Followers)
	)

	if err != nil {

		return m, err

	}

	m.AddText("> " + profile.TwitterName + " 粉丝数：\n")

	countUint, err := strconv.ParseUint(count, 10, 32)
	if err != nil {

		return m, err

	}

	var (
		CountUint = float64(countUint) / 10000

		// 目的：提高精确度
		// 实现：输出千分位，阻止百分位非 9 数进一，删除最后一位小数
		CountString = strconv.FormatFloat(CountUint, 'f', 3, 64)
	)

	CountString = CountString[:len(CountString)-1]

	m.AddText(count + " — 约 " + CountString + " 万")

	return m, nil

}

func (p *pushScheduleFollowersCount) getDatabaseData(id string) error {

	db, err := sql.Open("mysql", configs.Conf.Databaseurl)
	if err != nil {

		return err

	}

	var (
		t         = time.Now()
		yesterday = t.AddDate(0, 0, -1).Format("20060102")
		yuechu    = t.AddDate(0, 0, 1-t.Day()).Format("20060102")
		shangyue  = t.AddDate(0, -1, 0).Format("20060102")

		sql1 = fmt.Sprintf(`SELECT twitter_followers.* FROM twitter_followers WHERE twitter_followers.id LIKE '%v' ORDER BY twitter_followers.date DESC`, yesterday+id)
		sql2 = fmt.Sprintf(`SELECT twitter_followers.* FROM twitter_followers WHERE twitter_followers.id LIKE '%v' ORDER BY twitter_followers.date DESC`, yuechu+id)
		sql3 = fmt.Sprintf(`SELECT twitter_followers.* FROM twitter_followers WHERE twitter_followers.id LIKE '%v' ORDER BY twitter_followers.date DESC`, shangyue+id)
	)

	rows1, err := db.Query(sql1)
	if err != nil {

		return err

	}

	rows2, err := db.Query(sql2)
	if err != nil {

		return err

	}

	rows3, err := db.Query(sql3)
	if err != nil {

		return err

	}

	var (
		yesterdayData = new(database.TwitterFollowers)
		yuechuData    = new(database.TwitterFollowers)
		shangyueData  = new(database.TwitterFollowers)
	)

	if rows1.Next() {

		err = rows1.Scan(&yesterdayData.ID,
			&yesterdayData.Date,
			&yesterdayData.Account,
			&yesterdayData.FollowersCount,
			&yesterdayData.NewFollowersCount,
			&yesterdayData.HuanbiRate,
			&yesterdayData.YdayHuanbiRate,
			&yesterdayData.DingjiRate,
			&yesterdayData.SyueDingjiRate)

		if err != nil {

			return err

		}

	}

	if rows2.Next() {

		err = rows2.Scan(&yuechuData.ID,
			&yuechuData.Date,
			&yuechuData.Account,
			&yuechuData.FollowersCount,
			&yuechuData.NewFollowersCount,
			&yuechuData.HuanbiRate,
			&yuechuData.YdayHuanbiRate,
			&yuechuData.DingjiRate,
			&yuechuData.SyueDingjiRate)

		if err != nil {

			return err

		}

	}

	if rows3.Next() {

		err = rows3.Scan(&shangyueData.ID,
			&shangyueData.Date,
			&shangyueData.Account,
			&shangyueData.FollowersCount,
			&shangyueData.NewFollowersCount,
			&shangyueData.HuanbiRate,
			&shangyueData.YdayHuanbiRate,
			&shangyueData.DingjiRate,
			&shangyueData.SyueDingjiRate)

		if err != nil {

			return err

		}

	}

	p.yesterdayFollowersCount = yesterdayData.FollowersCount // 昨日粉丝数
	p.yuechuFollowersCount = yuechuData.FollowersCount       // 月初粉丝数
	p.yesterdayHuanbiRate = yesterdayData.HuanbiRate         // 昨日环比率
	p.shangyueDingjiRate = shangyueData.DingjiRate           // 上月定基率

	defer db.Close()

	return nil

}

func (p *pushScheduleFollowersCount) calcTwitterFollowersData() {

	p.newFollowersCount = int32(p.todayFollowersCount) - int32(p.yesterdayFollowersCount) // 增加的粉丝数

	p.toYesterdayHuanbi = float32(p.newFollowersCount) / float32(p.yesterdayFollowersCount) * 1000 // 较昨日环比

	p.toYesterdayHuanbiRate = p.toYesterdayHuanbi - p.yesterdayHuanbiRate // 较昨日环比率

	if p.yuechuFollowersCount != 0 {

		p.toYuechuDingji = (float32(p.todayFollowersCount)/float32(p.yuechuFollowersCount) - 1) * 1000 // 较本月月初定基

	}

	if p.shangyueDingjiRate != 0 {

		p.toShangyueDingjiRate = p.toYuechuDingji - p.shangyueDingjiRate // 较上月定基率同比

	}

}

func writeDB(a string, b string, c int32, d float32, e float32, f float32, g float32) {

	var (
		t    = time.Now()
		id   = t.Format("20060102") + a
		date = t.Format("2006-01-02")
	)

	insert, err := database.DB.Prepare(`INSERT INTO twitter_followers ( id, date, account, followersCount, newFollowersCount, huanbiRate, ydayHuanbiRate, dingjiRate, syueDingjiRate ) VALUES (?,	?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {

		log.Error().Err(err)

		return

	}
	defer insert.Close()

	_, err = insert.Exec(id, date, a, b, c, d, e, f, g)
	if err != nil {

		log.Error().Err(err)

		return

	}

}

func scheduleFollowersCount(name string) (*message.Chain, error) {

	var (
		m          = new(message.Chain)
		profile    = configs.FuzzyGetProfile(name)
		count, err = main(profile.Followers)
	)

	if err != nil {

		return m, err

	}

	follwersCount, err := strconv.ParseUint(count, 10, 64)
	if err != nil {

		return m, nil

	}

	data := new(pushScheduleFollowersCount)

	err = data.getDatabaseData(profile.ID)
	if err != nil {

		return m, err

	}

	data.todayFollowersCount = uint32(follwersCount)

	data.calcTwitterFollowersData()

	go writeDB(profile.ID,
		count,
		data.newFollowersCount,
		data.toYesterdayHuanbi,
		data.toYesterdayHuanbiRate,
		data.toYuechuDingji,
		data.toShangyueDingjiRate,
	)

	m.AddText("> " + profile.TwitterName + " 粉丝数数据：\n")
	m.AddText("早上好！数据日报订阅详情：\n")
	m.AddText("总数：" + count + "\n")
	m.AddText("新增：" + strconv.FormatInt(int64(data.newFollowersCount), 10) + "\n")
	m.AddText("较昨日环比：" + strconv.FormatFloat(float64(data.toYesterdayHuanbi), 'f', 2, 32) + "\n")
	m.AddText("较昨日环比率：" + strconv.FormatFloat(float64(data.toYesterdayHuanbiRate), 'f', 2, 32) + "\n")
	m.AddText("较本月月初定基：" + strconv.FormatFloat(float64(data.toYuechuDingji), 'f', 2, 32) + "\n")
	m.AddText("较上月定基率同比：" + strconv.FormatFloat(float64(data.toShangyueDingjiRate), 'f', 2, 32))

	return m, nil

}
