/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : message
*   File Name    : main.go
*   File Path    : internal/ayumu/message/
*   Author       : Qianjunakasumi
*   Description  : AYUMU 消息请求
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

package message

import "os"

// Message 发送的消息的结构
type Message struct {
	Target   uint64  // 目标
	IsCancel bool    // 取消发送
	Chain    []Chain // 消息链
}

type Chain struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

var wd, _ = os.Getwd()

// NewText 新建消息并新增文本
func NewText(s string) *Message {

	m := new(Message)
	m.AddText(s)

	return m

}

// NewImage 新建消息并新增图片
func NewImage(p string) *Message {

	m := new(Message)
	m.AddImage(p)

	return m

}

// NewVoice 新建消息并新增图片
func NewVoice(p string) *Message {

	m := new(Message)
	m.AddVoice(p)

	return m

}

// New 新建消息
func New() *Message { return new(Message) }

// AddText 新增文本
func (m *Message) AddText(s string) {

	m.Chain = append(m.Chain, Chain{
		"text",
		map[string]interface{}{
			"text": s,
		},
	})

}

// AddImage 新增图片
func (m *Message) AddImage(p string) {

	m.Chain = append(m.Chain, Chain{
		"image",
		map[string]interface{}{
			"file": "file:///" + wd + "/" + p,
		},
	})

}

// AddVoice 新增音频
func (m *Message) AddVoice(p string) {

	m.Chain = append(m.Chain, Chain{
		"record",
		map[string]interface{}{
			"file": "file:///" + wd + "/" + p,
		},
	})

}

// AddAt 新增提醒
func (m *Message) AddAt(u uint64) {

	m.Chain = append(m.Chain, Chain{
		"at",
		map[string]interface{}{
			"qq": u,
		},
	})
}
