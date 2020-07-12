/*
unit.go: 与Mirai通讯的最小单元模块
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

package uehara

import (
	"errors"
	"fmt"

	"github.com/qianjunakasumi/shizuku/configs"
	"github.com/qianjunakasumi/shizuku/internal/uehara/messageChain"

	"golang.org/x/net/websocket"
)

var (
	session string
)

func code(code float64) error {
	if code == 0 {
		return nil
	}

	fmt.Println(code)
	return errors.New("code不为0")
}

func auth() error {
	res, err := post("auth", Content{
		"authKey": configs.MiraiAuthKey,
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
		"qq":         configs.QQNumber,
	})
	if err != nil {
		return err
	}

	if err = code(res["code"].(float64)); err != nil {
		return err
	}

	return nil
}

// 释放SessionKey
func Release() error {
	res, err := post("release", Content{
		"sessionKey": session,
		"qq":         configs.QQNumber,
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
	ws, err := websocket.Dial("ws://"+configs.MiraiAddress+"/message?sessionKey="+session, "", "http://localhost/")
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

// 发送群消息
func SendGroupMessage(target uint32, message *messageChain.MessageChain) error {
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
