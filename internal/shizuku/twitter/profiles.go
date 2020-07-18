/*
profiles.go: 档案
Copyright (C) 2020-present  QianjuNakasumi

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package twitter

// Profiles 档案
type Profiles struct {
	name      string
	key       []string
	followers string
	tweets    string
}

var profiles = []Profiles{{
	"ラブライブ！シリーズ公式",
	[]string{"ラブライブ", "lovelive", "官方", "LoveLive_staff"},
	"lovelive_staff",
	"347849994",
}, {
	"前田 佳織里",
	[]string{"前田佳織里", "前田佳织里", "前田", "佳織里", "佳织里", "kaor1n_n", "加智力"},
	"kaor1n_n",
	"880621944101404672",
}, {
	"大西亜玖璃",
	[]string{"大西亜玖璃", "大西亚玖璃", "大西", "亜玖璃", "亚玖璃", "aguri_onishi", "阿兔嘭"},
	"aguri_onishi",
	"991283114885365761",
}, {
	"相良 茉優",
	[]string{"相良茉優", "相良茉优", "相良", "茉优", "MayuSgr", "麻油鸡"},
	"mayusgr",
	"1057615013282631680",
}, {
	"久保田未夢",
	[]string{"久保田未夢", "久保田未梦", "久保田", "未夢", "未梦", "iRis_k_miyu"},
	"iris_k_miyu",
	"2384783184",
}, {
	"村上奈津実",
	[]string{"村上奈津実", "村上奈津实", "村上", "奈津実", "奈津实", "natyaaaaaaan07"},
	"natyaaaaaaan07",
	"760000974005997568",
}, {
	"鬼頭明里",
	[]string{"鬼頭明里", "鬼头明里", "鬼頭", "鬼头", "明里", "kitoakari_1016"},
	"kitoakari_1016",
	"1141319903250534400",
}, {
	"楠木ともり",
	[]string{"楠木ともり", "楠木灯", "楠木", "ともり", "灯", "tomori_kusunoki"},
	"tomori_kusunoki",
	"847365153691582465",
}, {
	"指出 毬亜",
	[]string{"指出毬亜", "指出毬亚", "指出", "毬亜", "毬亚", "sashide_m"},
	"sashide_m",
	"1075210326990217216",
}, {
	"田中ちえ美",
	[]string{"田中ちえ美", "田中千惠美", "田中", "ちえ美", "千惠美", "t_chiemi1006"},
	"t_chiemi1006",
	"1176845285059747842",
}}

// GetKeys 获取所有关键字
func GetKeys() []string {
	var keys []string

	for _, v := range profiles {
		keys = append(keys, v.key...)
	}

	return keys
}

func getProfile(str string) Profiles {
	for _, v := range profiles {
		for _, v2 := range v.key {
			if v2 == str {
				return v
			}
		}
	}

	var i Profiles
	return i
}
