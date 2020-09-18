/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : shizuku
*   File Name    : main.go
*   File Path    : internal/shizuku/
*   Author       : Qianjunakasumi
*   Description  : SHIZUKU 关键功能
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
	"github.com/qianjunakasumi/project-shizuku/configs"
	"github.com/qianjunakasumi/project-shizuku/internal/utils/database"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type (
	// SHIZUKU SHIZUKU Robot
	SHIZUKU struct {
		QQID    uint64         // QQID QQ 号
		Rina    *Rina          // Rina QQ 客户端
		Job     map[uint64]job // Job 事务列表
		Conf    configs.App    // Conf 应用配置信息
		command []*AppInfo     // command 应用列表
		msgChan chan *QQMsg    // msgChan 消息管道
	}

	job struct {
		Enable  bool     // Enable 是否启用
		Pointer AppJober // Pointer 事务实例 指针
	}

	// Apper 应用接口
	Apper interface {
		OnCall(qm *QQMsg, sz *SHIZUKU) (rm *Message, err error)
	}

	// AppJober 应用事务接口
	AppJober interface {
		OnJobCall(qm *QQMsg, sz *SHIZUKU) (rm *Message, err error)
	}

	// AppTaskr 应用任务接口
	AppTaskr interface {
		OnTaskCall(sz *SHIZUKU) (rm *Message, err error)
	}

	// AppInfo 应用信息
	AppInfo struct {
		Name        string   // Name 应用名称
		DisplayName string   // DisplayName 应用显示名称
		Keys        []string // Keys 应用关键字
		Expand      Expand   // Expand 扩展
		Pointer     Apper    // Pointer 应用实例 指针
	}

	// Expand 扩展
	Expand []struct {
		Name        string   // 扩展名称
		DisplayName string   // 扩展显示名称
		Keys        []string // 扩展关键字
		Limit       []string // 扩展值限制
		Require     bool     // 是否必须
		Default     string   // 默认值
	}

	// AppTaskInfo 应用任务信息
	AppTaskInfo struct {
		Name    string   // 任务名称
		Spec    string   // 执行时间
		QQID    uint64   // 任务目标
		Pointer AppTaskr // 应用实例
	}

	task struct {
		Info *AppTaskInfo // 信息
	}
)

var (
	shizuku      *SHIZUKU       // 应用
	InitAppInfo  []*AppInfo     // 初始化时应用信息列表
	InitTaskInfo []*AppTaskInfo // 初始化时任务列表
)

// NewApp 新建应用
func NewApp(i *AppInfo) {

	InitAppInfo = append(InitAppInfo, i)
	log.Info().Str("命令", i.DisplayName).Msg("注册命令成功")

}

// NewApp 新建定时任务
func NewTask(i *AppTaskInfo) {

	InitTaskInfo = append(InitTaskInfo, i)
	log.Info().Str("任务", i.Name).Msg("注册定时任务成功")

}

// New 新建 SHIZUKU Robot
func New() {

	c := configs.GetAllConf()
	err := database.Connect(c.Databaseurl)
	if err != nil {
		log.Error().Err(err).Msg("无法连接至数据库")
	}

	var (
		ch = make(chan *QQMsg, 10)
		r  = newRina(c.QQID, c.QQPassword, &ch)
	)

	s := &SHIZUKU{
		QQID:    c.QQID,
		Rina:    r,
		Job:     make(map[uint64]job),
		Conf:    c.App,
		command: InitAppInfo,
		msgChan: ch,
	}

	cron2 := cron.New()
	for i := 0; i < len(InitTaskInfo); i++ {

		_, err := cron2.AddJob(InitTaskInfo[i].Spec, task{Info: InitTaskInfo[i]})
		if err != nil {
			log.Error().Err(err).Msg("注册任务失败")
		}

	}
	cron2.Start()

	go s.monitor()
	r.regEventHandle()

	shizuku = s

}

// Run 定时任务执行实现
func (t task) Run() {

	log.Info().Str("任务", t.Info.Name).Msg("执行定时任务")
	rm, err := t.Info.Pointer.OnTaskCall(shizuku)
	if err != nil {
		log.Error().Err(err).Msg("定时任务执行失败")
	}
	if rm == nil {
		return
	}

	shizuku.Rina.SendGroupMsg(rm.To(t.Info.QQID))

}

// OpenJob 开启一个事务
func (s *SHIZUKU) OpenJob(u uint64, p AppJober) { s.Job[u] = job{true, p} }

// CloseJob 关闭一个事务
func (s *SHIZUKU) CloseJob(u uint64) { delete(s.Job, u) }
