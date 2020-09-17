/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : shizuku
*   File Name    : info.go
*   File Path    : internal/app/shizuku/
*   Author       : Qianjunakasumi
*   Description  : 应用定义
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

package shizuku

import (
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"
)

type shizukuchan struct{}

type subApper interface {
	onSubCallByQQ(qp *shizuku.QQMsg, cl map[string]string) (*shizuku.Message, error)
}

var subList = map[string]subApper{
	"系统信息": new(sysInfo),
}

func init() {

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "shizuku",
		DisplayName: "小雫",
		Keys:        []string{"s", "小雫"},
		Expand: shizuku.Expand{
			{
				"func",
				"功能",
				[]string{"功能", "操作"},
				getKeys(),
				false,
				"",
			},
		},
		Pointer: new(shizukuchan),
	})

}

func getKeys() []string {

	var l []string
	for k := range subList {
		l = append(l, k)
	}

	return l

}

func (s shizukuchan) OnCall(qm *shizuku.QQMsg, _ *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	f := subList[qm.Call["func"]]
	if f == nil {
		return shizuku.NewText("欢迎使用小雫Project SHIZUKU"), nil
	}

	r, err := f.onSubCallByQQ(qm, qm.Call)
	if err != nil {
		return nil, err
	}

	return r, nil

}
