/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : guesssong
*   File Name    : info.go
*   File Path    : internal/app/guesssong/
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

package guesssong

import "github.com/qianjunakasumi/project-shizuku/internal/shizuku"

func init() {

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "guessSong",
		DisplayName: "阅词识曲",
		Keys:        []string{"阅词识曲", "阅词识歌", "阅词猜曲", "看词猜歌", "看词识歌", "看词猜歌", "ycsq"},
		Expand:      []shizuku.Expand{},
		Pointer:     new(guesssong).initApp(),
	})

}

func (g *guesssong) initApp() *guesssong {

	g.root = "assets/game/guesssong/"
	g.gameData = map[uint64]*game{}

	return g

}
