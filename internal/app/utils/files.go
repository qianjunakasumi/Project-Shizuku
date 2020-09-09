/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : utils
*   File Name    : files.go
*   File Path    : internal/app/utils/
*   Author       : Qianjunakasumi
*   Description  : 文件相关工具
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

package utils

import (
	"errors"
	"io/ioutil"
	"math/rand"
)

// GetFileNameByDir 随机获得目录下的一个文件的名称
func GetFileNameByDir(p string) (name string, err error) {

	files, err := ioutil.ReadDir(p)
	if err != nil {
		return
	}
	if len(files) < 1 {
		err = errors.New("文件夹为空")
		return
	}

	file := files[rand.Int31n(int32(len(files)))]
	name = p + file.Name()

	return

}
