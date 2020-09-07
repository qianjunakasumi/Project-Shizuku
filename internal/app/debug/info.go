/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : debug
*   File Name    : info.go
*   File Path    : internal/app/debug/
*   Author       : Qianjunakasumi
*   Description  : 应用信息注册
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

package debug

import "github.com/qianjunakasumi/project-shizuku/internal/shizuku"

type debug struct{}

func init() {

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "debug",
		DisplayName: "调试",
		Keys:        []string{"debug"},
		Expand: []shizuku.Expand{
			{
				"func",
				"功能",
				[]string{"功能", "操作"},
				[]string{},
				true,
				"",
			}, {
				"data",
				"数据",
				[]string{"数据"},
				[]string{},
				false,
				"",
			},
		},
		Pointer: new(debug),
	})

}

func (d debug) OnCall(qm *shizuku.QQMsg, _ *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	switch qm.Call["func"] {
	default:
		rm = shizuku.NewText("请输入正确的调试内容")

	}

	return

}
