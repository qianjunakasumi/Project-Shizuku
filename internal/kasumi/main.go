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
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

/*
Request 请求字段

Addr 和 Method 是必填字段

当上述字段留空时返回空指针，请谨慎使用
*/
type Request struct {
	Host      string     // 主机
	Addr      string     // 地址
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

	return &Network{Request: r}

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
		res         *http.Response
	)

	// 避免 deadlock 产生 panic ，注意可能 nil
	defer func() {
		c <- res
	}()

	if err != nil {
		log.Error().Err(err).Msg("设置客户端出错")
		return
	}

	req, err := http.NewRequest(n.Method, n.Addr, nil)
	if err != nil {
		log.Error().Err(err).Msg("新建请求出错")
		return
	}

	n.setHeader(req)

	res, err = client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("返回错误")
		return
	}

}
