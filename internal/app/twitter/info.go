/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : info.go
*   File Path    : internal/app/twitter/
*   Author       : Qianjunakasumi
*   Description  : 应用信息注册
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
	"time"

	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"

	"github.com/rs/zerolog/log"
)

func init() {

	Timer()

	shizuku.NewTask(&shizuku.AppTaskInfo{
		Name:    "定时获取推文",
		Spec:    "* * * * *",
		Pointer: new(test),
	})

	shizuku.NewTask(&shizuku.AppTaskInfo{
		Name:    "定时推送粉丝数",
		Spec:    "0 5 * * *",
		Pointer: new(test2),
	})

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "FetchTweet",
		DisplayName: "获取推文",
		Keys:        []string{"推文", "推特推文", "tweet", "tw"},
		Expand: shizuku.Expand{
			{
				"idol",
				"偶像",
				[]string{"偶像", "爱抖露"},
				[]string{},
				false,
				"_SHIZUKU默认检查专用",
			},
		},
		Pointer: new(tweet),
	})

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "FetchFollowersCount",
		DisplayName: "获取推特粉丝",
		Keys:        []string{"粉丝", "推特粉丝", "fans", "fs"},
		Expand: shizuku.Expand{
			{
				"idol",
				"偶像",
				[]string{"偶像", "爱抖露"},
				[]string{},
				false,
				"_SHIZUKU默认检查专用",
			},
		},
		Pointer: new(followers),
	})

}

type test struct{}

func (t test) OnTaskCall(sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	pushList := []struct {
		name    string
		targets []uint64
	}{
		{"ラブライブ", []uint64{289625710}},
		{"前田", []uint64{1050964896}},
		{"大西亜", []uint64{296973163}},
		{"小泉", []uint64{641985475}},
		{"相良", []uint64{522730499}},
	}

	for _, v := range pushList {

		ss, err := scheduleFetchTweets(v.name, sz)
		if err != nil {
			log.Error().Err(err).Msg("推送推文时出错")
			continue
		}
		if ss != nil {

			for _, vv := range v.targets {
				sz.Rina.SendGroupMsg(ss.To(vv))
			}

		}
		time.Sleep(time.Second)

	}

	return

}
