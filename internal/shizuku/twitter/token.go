/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : token.go
*   File Path    : internal/shizuku/twitter/
*   Author       : Qianjunakasumi
*   Description  : 获取 Twitter Guest Token
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     token -- Twitter Guest Token
*
*   func isRealToken(t string) bool -- 校验 Token 是否有效
*   func writeToken(t string)       -- 写入 Token 串
*
*   type extractTokener interface                               -- 提取 Token 计划的接口
*     type extractToken struct                                  -- 提取 Token Plan.Ⅰ 的对象
*       func (e extractToken2) extractToken(res *http.Response) -- 提取 Token Plan.Ⅰ 的方法
*     type extractToken2 struct                                 -- 提取 Token Plan.Ⅱ 的对象
*       func (e extractToken2) extractToken(res *http.Response) -- 提取 Token Plan.Ⅱ 的方法
*
*   func FetchToken() -- 获取含 Token 的原始内容
*   func Timer()      -- 定时获取 Token
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

package twitter

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/qianjunakasumi/project-shizuku/internal/utils/networkware"

	"github.com/rs/zerolog/log"
)

var token string

func isRealToken(t string) bool {

	_, err := strconv.ParseUint(t, 10, 64)
	if err != nil {

		log.Error().
			Str("Token", t).
			Msg("校验 Token 结果失败")

		return false

	}

	return true

}

func writeToken(t string) {

	token = t

	log.Info().
		Str("Token", t).
		Msg("成功获取 Token")

}

type extractTokener interface {
	extractToken(res *http.Response)
}

type extractToken struct {
	next extractTokener
}

func (e extractToken) extractToken(res *http.Response) {

	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {

		log.Error().
			Err(err).
			Msg("提取 Token 时发生错误")

		return

	}

	var (
		p     = bytes.LastIndex(cont, []byte("document.cookie = decodeURIComponent"))
		token = string(cont)[p+41 : p+60]
	)

	if isRealToken(token) {

		writeToken(token)

		return

	}

	e.next.extractToken(res)

}

type extractToken2 struct {
	next extractTokener
}

func (e extractToken2) extractToken(res *http.Response) {

	var (
		tokenContent = res.Header["Set-Cookie"][0]
		p            = bytes.IndexAny([]byte(tokenContent), "gt=") + 3
	)

	if p == -1 {

		log.Error().
			Msg("定位 Token 时发生错误")

		return

	}

	token := tokenContent[p : p+19]

	if isRealToken(token) {

		writeToken(token)

		return

	}

	log.Error().
		Str("Token", token).
		Msg("提取 Token 时发生错误")

}

// FetchToken 获取token
func FetchToken() {

	req := new(networkware.Networkware)

	req.Address = "https://twitter.com/kaor1n_n"
	req.Method = "GET"
	req.Header = [][]string{{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4128.3 Safari/537.36"}}
	req.Proxy = "http://127.0.0.1:10809"

	res, err := req.Send()
	if err != nil {
		log.Error().Err(err)
		return
	}

	var (
		ext  = new(extractToken)
		ext2 = new(extractToken2)
	)

	ext.next = ext2

	ext.extractToken(res)

}

// Timer 定时器
func Timer() {

	FetchToken()

	go func() {

		for range time.Tick(time.Hour) {

			FetchToken()

		}

	}()

}
