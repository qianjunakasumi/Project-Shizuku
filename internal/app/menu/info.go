package menu

import (
	"fmt"

	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"
)

type menu struct{}

func init() {

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "menu",
		DisplayName: "菜单",
		Keys:        []string{"菜单"},
		Expand:      shizuku.Expand{},
		Pointer:     new(menu),
	})

}

func (m menu) OnCall(_ *shizuku.QQMsg, _ *shizuku.SHIZUKU) (rm *shizuku.Message, err error) {

	r := shizuku.NewText("> PROJECT-SHIZUKU 菜单\n")

	for _, v := range shizuku.InitAppInfo {
		r.AddText("\n" + v.DisplayName + ":" + fmt.Sprintf("%v", v.Keys))
	}

	return r, nil

}
