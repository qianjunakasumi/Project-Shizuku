/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : uehara
*   File Name    : networkMiddleware.go
*   File Path    : internal/uehara/
*   Author       : Qianjunakasumi
*   Description  : 适用于 Mirai 的网络中间件
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   Variables:
*     json -- JSON 解析器
*
*   type Content map[string]interface{} -- Mirai API 返回的内容结构
*
*   func post(address string, content Content) (Content, error)  -- GET Mirai API Content
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

package uehara

import (
	"io/ioutil"

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/networkware"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Content Mirai 返回的内容
type Content map[string]interface{}

func post(address string, content Content) (Content, error) {

	bytes, err := json.Marshal(&content)
	if err != nil {
		return nil, err
	}

	req := new(networkware.Networkware)
	req.Address = "http://" + configs.Conf.MiraiAddress + "/" + address
	req.Body = bytes
	req.Method = "POST"
	req.Header = [][]string{
		{"Content-Type", "application/json; charset=utf-8"},
	}

	res, err := req.Send()
	if err != nil {
		return nil, err
	}
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content2 := make(Content)
	if err = json.Unmarshal(cont, &content2); err != nil {
		return nil, err
	}

	return content2, nil

}
