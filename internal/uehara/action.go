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
*   along with this program.  If not, see https://github.com/qianjunakasumi/project-shizuku/blob/master/LICENSE.
*----------------------------------------------------------------------------------------------------------------------*/

package uehara

import (
	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku/llas"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku/meme"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku/shizuku"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku/twitter"
	"github.com/qianjunakasumi/project-shizuku/internal/uehara/messagechain"
	"github.com/qianjunakasumi/project-shizuku/pkg/database"

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

type action struct {
	name   string
	key    []string
	expand []expand
	fun    *func(calls map[string]string, info *messagechain.MessageInfo) (*messagechain.MessageChain, error)
}

var actions = []action{
	{
		"FetchTweet",
		[]string{"查询推文", "推文", "获取推文", "推文查询", "推文获取", "cxtw", "tw", "tweet"},
		[]expand{
			{
				"account",
				"帐号",
				[]string{"帐号", "声优", "偶像"},
				twitter.GetKeys(),
				false,
				"前田",
			},
			{
				"sequence",
				"次第",
				[]string{"第", "次第"},
				[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"},
				false,
				"1",
			},
		},
		&twitter.FetchTweets,
	},
	{
		"FetchFollowersCount",
		[]string{"查询粉丝", "粉丝", "获取粉丝", "粉丝查询", "粉丝获取", "cxfs", "fs", "fscx", "fans"},
		[]expand{{
			"account",
			"帐号",
			[]string{"帐号", "声优", "偶像"},
			twitter.GetKeys(),
			false,
			"前田",
		}},
		&twitter.FetchFollowersCount,
	},
	{
		"SendRandomMeme",
		[]string{"表情", "表情包", "随机表情", "随机表情包", "bq", "bqb", "meme"},
		[]expand{{
			"type",
			"种类",
			[]string{"种类", "类型", "偶像", "声优"},
			[]string{},
			false,
			"雫",
		}},
		&meme.SendRandomMeme,
	},
	{
		"UploadMeme",
		[]string{"上传表情", "表情上传", "scbq"},
		[]expand{},
		&meme.UploadMeme,
	},
	{
		"SendRandomStill",
		[]string{"来一张"},
		[]expand{{
			"type",
			"种类",
			[]string{"种类", "类型", "偶像", "声优"},
			[]string{},
			false,
			"上原步梦",
		}},
		&llas.RandomCard,
	},
	{
		"SHIZUKU",
		[]string{"S", "SHIZUKU", "s", "shizuku", "小雫"},
		[]expand{{
			"sentence",
			"句",
			[]string{"句", "话"},
			[]string{},
			true,
			"",
		}},
		&shizuku.TEST,
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

		SendGroupMessage(1050964896, m)
	})
	if err != nil {
		log.Error().Err(err)
		return
	}

	_, err = c.AddFunc("* * * * *", func() {
		m, err := twitter.ScheduleFetchTweets("lovelive")
		if err != nil {
			return
		}

		SendGroupMessage(289625710, m)
	})
	if err != nil {
		log.Error().Err(err)
		return
	}

	_, err = c.AddFunc("* * * * *", func() {
		m, err := twitter.ScheduleFetchTweets("前田")
		if err != nil {
			return
		}

		SendGroupMessage(1050964896, m)
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
	err := database.Connect()
	if err != nil {
		log.Error().Err(err).Msg("连接数据库时出错，请检查")
	}

}
