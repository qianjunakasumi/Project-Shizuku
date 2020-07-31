/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : uehara
*   File Name    : action.go
*   File Path    : internal/uehara/
*   Author       : Qianjunakasumi
*   Description  : SHIZUKU 应用的定义和引用
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     actions []map[string]interface{} -- SHIZUKU 命令的定义
*
*   type expand struct -- SHIZUKU 命令的结构
*
*   func schedule() -- SHIZUKU 定时任务执行函数
*   func initApp()  -- 初始化 SHIZUKU 应用运行数据
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

package uehara

import (
	"github.com/qianjunakasumi/shizuku/configs"
	"github.com/qianjunakasumi/shizuku/internal/shizuku/twitter"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type expand struct {
	name     string
	display  string
	key      []string
	limit    []string
	required bool
	defaults string
}

var actions = []map[string]interface{}{
	{
		"name": "FetchTweet",
		"key":  []string{"查询推文", "推文", "获取推文", "推文查询", "推文获取"},
		"expand": []expand{{
			"account",
			"帐号",
			[]string{"帐号", "声优", "偶像"},
			twitter.GetKeys(),
			false,
			"前田",
		}},
		"func": &twitter.FetchTweets,
	},
	{
		"name": "FetchFollowersCount",
		"key":  []string{"查询粉丝", "粉丝", "获取粉丝", "粉丝查询", "粉丝获取"},
		"expand": []expand{{
			"account",
			"帐号",
			[]string{"帐号", "声优", "偶像"},
			twitter.GetKeys(),
			false,
			"前田",
		}},
		"func": &twitter.FetchFollowersCount,
	},
}

func schedule() {

	if configs.Conf.Development {
		return
	}

	c := cron.New()

	_, err := c.AddFunc("0 5 * * *", func() {
		m, err := twitter.ScheduleFollowersCount("前田")
		if err != nil {
			return
		}

		err = SendGroupMessage(1050964896, m)
		if err != nil {
			log.Error().Err(err)
		}
	})
	_, err = c.AddFunc("* * * * *", func() {
		m, err := twitter.ScheduleFetchTweets("lovelive")
		if err != nil {
			return
		}

		err = SendGroupMessage(289625710, m)
		if err != nil {
			log.Error().Err(err)
		}
	})
	_, err = c.AddFunc("* * * * *", func() {
		m, err := twitter.ScheduleFetchTweets("前田")
		if err != nil {
			return
		}

		err = SendGroupMessage(1050964896, m)
		if err != nil {
			log.Error().Err(err)
		}
	})

	if err != nil {
		log.Error().Err(err)
		return
	}

	c.Start()

}

func initApp() {

	twitter.Timer()
	schedule()

}
