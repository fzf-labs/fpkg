package monitor

import (
	"fmt"
	"os"
	"time"

	"mosn.io/holmes"
)

func Holmes() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("5s"), // 每5s采集一次当前应用的各项指标，该值建议设置为大于1s。
		holmes.WithDumpPath("/tmp"),      // profile文件保存路径。
		holmes.WithTextDump(),            // 以文本格式保存profile内容。
		holmes.WithCPUDump(10, 25, 80, time.Minute),
		holmes.WithMemDump(30, 25, 80, time.Minute),
		holmes.WithGCHeapDump(10, 20, 40, time.Minute),
		holmes.WithGoroutineDump(500, 25, 20000, 0, time.Minute),
		holmes.WithCGroup(isCGroup()), // set cgroup to true
	)
	h.EnableCPUDump().EnableGoroutineDump().EnableMemDump().EnableGCHeapDump().Start()
}

type HolmesReporter struct{}

func NewHolmesReporter() *HolmesReporter {
	return &HolmesReporter{}
}

//nolint:gocritic
func (r *HolmesReporter) Report(pType string, buf []byte, reason string, eventID string) error {
	fmt.Println(pType)
	fmt.Println(buf)
	fmt.Println(reason)
	fmt.Println(eventID)
	return nil
}

func isCGroup() bool {
	if _, err := os.Stat("/proc/self/cgroup"); err == nil {
		return true
	}
	return false
}
