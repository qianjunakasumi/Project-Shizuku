/***********************************************************************************************************************
* Basic:
*
*   Package Name : twitter
*   File Name    : tweets.go
*   File Path    : internal/shizuku/twitter/
*   Author       : Qianjunakasumi
*   Description  : 获取并解析 Tweet
*
*----------------------------------------------------------------------------------------------------------------------*/

package twitter

import "github.com/qianjunakasumi/project-shizuku/internal/shizuku"

// fetchTweet 获取推文实例
type fetchTweet struct{}

// OnCall 处理调用
func (f fetchTweet) OnCall(qm *shizuku.QQMsg, sz *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	if c := qm.Call["idol"]; c != "_SHIZUKU默认检查专用" {
		qm.Type = sz.FuzzyGetIdol(c)
	}

	sb, err := newTweetAPIParser([]byte{})
	if err != nil {
		return
	}

	full := sb.getText(0)
	rm = shizuku.NewText(full)

	return
}
