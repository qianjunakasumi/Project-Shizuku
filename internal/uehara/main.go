/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : uehara
*   File Name    : main.go
*   File Path    : internal/uehara/
*   Author       : Qianjunakasumi
*   Description  : 消息主入口命令处理
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   type Message map[string]interface{} -- Mirai API 消息的结构
*
*   func writeDefaults(expand2 []expand) map[string]string                                             -- 写入命令的默认值
*   func isExistKey(keys []expand, str string) string                                                  -- 判断是否存在字段
*   func isLimit(limit []string, str string) bool                                                      -- 判断是否超出限制
*   func writeCalls(fields *map[string]string, calls *[]string, action *map[string]interface{}) string -- 写入命令的输入值
*   func handle(calls *[]string, msg *Message, action *map[string]interface{})                         -- 处理命令
*   func receive(msg Message) error                                                                    -- 提取消息基本信息
*   func Connect() error                                                                               -- 启动 UEHARA
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
	"fmt"
	"strings"

	"github.com/qianjunakasumi/project-shizuku/internal/uehara/messagechain"

	"github.com/rs/zerolog/log"
)

// Message 消息
type Message map[string]interface{}

func writeDefaults(expand2 []expand) map[string]string {

	fields := make(map[string]string)

	for _, v := range expand2 {
		fields[v.name] = v.defaults
	}

	return fields

}

func isExistKey(keys []expand, str string) (string, int) {

	for k, v := range keys {
		for _, v2 := range v.key {
			if v2 == str {
				return v.name, k
			}
		}
	}

	return "", 0

}

func isLimit(limit []string, str string) bool {

	if len(limit) == 0 { // 无限制情形
		return true
	}

	for _, v := range limit {
		if v == str {
			return true
		}
	}

	return false

}

func writeCalls(fields *map[string]string, calls *[]string, expand2 []expand) string {

	i := 0

	for k := range *fields {

		if i == len(*calls) { // 如果索引等于传入的值长度时则退出，这时已无法取到传入的值
			break
		}

		callSplit := strings.Split((*calls)[i], "：") // 以“：”分割字符串，当长度为 2 时则是指定参数方式

		if len(callSplit) < 2 {

			if isLimit(expand2[i].limit, callSplit[0]) {
				(*fields)[k] = callSplit[0]
			} else {
				return "您输入的值不符合标准，我们推荐以下值：" + fmt.Sprintf("%v", expand2[i].limit)
			}

		} else {

			name, ii := isExistKey(expand2, callSplit[0])

			if name == "" {
				return "您输入的字段不符合标准，我们推荐以下字段：" + fmt.Sprintf("%v", expand2[i].key)
			}

			if isLimit(expand2[ii].limit, callSplit[1]) {
				(*fields)[name] = callSplit[1]
			} else {
				return "您输入的值不符合标准，我们推荐以下值：" + fmt.Sprintf("%v", expand2[ii].limit)
			}

		}

		i++

	}

	return ""

}

func handle(calls *[]string, msgInfo *messagechain.MessageInfo, action *action) {

	fields := writeDefaults(action.expand)
	errMsg := writeCalls(&fields, calls, action.expand)
	var msgChain *messagechain.MessageChain

	log.Info().Msg("查询详情：" + fmt.Sprintf("%v", fields))

	if errMsg != "" {

		msgChain = new(messagechain.MessageChain)
		msgChain.AddText(errMsg)

	} else {

		var err error

		// 取出函数指针的值，执行函数
		msgChain, err = (*(action.fun))(fields, msgInfo)
		if err != nil {
			msgChain.AddText("\n执行时发生错误，调试信息：" + fmt.Sprintf("%v", err))
		}

	}

	err := SendGroupMessage(msgInfo.GroupId, msgChain)
	if err != nil {
		return
	}

}

func receive(msg Message) {

	if msg["type"] != "GroupMessage" {
		return
	}

	msgChain := msg["messageChain"].([]interface{})

	if len(msgChain) < 2 {
		return
	}

	mainMsg := msgChain[1].(map[string]interface{})

	msgInfo := new(messagechain.MessageInfo)

	msgInfo.UserId = uint32((msg["sender"].(map[string]interface{}))["id"].(float64))
	msgInfo.UserName = (msg["sender"].(map[string]interface{}))["memberName"].(string)
	msgInfo.GroupId = uint32(((msg["sender"].(map[string]interface{}))["group"].(map[string]interface{}))["id"].(float64))
	msgInfo.GroupName = ((msg["sender"].(map[string]interface{}))["group"].(map[string]interface{}))["name"].(string)

	switch mainMsg["type"] {
	case "Plain":

		calls := strings.Fields(mainMsg["text"].(string))
		if len(calls) < 1 {
			return
		}

		for _, v := range actions {

			for _, v2 := range v.key {

				if calls[0] == v2 {

					calls2 := calls[1:]
					handle(&calls2, msgInfo, &v)

				}

			}

		}

	case "Image":

	}

}

// Connect 连接至Mirai
func Connect() error {

	if err := auth(); err != nil {
		return err
	}
	if err := verify(); err != nil {
		return err
	}
	if err := listen(); err != nil {
		return err
	}

	initApp()

	return nil

}
