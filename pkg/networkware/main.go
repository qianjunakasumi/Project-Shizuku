/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : networkware
*   File Name    : main.go
*   File Path    : pkg/networkware/
*   Author       : Qianjunakasumi
*   Description  : 基本封装的网络件
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   type Networkware struct -- 存储请求信息和提供相应方法的容器

*   func (n Networkware) transport(tran **http.Transport) -- 设置代理
*   func (n Networkware) header(req **http.Request)       -- 写入请求头部信息
*   func (n Networkware) Send() (*http.Response, error)   -- 发送请求并返回相应
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

package networkware

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

// Networkware 网络件
type Networkware struct {
	Address string
	Body    []byte
	Method  string

	Header [][]string
	Proxy  string
}

func (n Networkware) transport(tran **http.Transport) {
	if n.Proxy != "" {
		proxy, _ := url.Parse(n.Proxy)
		*tran = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
}

func (n Networkware) header(req **http.Request) {
	for i := 0; i < len(n.Header); i++ {
		(*req).Header.Set(n.Header[i][0], n.Header[i][1])
	}
}

// Send 发送请求
func (n Networkware) Send() (*http.Response, error) {
	if n.Address == "" {
		return nil, errors.New("请求地址为空")
	}

	transport := &http.Transport{}
	n.transport(&transport)
	client := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest(n.Method, n.Address, bytes.NewBuffer(n.Body))
	if err != nil {
		return nil, err
	}
	n.header(&req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
