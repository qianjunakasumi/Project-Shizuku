/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : guesssong
*   File Name    : main.go
*   File Path    : internal/app/guesssong/
*   Author       : Qianjunakasumi
*   Description  : 阅词识曲
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

import (
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/qianjunakasumi/project-shizuku/internal/app/utils"
	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type guesssong struct {
	root     string
	gameData map[uint64]*game
}

type song struct {
	Name     string   `yaml:"name"`     // 歌曲名
	Singer   string   `yaml:"singer"`   // 五阿泰
	Key      []string `yaml:"key"`      // 关键字
	CoverImg string   `yaml:"coverimg"` // 封面
	Lyrics   []struct {
		Point string `yaml:"point"` // 时间定位
		Text  string `yaml:"text"`  // 歌词文本
	} `yaml:"lyrics"` // 歌词
}

type game struct {
	times uint8 // 游玩次数
	point int32 // 歌词定位
	song  song  // 歌曲信息
}

func (g *guesssong) OnCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	if g.gameData[qm.Group.ID] != nil {
		rm = shizuku.NewText("正在游戏中哦，您可以发送“不玩了”取消游戏")
		return
	}

	songPath, err := utils.GetFileNameByDir(g.root)
	if err != nil {
		return
	}

	songFile, err := ioutil.ReadFile(songPath)
	if err != nil {
		return
	}

	var song song
	err = yaml.Unmarshal(songFile, &song)
	if err != nil {
		return
	}

	n := rand.Int31n(int32(len(song.Lyrics)))
	rm = shizuku.NewText("Hi~ o(*￣▽￣*)ブ欢迎游玩！歌词找好啦「" + song.Lyrics[n].Text + "」快来猜猜看是哪首歌吧~（长按头像@我并带上答案就可以参与游戏哦)")

	g.gameData[qm.Group.ID] = &game{0, n, song}
	sz.OpenJob(qm.Group.ID, g)
	log.Info().Strs("答案", song.Key).Msg("阅词识曲")

	return

}

func (g *guesssong) OnJobCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	c := qm.Chain
	if c[0].QQ != sz.QQID || len(c) < 2 {
		return
	}
	if c[1].Type != "text" {
		return
	}

	var (
		t    = strings.ToLower(strings.Replace(c[1].Text, " ", "", 1))
		song = g.gameData[qm.Group.ID]
	)

	if song.times > 20 {
		sz.CloseJob(qm.Group.ID)
		rm = shizuku.NewText("啊这，猜的次数太多啦。辛苦了，请休息一下再来游玩趴~")
	}

	switch t {
	case "提示":
		rm = shizuku.NewText("给你一点小提示吧，这首歌是" + song.song.Singer + "唱的哦") // TODO 更牛逼的做法
		return                                                           // 建议该模块使用抽象工厂实现
	}

	for _, v := range song.song.Key {

		v = strings.ToLower(v)
		if v != t {
			continue
		}

		sz.CloseJob(qm.Group.ID)
		delete(g.gameData, qm.Group.ID)
		rm = shizuku.NewText("恭喜您猜对啦~ o(*￣▽￣*)ブ，这首是「" + song.song.Name + "」，").
			AddText("歌词出处在「" + song.song.Lyrics[song.point].Point + "」哦\n======\n")

		for i := -2; i < 3; i++ {

			p := int(song.point) + i
			if p < 0 {
				p = 0
			} else if p >= len(song.song.Lyrics) {
				break
			}

			rm.AddText(song.song.Lyrics[p].Text + "\n")

		}

		rm.AddText("======\n专辑封面：").AddImage("assets/images/cd/" + song.song.CoverImg)

		return

	}

	song.times++
	rm = shizuku.NewText("啊哦，您回答的不是歌词「" + song.song.Lyrics[song.point].Text + "」的答案哦，请再接再厉！")

	return

}
