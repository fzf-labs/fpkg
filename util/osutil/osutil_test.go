package osutil

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"testing"
)

func TestIsWindows(t *testing.T) {
	//fmt.Println(IsWindows())
	//fmt.Println(IsLinux())
	//fmt.Println(IsMac())
	info, err := cpu.Info()
	if err != nil {
		return
	}
	fmt.Println(info)
}
