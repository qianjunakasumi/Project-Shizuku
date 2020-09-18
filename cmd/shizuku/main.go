/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : main
*   File Name    : main.go
*   File Path    : cmd/shizuku/
*   Author       : Qianjunakasumi
*   Description  : 程序入点
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
	_ "github.com/qianjunakasumi/project-shizuku/cmd/shizuku/basic" // 全局初始化
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"

	_ "github.com/qianjunakasumi/project-shizuku/internal/app/debug"     // 调试 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/guesssong" // 阅词识曲 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/llas"      // 来一张场景 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/meme"      // 表情 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/menu"      // 菜单 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/shizuku"   // 小雫 应用
	_ "github.com/qianjunakasumi/project-shizuku/internal/app/twitter"   // Twitter 应用
)

func main() {

	shizuku.New()

	// 调试程序使用
	/*go func() {
		log.Error().Err(http.ListenAndServe(":520", nil)).Msg("初始化错误")
	}()*/

	select {}

}
