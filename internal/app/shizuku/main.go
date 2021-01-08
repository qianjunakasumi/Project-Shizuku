package shizuku

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/qianjunakasumi/project-shizuku/internal/shizuku"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type sysInfo struct{}

func (s sysInfo) onSubCallByQQ(_ *shizuku.QQMsg, _ map[string]string) (*shizuku.Message, error) {

	cp, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	r := shizuku.NewText("CPU占用：" + strconv.FormatFloat(cp[0], 'f', 2, 64) + " %\n")

	me, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	r.AddText("内存占用：" + strconv.FormatFloat(me.UsedPercent, 'f', 2, 64) + " %\n")
	r.AddText("可用内存：" + strconv.FormatFloat(float64(me.Free)/1073741824, 'f', 4, 64) + " G\n")

	wd, _ := os.Getwd()
	d, err := disk.Usage(wd)
	if err != nil {
		return nil, err
	}

	r.AddText("磁盘占用：" + strconv.FormatFloat(d.UsedPercent, 'f', 2, 64) + " %\n")
	r.AddText("可用存储：" + strconv.FormatFloat(float64(d.Free)/1073741824, 'f', 4, 64) + " G\n")

	hp, _, hv, _ := host.PlatformInformation()
	r.AddText("系统信息：" + hp + " " + hv + "\n")

	t, _ := host.BootTime()
	r.AddText("已开机时间：" + fmt.Sprintf("%v", time.Since(time.Unix(int64(t), 0))))

	return r, nil
}
