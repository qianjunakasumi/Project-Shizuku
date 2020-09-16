/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : basic
*   File Name    : main.go
*   File Path    : cmd/shizuku/basic/
*   Author       : Qianjunakasumi
*   Description  : 全局初始化
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

package basic

import (
	"fmt"
	"os"

	"github.com/qianjunakasumi/project-shizuku/configs"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {

	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "01-02|15:04:05"}
	writer.FormatCaller = func(i interface{}) string {
		return fmt.Sprintf("\x1b[1m%v\x1b[0m \x1b[36m>\x1b[0m", i.(string)[39:])
	}

	log.Logger = log.Output(writer).With().Caller().Logger()

	log.Info().Msg("Copyright (C) 2020-present  QianjuNakasumi  AGPL-3.0 License | version：" +
		configs.Version + "+" + configs.CommitId)

}
