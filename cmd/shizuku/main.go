/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : main
*   File Name    : main.go
*   File Path    : cmd/shizuku/
*   Author       : Qianjunakasumi
*   Description  : 程序主入口，启动 SHIZUKU 应用
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
*   func main() -- 设置全局 LOG 样式，加载配置文件，启动 SHIZUKU 应用
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

package main

import (
	"fmt"
	"os"

	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "01-02|15:04:05"})

	err := configs.SetConfigs()
	if err != nil {
		log.Fatal().Msg("加载配置文件时出现错误")
		return
	}

	log.Info().Msg("Copyright (C) 2020-present  QianjuNakasumi  AGPL-3.0 License | Release version：" +
		configs.Version + "+" +
		configs.CommitId + "." +
		configs.BuildTime)

	if err := shizuku.Start(); err != nil {
		log.Fatal().Msg(fmt.Sprintf("%v", err))
	}
	log.Info().
		Msg("启动SHIZUKU完毕")

	select {}
}
