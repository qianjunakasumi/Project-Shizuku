/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : ayumu
*   File Name    : main.go
*   File Path    : internal/ayumu/
*   Author       : Qianjunakasumi
*   Description  : AYUMU 基础结构
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

package ayumu

import (
	"golang.org/x/net/websocket"
)

type Ayumu struct {
	wsAddr string          // Websocket 地址
	ws     *websocket.Conn // Websocket 会话
	Exit   chan uint8      // 退出信号
}

// New 新建一个 AYUMU
func New(addr string) *Ayumu {

	return &Ayumu{
		wsAddr: addr,
	}

}
