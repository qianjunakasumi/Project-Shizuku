/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : kasumi
*   File Name    : wrapper.go
*   File Path    : internal/kasumi/
*   Author       : Qianjunakasumi
*   Description  : 网络请求件封装
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

package kasumi

import (
	"errors"
	"net/http"
)

// MiraiReq 适用于 MiraiReq 的网络请求
func (n *Network) MiraiReq() (C, error) {

	n.Addr = "http://" + n.Host + "/" + n.Addr
	n.Header = [][]string{
		{"Content-Type", "application/json; charset=utf-8"},
	}

	c := make(chan *http.Response)

	go n.send(c)

	res := <-c
	if res == nil {

		return nil, errors.New("网络请求失败")

	}

	defer res.Body.Close()

	content, err := n.jsonToMap(res)
	if err != nil {

		return nil, err

	}

	return content, nil

}

// TwitterReq 适用于 Twitter 的网络请求
func (n Network) TwitterReq() {

}
