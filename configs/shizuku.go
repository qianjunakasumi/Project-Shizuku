/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : configs
*   File Name    : shizuku.go
*   File Path    : configs/
*   Author       : Qianjunakasumi
*   Description  : 解析 SHIZUKU 配置参数
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   待完善
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

package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// conf configs.yml结构
type conf struct {
	Development      bool   // 开发模式
	QQNumber         uint32 // Robot QQ号
	MiraiAddress     string // Mirai API HTTP URL地址
	MiraiAuthKey     string // Mirai API HTTP AuthKey
	Databaseurl      string // 数据库地址
	TranslationAppID string // 百度翻译 APP ID
	TranslationKey   string // 百度翻译 Key
}

var (
	Conf      conf           // 配置文件
	Version   = "0.1.0-beta" // 版本
	BuildTime string         // 编译时的日期和时间
	CommitId  string         // 存储库最新提交的短SHA1
)

// SetConfigs 配置配置参数
func SetConfigs() error {
	file, err := ioutil.ReadFile("configs/configs.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return err
	}

	return nil
}
