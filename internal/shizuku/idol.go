/***********************************************************************************************************************
***  P R O J E C T  --  S H I Z U K U                                                   Q I A N J U N A K A S U M I  ***
************************************************************************************************************************
* Basic:
*
*   Package Name : shizuku
*   File Name    : idol.go
*   File Path    : internal/shizuku/
*   Author       : Qianjunakasumi
*   Description  : SHIZUKU 应用定义的关键信息
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

import "math"

type Idol struct {
	ID       string   // 编号
	Name     string   // 名
	PickName string   // 角色名
	Key      []string // 关键字
	Twitter  twitter
}

type twitter struct {
	Followers string
	Tweets    string
	Push      func(x float64) float64
}

var Idols = []Idol{{
	"00",
	"ラブライブ！シリーズ公式",
	"ラブライブ！シリーズ公式",
	[]string{"ラブライブ", "lovelive", "官方", "LoveLive_staff"},
	twitter{
		"lovelive_staff",
		"347849994",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"01",
	"樱坂雫",
	"前田 佳織里",
	[]string{"樱坂", "雫", "前田", "佳織里", "佳织里", "kaor1n_n", "加智力"},
	twitter{
		"kaor1n_n",
		"880621944101404672",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"02",
	"上原步梦",
	"大西亜玖璃",
	[]string{"上原", "步梦", "大西亜", "大西亚", "玖璃", "aguri_onishi", "阿兔嘭"},
	twitter{
		"aguri_onishi",
		"991283114885365761",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"03",
	"中须霞",
	"相良 茉優",
	[]string{"相良", "茉優", "茉优", "MayuSgr", "麻油鸡"},
	twitter{
		"mayusgr",
		"1057615013282631680",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"04",
	"朝香果林",
	"久保田未夢",
	[]string{"久保田", "未夢", "未梦", "iRis_k_miyu"},
	twitter{
		"iris_k_miyu",
		"2384783184",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"05",
	"宫下爱",
	"村上奈津実",
	[]string{"村上", "奈津実", "奈津实", "natyaaaaaaan07"},
	twitter{
		"natyaaaaaaan07",
		"760000974005997568",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"06",
	"近江彼方",
	"鬼頭明里",
	[]string{"鬼頭", "鬼头", "明里", "kitoakari_1016"},
	twitter{
		"kitoakari_1016",
		"1141319903250534400",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"07",
	"优木雪菜",
	"楠木ともり",
	[]string{"楠木", "ともり", "灯", "tomori_kusunoki"},
	twitter{
		"tomori_kusunoki",
		"847365153691582465",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"08",
	"艾玛·维尔德",
	"指出 毬亜",
	[]string{"指出", "毬亜", "毬亚", "sashide_m"},
	twitter{
		"sashide_m",
		"1075210326990217216",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"09",
	"天王寺璃奈",
	"田中ちえ美",
	[]string{"田中", "ちえ美", "千惠美", "t_chiemi1006"},
	twitter{
		"t_chiemi1006",
		"1176845285059747842",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}, {
	"10",
	"三船栞子",
	"小泉萌香",
	[]string{"三船", "栞子", "小泉", "萌香", "k_moeka_", "萌p"},
	twitter{
		"k_moeka_",
		"4110103573",
		func(x float64) float64 {
			return -0.0047*math.Pow(x, 4) + 0.1544*math.Pow(x, 3) - 1.1701*math.Pow(x, 2) + 2.8274*x + 4.8613
		}},
}}

// TODO 干活 模糊匹配
func XXXXX() {

	for _, _ = range Idols {

	}

}
