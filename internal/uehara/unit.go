/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : uehara
*   File Name    : unit.go
*   File Path    : internal/uehara/
*   Author       : Qianjunakasumi
*   Description  : 与 Mirai 通讯的最小单元模块
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     session -- 令牌
*
*   func code(code float64) error                                                  -- 状态码解析
*   func auth() error                                                              -- 认证 Session
*   func verify() error                                                            -- 验证 Session
*   func Release() error                                                           -- 释放 Session
*   func listen() error                                                            -- 监听 WebSocket 消息
*   func SendGroupMessage(target uint32, message *messagechain.MessageChain) error -- 发送群消息
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
	"errors"

	"github.com/qianjunakasumi/shizuku/configs"
	"github.com/qianjunakasumi/shizuku/internal/uehara/messagechain"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

var session string

func code(code float64) error {

	if code == 0 {
		return nil
	}

	log.Error().
		Float64("状态码", code).
		Msg("发送消息失败")
	return errors.New("code不为0")

}

func auth() error {

	res, err := post("auth", Content{
		"authKey": configs.Conf.MiraiAuthKey,
	})
	if err != nil {
		return err
	}

	if err = code(res["code"].(float64)); err != nil {
		return err
	}
	session = res["session"].(string)

	return nil
}

func verify() error {

	res, err := post("verify", Content{
		"sessionKey": session,
		"qq":         configs.Conf.QQNumber,
	})
	if err != nil {
		return err
	}

	if err = code(res["code"].(float64)); err != nil {
		return err
	}

	return nil

}

// Release 释放SessionKey
func Release() error {

	res, err := post("release", Content{
		"sessionKey": session,
		"qq":         configs.Conf.QQNumber,
	})
	if err != nil {
		return err
	}

	if err = code(res["code"].(float64)); err != nil {
		return err
	}

	return nil

}

func listen() error {

	ws, err := websocket.Dial("ws://"+configs.Conf.MiraiAddress+"/message?sessionKey="+session, "", "http://localhost/")
	if err != nil {
		return err
	}

	go func() {
		for {
			msg := make(Message)

			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				return
			}
			if err = receive(msg); err != nil {
				return
			}
		}
	}()

	return nil

}

// SendGroupMessage 发送群消息
func SendGroupMessage(target uint32, message *messagechain.MessageChain) error {

	if message.Cancel {
		return nil
	}

	res, err := post("sendGroupMessage", Content{
		"sessionKey":   session,
		"target":       target,
		"messageChain": message.Content,
	})
	if err != nil {
		return err
	}

	if err = code(res["code"].(float64)); err != nil {
		return err
	}

	return nil

}
