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
*   func handlePlain(calls *[]string, msg *Message, action *map[string]interface{})                    -- 处理命令
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

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/uehara/message"

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

func writeCalls(
	fields *map[string]string, calls *[]string, expand2 []expand,
) string {

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

func handlePlain(mainMsg map[string]interface{}, msgInfo *message.MessageInfo) {

	var (
		a       bool // 指示是否匹配到命令，若 false 则退出
		calls   = strings.Fields(mainMsg["text"].(string))
		calls2  []string
		action2 = new(action)
	)

	if len(calls) < 1 {

		return

	}

	for _, v := range actions {

		for _, v2 := range v.key {

			if calls[0] == v2 {

				a = true
				calls2 = calls[1:]
				*action2 = v

			}

		}

	}

	if !a {

		return

	}

	var (
		fields   = writeDefaults(action2.expand)
		errMsg   = writeCalls(&fields, &calls2, action2.expand)
		msgChain *message.Chain
	)

	log.Info().
		Str("命令", action2.name).
		Str("参数", fmt.Sprintf("%v", fields)).
		Str("用户", msgInfo.UserName).
		Str("群组", msgInfo.GroupName).
		Msg("成功调用")

	if errMsg != "" {

		msgChain = new(message.Chain)
		msgChain.AddText(errMsg)

	} else {

		var err error

		// 取出函数指针的值，执行函数
		msgChain, err = (*(action2.fun))(fields, msgInfo)
		if err != nil {

			msgChain.AddText("\n\n=== :( 发生错误 ===\n调试信息：" + fmt.Sprintf("%v", err))

		}

	}

	sendGroupMessage(msgInfo.GroupId, msgChain)

}

func handleImage(msgInfo *message.MessageInfo, mainMsg map[string]interface{}) {

	if configs.ImageJob[msgInfo.UserId] == nil {

		return

	}

	url := mainMsg["url"].(string)

	var (
		fun = configs.ImageJob[msgInfo.UserId]
		m   *message.Chain
	)

	m, err := (*fun)(url, msgInfo)
	if err != nil {

		m.AddText("\n\n=== :( 发生错误 ===\n调试信息：" + fmt.Sprintf("%v", err))

	}

	sendGroupMessage(msgInfo.GroupId, m)

}

func handleQuote(mainMsg map[string]interface{}, msgChain []interface{}, msgInfo *message.MessageInfo) {

	if uint32(mainMsg["senderId"].(float64)) != configs.Conf.QQNumber {

		return

	}

	if configs.QuoteJob[msgInfo.UserId] == nil {

		return

	}

	var text string

	for i := 2; i < len(msgChain); i++ {

		if (msgChain[i].(map[string]interface{}))["type"] == "Plain" {

			text = (msgChain[i].(map[string]interface{}))["text"].(string)

		}

	}

	if text == "" {

		return

	}

	var (
		fun = configs.QuoteJob[msgInfo.UserId]
		m   *message.Chain
	)

	m, err := (*fun)(text, msgInfo)
	if err != nil {

		m.AddText("\n\n=== :( 发生错误 ===\n调试信息：" + fmt.Sprintf("%v", err))

	}

	sendGroupMessage(msgInfo.GroupId, m)

}

func receive(msg Message) {

	if msg["type"] != "GroupMessage" {

		return

	}

	msgChain := msg["messageChain"].([]interface{})

	if len(msgChain) < 2 {

		return

	}

	var (
		mainMsg = msgChain[1].(map[string]interface{})
		msgInfo = &message.MessageInfo{
			UserName:  (msg["sender"].(map[string]interface{}))["memberName"].(string),
			UserId:    uint32((msg["sender"].(map[string]interface{}))["id"].(float64)),
			GroupName: ((msg["sender"].(map[string]interface{}))["group"].(map[string]interface{}))["name"].(string),
			GroupId:   uint32(((msg["sender"].(map[string]interface{}))["group"].(map[string]interface{}))["id"].(float64)),
		}
	)

	// 若匿名用户则禁止
	if msgInfo.UserId == 80000000 {

		return

	}

	switch mainMsg["type"] {

	case "Plain":

		handlePlain(mainMsg, msgInfo)

	case "Image":

		handleImage(msgInfo, mainMsg)

	case "Quote":

		handleQuote(mainMsg, msgChain, msgInfo)

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
