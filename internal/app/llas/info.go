package llas

import "github.com/qianjunakasumi/project-shizuku/internal/shizuku"

func init() {

	a := new(randomStill)
	a.initApp()

	shizuku.NewApp(&shizuku.AppInfo{
		Name:        "RandomStill",
		DisplayName: "来一张场景",
		Keys:        []string{"来一张场景", "lyzcj", "still"},
		Expand: []shizuku.Expand{
			{
				"idol",
				"偶像",
				[]string{"偶像", "爱抖露"},
				[]string{},
				false,
				"_SHIZUKU专用",
			},
		},
		Pointer: a,
	})

}

func (r *randomStill) initApp() {

	r.root = "assets/images/llas/"

}
