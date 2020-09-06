/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : shizuku
*   File Name    : qq.go
*   File Path    : internal/shizuku/
*   Author       : Qianjunakasumi
*   Description  : Rina QQ 协议功能
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
	"io/ioutil"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	m2 "github.com/Mrs4s/MiraiGo/message"
	"github.com/rs/zerolog/log"
)

// Rina Rina // TODO 更多
type Rina struct {
	c       *client.QQClient // 客户端
	msgChan *chan *QQMsg     // 消息管道
}

// QQMsg 接收的 QQ 消息
type QQMsg struct {
	Type  *Idol
	Chain []Chain // 消息链
	Call  map[string]string
	Group struct {
		ID   uint64 // 群号
		Name string // 群名
	}
	User struct {
		ID   uint64 // QQ号
		Name string // QQ名
	}
}

// Chain 消息链
type Chain struct {
	Type string // 类型：text、image、at
	Text string // text
	URL  string // image
	QQ   uint64 // at
}

// Message 返回的 QQ 消息
type Message struct {
	target uint64            // 目标
	chain  m2.SendingMessage // 消息链
}

// newRina 新增 Rina
func newRina(i uint64, p string, ch *chan *QQMsg) (r *Rina) {

	err := client.SystemDeviceInfo.ReadJson([]byte("{\"display\":\"MIRAI.666470.001\",\"product\":\"mirai\",\"device\":\"mirai\",\"board\":\"mirai\",\"model\":\"mirai\",\"finger_print\":\"mamoe/mirai/mirai:10/MIRAI.200122.001/5696651:user/release-keys\",\"boot_id\":\"58fe8ac7-4de7-71ec-073d-07eb3187a533\",\"proc_version\":\"Linux version 3.0.31-HxHC3WtY (android-build@xxx.xxx.xxx.xxx.com)\",\"imei\":\"351912693210254\"}"))
	if err != nil {
		log.Panic().Err(err).Msg("设置设备信息失败")
	}

	c := client.NewClient(int64(i), p)
	c.OnLog(func(q *client.QQClient, e *client.LogEvent) {
		switch e.Type {

		case "INFO":
			log.Info().Str("信息", e.Message).Msg("协议")

		case "ERROR":
			log.Error().Str("信息", e.Message).Msg("协议")
		}
	})

	r = &Rina{c: c, msgChan: ch}
	if err := r.login(); err != nil {
		log.Panic().Msg("登录失败")
	}

	return

}

func (r Rina) login() (err error) {

	re, err := r.c.Login()
	if err != nil {
		return
	}
	if !re.Success {
		return
	}

	log.Info().Msg("登录成功：" + r.c.Nickname)

	err = r.c.ReloadGroupList()
	if err != nil {
		log.Error().Err(err).Msg("加载群列表失败")
		return
	}

	err = r.c.ReloadFriendList()
	if err != nil {
		log.Error().Err(err).Msg("加载好友列表失败")
		return
	}

	log.Info().Int("个数", len(r.c.GroupList)).Msg("加载群列表成功")
	log.Info().Int("个数", len(r.c.FriendList)).Msg("加载好友列表成功")

	return

}

func (r Rina) regEventHandle() {

	r.c.OnGroupMessage(r.onGroupMsg)

	// 断线重连
	r.c.OnDisconnected(func(q *client.QQClient, e *client.ClientDisconnectedEvent) {
		for {

			log.Warn().Msg("啊哦连接丢失了，准备重连中...1s")
			time.Sleep(time.Second)
			if err := r.login(); err != nil {
				log.Warn().Msg("重登录失败，再次尝试中...")
				continue
			}

			return

		}
	})

}

func (r Rina) onGroupMsg(_ *client.QQClient, m *m2.GroupMessage) {

	msg := &QQMsg{
		Type:  FuzzyGetIdol(m.GroupName),
		Chain: []Chain{},
		Group: struct {
			ID   uint64
			Name string
		}{
			uint64(m.GroupCode),
			m.GroupName,
		},
		User: struct {
			ID   uint64
			Name string
		}{
			uint64(m.Sender.Uin),
			m.Sender.Nickname,
		},
	}

	for _, v := range m.Elements {
		switch e := v.(type) {
		case *m2.TextElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "text",
				Text: e.Content,
			})

		case *m2.AtElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "at",
				QQ:   uint64(e.Target),
			})

		case *m2.ImageElement:
			msg.Chain = append(msg.Chain, Chain{
				Type: "image",
				URL:  e.Url,
			})

		}
	}

	log.Info().
		Interface("群类型", msg.Type.PickName).
		Str("群名", msg.Group.Name).
		Str("昵称", msg.User.Name).
		Interface("原文", msg.Chain).
		Msg("收到群消息")

	*r.msgChan <- msg
}

// NewMsg 新建消息结构体
func NewMsg() *Message { return &Message{chain: m2.SendingMessage{}} }

// NewText 新建文本消息结构体
func NewText(t string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddText(t) }

// NewImage 新建图片消息结构体
func NewImage(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddImage(p) }

// NewAudio 新建音频消息结构体
func NewAudio(p string) *Message { m := &Message{chain: m2.SendingMessage{}}; return m.AddAudio(p) }

// AddText 添加文本
func (m *Message) AddText(t string) *Message { m.chain.Append(m2.NewText(t)); return m }

// AddImage 添加图片
func (m *Message) AddImage(p string) *Message {

	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error().Err(err).Msg("读取图片失败")
		return m
	}

	m.chain.Append(m2.NewImage(b))

	return m

}

// AddAudio 添加音频
func (m *Message) AddAudio(p string) *Message {

	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error().Err(err).Msg("读取语音失败")
		return m
	}

	m.chain.Append(&m2.VoiceElement{Data: b})

	return m

}

// To 发送的目标
func (m *Message) To(i uint64) *Message { m.target = i; return m }

// SendGroupMsg 发送群消息
func (r Rina) SendGroupMsg(m *Message) {

	for k, v := range m.chain.Elements {
		if nm, ok := v.(*m2.ImageElement); ok {
			am, err := r.c.UploadGroupImage(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传图片失败")
			} else {
				m.chain.Elements[k] = am
			}

		}

		if nm, ok := v.(*m2.VoiceElement); ok {
			am, err := r.c.UploadGroupPtt(int64(m.target), nm.Data)
			if err != nil {
				log.Error().Err(err).Msg("上传语音失败")
			} else {
				m.chain.Elements[k] = am
			}
		}
	}

	log.Info().Msg("发送群消息")

	r.c.SendGroupMessage(int64(m.target), &m.chain)

}
