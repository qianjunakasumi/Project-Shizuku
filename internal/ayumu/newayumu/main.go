/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : ayumu
*   File Name    : main.go
*   File Path    : internal/ayumu/
*   Author       : Qianjunakasumi
*   Description  : AYUMU 主要功能
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

package newayumu

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/rs/zerolog/log"
	"time"
)

// Ayumu Ayumu // TODO 更多
type Ayumu struct {
	id uint64 // QQ 帐号
	ps string // QQ 密码
	c  *client.QQClient
}

// New 新增AYUMU
func New(i uint64, p string) {

	err := client.SystemDeviceInfo.ReadJson([]byte("{\"display\":\"MIRAI.666470.001\",\"product\":\"mirai\",\"device\":\"mirai\",\"board\":\"mirai\",\"model\":\"mirai\",\"finger_print\":\"mamoe/mirai/mirai:10/MIRAI.200122.001/5696651:user/release-keys\",\"boot_id\":\"58fe8ac7-4de7-71ec-073d-07eb3187a533\",\"proc_version\":\"Linux version 3.0.31-HxHC3WtY (android-build@xxx.xxx.xxx.xxx.com)\",\"imei\":\"351912693210254\"}"))
	if err != nil {
		log.Panic().Err(err).Msg("设置设备信息失败")
	}

	c := client.NewClient(int64(i), p)
	c.OnLog(func(q *client.QQClient, e *client.LogEvent) {
		switch e.Type {

		case "INFO":
			log.Info().Str("信息", e.Message).Msg("协议")

		case "ERROR":
			log.Error().Str("信息", e.Message).Msg("协议")
		}
	})

	a := &Ayumu{id: i, ps: p, c: c}
	a.login()

	go func() {
		c.OnDisconnected(func(q *client.QQClient, e *client.ClientDisconnectedEvent) {

			for {
				log.Warn().Msg("啊哦掉线了，准备重连中...2s")
				time.Sleep(time.Second * 2)
				a.login()

				return
			}

		})
	}()

}

func (a Ayumu) login() {

	r, err := a.c.Login()
	if err != nil {
		log.Panic().Err(err).Msg("登录失败")
	}
	if !r.Success {
		log.Panic().Str("错误", r.ErrorMessage).Msg("登录失败")
	}

	log.Info().Msg("登录成功：" + a.c.Nickname)

	err = a.c.ReloadGroupList()
	if err != nil {
		log.Panic().Err(err).Msg("加载群列表失败")
	}

	err = a.c.ReloadFriendList()
	if err != nil {
		log.Panic().Err(err).Msg("加载好友列表失败")
	}

	log.Info().Int("个数", len(a.c.GroupList)).Msg("加载群列表成功")
	log.Info().Int("个数", len(a.c.FriendList)).Msg("加载好友列表成功")

}
