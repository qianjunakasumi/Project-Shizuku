/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : llas
*   File Name    : info.go
*   File Path    : internal/app/llas/
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

package llas

import "github.com/qianjunakasumi/project-shizuku/internal/shizuku"

func init() {

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "RandomStill",
		DisplayName: "来一张场景",
		Keys:        []string{"来一张场景", "lyzcj", "still"},
		Expand: []shizuku.Expand{
			{
				"idol",
				"偶像",
				[]string{"偶像", "爱抖露"},
				[]string{},
				false,
				"_SHIZUKU默认检查专用",
			},
		},
		Pointer: new(randomStill).initApp(),
	})

}

func (r *randomStill) initApp() *randomStill {

	r.root = "assets/images/llas/"
	return r

}
