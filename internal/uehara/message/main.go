/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : message
*   File Name    : main.go
*   File Path    : internal/uehara/messageChain/
*   Author       : Qianjunakasumi
*   Description  : 适用于 Mirai 的消息链
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   type Chain struct -- 存储消息信息和提供相关方法的容器
*
*   func (m *Chain) Cancel()              -- 取消发送
*   func (m *Chain) AddText(str string)   -- 添加文本
*   func (m *Chain) AddAt(target uint32)  -- 添加提醒
*   func (m *Chain) AddImage(path string) -- 添加图片
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

// Chain 消息链
type Chain struct {
	Content  []map[string]interface{}
	IsCancel bool
}

type MessageInfo struct {
	UserName  string
	UserId    uint32
	GroupName string
	GroupId   uint32
}

func (m *Chain) Cancel() {
	m.IsCancel = true
}

// AddText 插入文本
func (m *Chain) AddText(str string) {
	m.Content = append(m.Content, map[string]interface{}{
		"type": "Plain",
		"text": str,
	})
}

// AddAt 插入提醒
func (m *Chain) AddAt(target uint32) {
	m.Content = append(m.Content, map[string]interface{}{
		"type":    "At",
		"target":  target,
		"display": "@",
	})
}

// AddImage 插入图片
func (m *Chain) AddImage(path string) {
	m.Content = append(m.Content, map[string]interface{}{
		"type": "Image",
		"path": "../../../../../" + path,
	})
}
