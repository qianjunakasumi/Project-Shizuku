/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : kasumi
*   File Name    : expand.go
*   File Path    : internal/kasumi/
*   Author       : Qianjunakasumi
*   Description  : 网络请求件扩展
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

package kasumi

import (
	"io/ioutil"
	"net/http"

	"github.com/qianjunakasumi/project-shizuku/internal/utils/json"
)

func (n Network) jsonToMap(r *http.Response) (C, error) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {

		return nil, err

	}

	content := make(C)

	err = json.JSON.Unmarshal(b, &content)
	if err != nil {

		return nil, err

	}

	return content, nil

}
