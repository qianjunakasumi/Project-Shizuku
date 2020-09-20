/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : shizuku
*   File Name    : message.go
*   File Path    : internal/shizuku/
*   Author       : Qianjunakasumi
*   Description  : SHIZUKU 消息处理器
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

package shizuku

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

var sorry = "=== :( 失败 ===\n错误已被记录，感谢您的理解和支持。\n\n"

func (s SHIZUKU) monitor() {
	for {
		go s.process(<-s.msgChan)
	}
}

func (s *SHIZUKU) process(m *QQMsg) {

	var (
		r   *Message
		err error
		ok  bool
	)

	if s.block(m) {
		return
	}

	if j := s.Job[m.Group.ID]; j.Enable {

		r, err = j.Pointer.OnJobCall(m, s)
		if err != nil {
			r = NewText(sorry + "执行：事务\n信息：" + fmt.Sprintf("%v", err))
		}
		if r != nil {
			ok = true
		}

	}

	if !ok {

		r, err = s.processCommand(m)
		if err != nil {
			r = NewText(sorry + "执行：普通\n信息：" + fmt.Sprintf("%v", err))
		}
		if r == nil {
			return
		}

	}

	s.Rina.SendGroupMsg(r.To(m.Group.ID))

}

func (s SHIZUKU) block(m *QQMsg) bool {

	// 当长度小于1时消息无法获取
	if len(m.Chain) < 1 {
		return true
	}

	// 匿名用户禁止
	if m.User.ID == 80000000 {
		return true
	}

	// 史诗：黑名单 | 有需要的时候再写 // 咕咕咕
	// 实现：在启动的时候载入黑名单列表，若匹配成功则阻挡

	return false

}

func (s *SHIZUKU) processCommand(m *QQMsg) (*Message, error) {

	calls := strings.Fields(m.Chain[0].Text)
	if len(calls) < 1 {
		return nil, nil
	}

	for _, v := range s.command {
		for _, v2 := range v.Keys {

			if calls[0] != v2 {
				continue
			}

			c, err := s.callParser(calls[1:], v)
			if err != "" {
				return NewText(err), nil
			}

			m.Call = c
			log.Info().Interface("数据", m.Call).Msg("命令请求")

			r, err2 := v.Pointer.OnCall(m, s)
			if err2 != nil {
				return nil, err2
			}

			return r, nil

		}
	}

	return nil, nil

}

func (s SHIZUKU) callParser(call []string, i *AppInfo) (map[string]string, string) {

	f := make(map[string]string)

	for _, v := range i.Expand {
		f[v.Name] = v.Default
	}

	n := 0
	for k := range f {

		if n == len(call) {
			if i.Expand[n].Require {
				return f, i.Expand[n].DisplayName + " 是必填字段"
			}
			break
		}

		str := strings.Split(call[n], "：")
		if len(str) == 1 {
			if l := i.Expand[n].Limit; !s.isLimit(call[n], l) {
				return f, "您的输入不被允许，但您可尝试下列值：" + fmt.Sprintf("%v", l)
			}

			f[k] = call[n]
		} else {
			//TODO Map方式
		}

		n++
	}

	return f, ""

}

func (s SHIZUKU) isLimit(t string, l []string) bool {

	if len(l) == 0 { // 不做限制
		return true
	}

	for _, v := range l {
		if t != v {
			continue
		}

		return true
	}

	return false

}
