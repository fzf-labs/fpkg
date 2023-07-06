package monitor

import (
	"time"

	"mosn.io/holmes"
)

func Holmes() {
	h, _ := holmes.New(
		holmes.WithCollectInterval("5s"), //每5s采集一次当前应用的各项指标，该值建议设置为大于1s。
		holmes.WithDumpPath("/tmp"),      //profile文件保存路径。
		holmes.WithTextDump(),            //以文本格式保存profile内容。
		holmes.WithCPUDump(10, 25, 80, time.Minute),
		holmes.WithMemDump(30, 25, 80, time.Minute),
		holmes.WithGCHeapDump(10, 20, 40, time.Minute),
		holmes.WithGoroutineDump(500, 25, 20000, 0, time.Minute),
		holmes.WithCGroup(true), // set cgroup to true
	)
	h.EnableCPUDump().EnableGoroutineDump().EnableMemDump().EnableGCHeapDump().Start()
}

type HolmesReporter struct{}

func NewHolmesReporter() *HolmesReporter {
	return &HolmesReporter{}
}

func (r *HolmesReporter) Report(pType string, buf []byte, reason string, eventID string) error {

	return nil
}
