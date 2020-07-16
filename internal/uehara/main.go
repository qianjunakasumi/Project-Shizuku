/*
main.go: 主要入口 | 消息解析校验和路由
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
	"fmt"
	"strings"

	"github.com/qianjunakasumi/shizuku/internal/uehara/messagechain"

	"github.com/rs/zerolog/log"
)

type Message map[string]interface{}

func writeDefaults(expand2 []expand) map[string]string {
	fields := make(map[string]string)

	for _, v := range expand2 {
		fields[v.name] = v.defaults
	}

	return fields
}

func isExistKey(keys []expand, str string) string {
	for _, v := range keys {
		for _, v2 := range v.key {
			if v2 == str {
				return v.name
			}
		}
	}

	return ""
}

func isLimit(limit []string, str string) bool {
	if len(limit) == 0 { //无限制情形
		return true
	}

	for _, v := range limit {
		if v == str {
			return true
		}
	}

	return false
}

func writeCalls(fields *map[string]string, calls *[]string, action *map[string]interface{}) string {
	callsLen := len(*calls)
	i := 0

	for k := range *fields {
		if i == callsLen { // 如果索引等于传入的值长度时则退出，这时已无法取到传入的值
			break
		}

		callSplit := strings.Split((*calls)[i], "：") //以“：”分割字符串，当长度为 2 时则是指定参数方式

		expand2 := (*action)["expand"].([]expand)

		if len(callSplit) < 2 {
			if isLimit(expand2[i].limit, callSplit[0]) {
				(*fields)[k] = callSplit[0]
			} else {
				return "您输入的值不符合标准，我们推荐以下值：" + fmt.Sprintf("%v", expand2[i].limit)
			}
		} else {
			name := isExistKey(expand2, callSplit[0])
			if name == "" {
				return "您输入的字段不符合标准，我们推荐以下字段：" + fmt.Sprintf("%v", expand2[i].key)
			}

			if isLimit(expand2[i].limit, callSplit[1]) {
				(*fields)[name] = callSplit[1]
			} else {
				return "您输入的值不符合标准，我们推荐以下值：" + fmt.Sprintf("%v", expand2[i].limit)
			}
		}

		i++
	}

	return ""
}

func handle(calls *[]string, msg *Message, action *map[string]interface{}) {
	fields := writeDefaults((*action)["expand"].([]expand))
	errMsg := writeCalls(&fields, calls, action)
	var msgChain *messagechain.MessageChain

	fmt.Println("查询的数据：", fields)

	if errMsg != "" {
		msgChain = new(messagechain.MessageChain)
		msgChain.AddText(errMsg)
	} else {
		var err error

		// 取出(*actions)的值，定位"func"字段，类型断言为相同的函数类型,取出函数指针的值，执行函数
		msgChain, err = (*((*action)["func"].(*func(calls map[string]string) (*messagechain.MessageChain, error))))(fields)
		if err != nil {
			msgChain.AddText("\n执行时发生错误，调试信息：" + fmt.Sprintf("%v", err))
		}
	}

	err := SendGroupMessage(uint32((*msg)["sender"].(map[string]interface{})["group"].(map[string]interface{})["id"].(float64)), msgChain)
	if err != nil {
		return
	}

}

func receive(msg Message) error {
	if msg["type"] != "GroupMessage" {
		return nil
	}

	msgChain := msg["messageChain"].([]interface{})

	if len(msgChain) < 2 {
		return nil
	}

	mainMsg := msgChain[1].(map[string]interface{})

	if mainMsg["type"] != "Plain" {
		return nil
	}

	calls := strings.Fields(mainMsg["text"].(string))
	if len(calls) < 1 {
		return nil
	}

	for _, v := range actions {
		for _, v2 := range v["key"].([]string) {
			if calls[0] == v2 {
				calls2 := calls[1:]
				handle(&calls2, &msg, &v)
			}
		}
	}

	return nil
}

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

	log.Info().
		Msg("连接至Mirai成功")

	initApp()

	return nil
}
