/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : kasumi
*   File Name    : main.go
*   File Path    : internal/kasumi/
*   Author       : Qianjunakasumi
*   Description  : 网络请求件
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
	"bytes"
	"net/http"
	"net/url"

	"github.com/qianjunakasumi/project-shizuku/internal/utils/json"

	"github.com/qianjunakasumi/project-shizuku/configs"
)

// C 内容
type C map[string]interface{}

/*
Request 请求字段

Addr 和 Method 是必填字段

当上述字段留空时返回空指针，请谨慎使用
*/
type Request struct {
	Addr      string     // 地址
	Body      C          // 内容
	Method    string     // 模式
	Header    [][]string // 请求头
	ProxyAddr string     // 代理地址
}

// Network 网络请求件
type Network struct {
	*Request // 请求字段
}

// New 新建网络请求件
func New(r *Request) *Network {

	if r.Addr == "" || r.Method == "" {

		return nil

	}

	return &Network{
		Request: r,
	}

}

func (n Network) setClient() (*http.Client, error) {

	var tran *http.Transport

	if a := n.ProxyAddr; a == "" {

		tran = new(http.Transport)

	} else { // 使用代理地址请求

		u, err := url.Parse(a)
		if err != nil {

			return nil, err

		}

		tran = &http.Transport{
			Proxy: http.ProxyURL(u),
		}

	}

	return &http.Client{
		Transport: tran,
	}, nil

}

func (n Network) setHeader(r *http.Request) {

	for i := 0; i < len(n.Header); i++ {

		r.Header.Set(n.Header[i][0], n.Header[i][1])

	}

}

func (n Network) send(c chan *http.Response) {

	var (
		client, err = n.setClient()
		b           []byte
	)

	if err != nil {

		// hanld
		return

	}

	if n.Body != nil {

		b, err = json.JSON.Marshal(&n.Body)
		if err != nil {

			return

		}

	}

	req, err := http.NewRequest(n.Method, n.Addr, bytes.NewBuffer(b))
	if err != nil {

		return

	}

	n.setHeader(req)

	res, err := client.Do(req)
	if err != nil {

		return

	}

	c <- res

}

// MiraiReq 适用于 MiraiReq 的网络请求
func (n *Network) MiraiReq() (C, error) {

	n.Addr = "http://" + configs.Conf.MiraiAddress + "/" + n.Addr
	n.Header = [][]string{
		{"Content-Type", "application/json; charset=utf-8"},
	}

	c := make(chan *http.Response)

	go n.send(c)

	res := <-c

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
