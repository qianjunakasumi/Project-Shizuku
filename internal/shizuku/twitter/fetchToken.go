/*
fetchToken.go: 获取token
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
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/qianjunakasumi/shizuku/pkg/networkware"

	"github.com/rs/zerolog/log"
)

var token string

func checkToken(token string) bool {
	_, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		log.Warn().
			Str("功能", "twitter").
			Str("token", token).
			Msg("校验 token 结果失败")
		return false
	}

	return true
}

func fetchToken(res *http.Response) (string, error) {
	cont, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	p := bytes.LastIndex(cont, []byte("document.cookie = decodeURIComponent"))
	token := string(cont)[p+41 : p+60]

	if checkToken(token) {
		return token, nil
	}

	return "", errors.New("否")
}

func fetchToken2(res *http.Response) (string, error) {
	tokenContent := res.Header["Set-Cookie"][0]
	p := bytes.IndexAny([]byte(tokenContent), "gt=") + 3
	if p == -1 {
		return "", errors.New("定位 token 时出现错误")
	}
	token := tokenContent[p : p+19]

	if checkToken(token) {
		return token, nil
	}

	return "", errors.New("否")
}

// FetchToken 获取token
func FetchToken() {
	req := new(networkware.Networkware)
	req.Address = "https://twitter.com/kaor1n_n"
	req.Method = "GET"
	req.Header = [][]string{{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4128.3 Safari/537.36"}}
	req.Proxy = "http://127.0.0.1:10809"
	res, err := req.Send()
	if err != nil {
		log.Error().Err(err)
		return
	}

	token2, err := fetchToken(res)
	if err != nil {
		log.Info().
			Str("功能", "twitter").
			Msg("第一次获取 token 时出错， 责任下传")

		token2, err = fetchToken2(res)
		if err != nil {
			log.Error().
				Str("功能", "twitter").
				Msg("获取 token 失败")
			return
		}
	}
	token = token2
	log.Info().
		Str("功能", "twitter").
		Str("token", token).
		Msg("成功获取 token")
}

// Timer 定时器
func Timer() {
	FetchToken()
	go func() {
		for range time.Tick(time.Hour) {
			FetchToken()
		}
	}()
}
