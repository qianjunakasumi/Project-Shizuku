/*
main.go: 消息链
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

package messageChain

type MessageChain struct {
	Content []map[string]interface{}
}

func (m *MessageChain) AddText(str string) {
	m.Content = append(m.Content, map[string]interface{}{
		"type": "Plain",
		"text": str,
	})
}

func (m *MessageChain) AddAt(target uint32) {
	m.Content = append(m.Content, map[string]interface{}{
		"type":    "At",
		"target":  target,
		"display": "@",
	})
}
