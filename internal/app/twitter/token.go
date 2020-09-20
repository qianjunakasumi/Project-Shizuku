/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : token.go
*   File Path    : internal/app/twitter/
*   Author       : Qianjunakasumi
*   Description  : Token 相关
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

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/kasumi"

	"github.com/rs/zerolog/log"
)

var token string

type (
	extractTokener interface {
		extractToken(res *http.Response)
	}

	extractToken struct {
		next extractTokener
	}

	extractToken2 struct{}
)

func (e extractToken) extractToken(res *http.Response) {

	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("提取 Token 时发生错误")
		return
	}

	p := bytes.LastIndex(cont, []byte("document.cookie = decodeURIComponent"))
	if p < 0 {
		e.next.extractToken(res)
		return
	}

	token := string(cont)[p+41 : p+60]
	if !isRealToken(token) {
		e.next.extractToken(res)
		return
	}

	writeToken(token)

}

func (e extractToken2) extractToken(res *http.Response) {

	var (
		tokenContent = res.Header["Set-Cookie"][0]
		p            = bytes.IndexAny([]byte(tokenContent), "gt=") + 3
	)

	if p == -1 {
		log.Error().Msg("定位 Token 时发生错误")
		return
	}

	token := tokenContent[p : p+19]
	if isRealToken(token) {
		writeToken(token)
		return
	}

	log.Error().Msg("无法获取 Token")

}

// Timer 定时更新 Token
func Timer() {

	FetchToken()
	go func() {

		for range time.Tick(time.Hour) {
			FetchToken()
		}

	}()

}

// FetchToken 获取 Token
func FetchToken() {

	res := kasumi.New(&kasumi.Request{
		Addr:   "twitter.com/kaor1n_n",
		Method: "GET",
		Header: [][]string{},
	}).TwitterReq(configs.GetProxyAddr())
	if res == nil {
		log.Error().Msg("请求数据出错")
		return
	}

	var (
		ext  = new(extractToken)
		ext2 = new(extractToken2)
	)

	ext.next = ext2
	ext.extractToken(res)

}

func isRealToken(t string) bool {

	_, err := strconv.ParseUint(t, 10, 64)
	if err != nil {
		log.Error().Str("Token", t).Msg("校验 Token 结果失败")
		return false
	}

	return true

}

func writeToken(t string) {

	token = t
	log.Info().Str("Token", t).Msg("成功获取 Token")

}
