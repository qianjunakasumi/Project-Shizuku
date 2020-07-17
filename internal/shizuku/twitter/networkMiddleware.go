/*
networkMiddleware.go: 适用于Twitter的网络中间件
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

import (
	"io/ioutil"

	"github.com/qianjunakasumi/shizuku/pkg/networkware"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Content map[string]interface{}

func get(address string) (Content, error) {
	req := new(networkware.Networkware)
	req.Address = "https://api.twitter.com/" + address
	req.Method = "GET"
	req.Header = [][]string{
		{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4128.3 Safari/537.36"},
		{"authorization", "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"},
		{"x-guest-token", token},
	}
	req.Proxy = "http://127.0.0.1:10809"
	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	content := make(map[string]interface{})
	if err = json.Unmarshal(cont, &content); err != nil {
		return nil, err
	}
	return content, nil
}
