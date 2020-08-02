/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : database
*   File Name    : main.go
*   File Path    : pkg/database/
*   Author       : Qianjunakasumi
*   Description  :
*
*----------------------------------------------------------------------------------------------------------------------*
* Summary:
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

package database

import (
	"database/sql"

	"github.com/qianjunakasumi/project-shizuku/configs"

	_ "github.com/go-sql-driver/mysql" // 连接数据库需要的包
	"github.com/rs/zerolog/log"
)

var DB *sql.DB

func Connect() error {

	var err error

	DB, err = sql.Open("mysql", configs.Conf.Databaseurl)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	return nil

}

func Close() {

	err := DB.Close()
	if err != nil {
		log.Error().Err(err)
		return
	}

}
