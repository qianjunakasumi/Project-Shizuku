/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : configs
*   File Name    : main.go
*   File Path    : configs/
*   Author       : Qianjunakasumi
*   Description  : 读取配置信息
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

var (
	Version  = "2.0.0" // Version 版本号
	CommitId string    // CommitId 提交的短ID
	confs    *Conf     // confs 配置信息
)

type (
	// Conf 配置
	Conf struct {
		QQID        uint64 `yaml:"qqID"`        // QQID QQ帐号
		QQPassword  string `yaml:"qqPassword"`  // QQPassword QQ密码
		Databaseurl string `yaml:"databaseURL"` // Databaseurl 数据库地址
		App         App    `yaml:"app"`         // App 应用配置
	}

	// App 应用配置
	App struct {
		ProxyAddr        string `yaml:"proxyAddr"`        // ProxyAddr 代理地址
		TranslationAppID string `yaml:"translationAppid"` // TranslationAppID 百度翻译 APP ID
		TranslationKey   string `yaml:"translationKey"`   // TranslationKey 百度翻译 Key
	}
)

// ReadConfigs 读取配置文件
func ReadConfigs() (err error) {

	f, err := ioutil.ReadFile("configs/configs.yml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(f, &confs)
	if err != nil {
		return
	}

	return

}

// GetAllConf 获取所有配置
func GetAllConf() *Conf { return confs }

// GetProxyAddr 获取代理配置
func GetProxyAddr() string { return confs.App.ProxyAddr }
