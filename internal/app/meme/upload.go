/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : meme
*   File Name    : upload.go
*   File Path    : internal/app/meme/
*   Author       : Qianjunakasumi
*   Description  : 上传表情
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

package meme

import (
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"
)

type uploadMeme struct {
	//root string
}

func (r uploadMeme) OnCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	// 字段是必须的

	if c := qm.Call["idol"]; c != "_SHIZUKU默认检查专用" {
		qm.Type = sz.FuzzyGetIdol(c)
	}

	rm = shizuku.NewText("欢迎使用表情上传，请发送您将上传的图片\n\n您上传的图片由 Project-SHIZUKU 服务器保存在中国香港。上传表情的同时会收集您的帐号信息：您的QQ号和当前群号。如您不同意上述方法，请勿使用本功能，使用“取消”退出")

	return

}

func (r uploadMeme) OnJobCall(_ *shizuku.QQMsg, _ *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	return
}

/*
var (
	UploadMeme = uploadMeme
	Uploading  = uploading
	Agree      = agree
)

func checkAgreeStatus(id uint64) (bool, error) {

	rows, err := database.DB.Query(`SELECT * FROM	meme_upload_agreeer WHERE meme_upload_agreeer.qqID = ?`, id)
	if err != nil {

		return false, err

	}
	defer rows.Close()

	if !rows.Next() {

		return false, nil

	}

	return true, nil

}

func agreeTerm(id uint64, name string, id2 uint64) error {

	insert, err := database.DB.Prepare(`INSERT INTO meme_upload_agreeer ( qqID, agreeInfo ) VALUES ( ?, ? )`)
	if err != nil {

		return err

	}
	defer insert.Close()

	_, err = insert.Exec(
		id,
		`{"groupName": "`+
			name+`", "groupID": `+
			strconv.FormatUint(id2, 10)+`, "time": "`+
			time.Now().Format(time.RFC3339)+`"}`,
	)
	if err != nil {

		return err

	}

	return nil

}

func agree(call string, info *message.Info) (*message.Chain, error) {

	call = strings.ReplaceAll(call, " ", "")

	m := new(message.Chain)

	if call != "我同意" {

		m.Cancel()
		return m, nil

	}

	err := agreeTerm(info.UserId, info.GroupName, info.GroupId)
	if err != nil {

		return m, nil

	}

	configs.QuoteJob[info.UserId] = nil
	configs.ImageJob[info.UserId] = &Uploading

	m.AddText("您已同意协议。请您发送欲上传的图片")

	return m, nil

}

func uploading(call string, info *message.Info) (*message.Chain, error) {

	m := new(message.Chain)

	m.AddText("您的图片已经上传。")

	m.AddText("\n图片的地址为：" + call)

	configs.ImageJob[info.UserId] = nil

	// 下载图片到内存中

	// 计算图片的SHA1值

	// 读取数据库检测是否有匹配的sha1

	// 若有则返回，图片已存在

	// 若无则将图片保存至目录下，以sha1值命名

	// 向数据库写入事件

	return m, nil

}

func uploadMeme(calls map[string]string, info *message.Info) (*message.Chain, error) {

	var (
		m            = new(message.Chain)
		isAgree, err = checkAgreeStatus(info.UserId)
	)

	if info.UserId > 2200000000 {

		m.AddText("您未在灰度测试范围中，请您耐心等待下一版本更新，谢谢。")

		return m, nil

	}

	if err != nil {

		return m, err

	}

	if isAgree {

		configs.ImageJob[info.UserId] = &Uploading

		m.AddText("请您发送欲上传的图片")

		return m, nil

	}

	m.AddAt(info.UserId)
	m.AddText("请您阅读以下协议，同意后方可开始上传\n")
	m.AddImage("assets/terms/uploadMeme.jpg")
	m.AddText("\n\n同意协议请回复本条消息：“我同意”")

	configs.QuoteJob[info.UserId] = &Agree

	return m, nil

}
*/
