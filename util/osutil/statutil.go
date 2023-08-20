package osutil

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/fzf-labs/fpkg/conv"
	"github.com/fzf-labs/fpkg/util/cmdutil"
	"github.com/fzf-labs/fpkg/util/iputil"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

func GetFmtSize(data uint64) string {
	var factor float64 = 1024
	res := float64(data)
	for _, unit := range []string{"", "K", "M", "G", "T", "P"} {
		if res < factor {
			return fmt.Sprintf("%.2f%sB", res, unit)
		}
		res /= factor
	}
	return fmt.Sprintf("%.2f%sB", res, "P")
}

func GetFmtTime(ms int64) (res string) {
	rem := ms / 1000
	days, rem := rem/86400, rem%86400
	hours, rem := rem/3600, rem%3600
	minutes := rem / 60
	res = strconv.FormatInt(minutes, 10) + "分钟"
	if hours > 0 {
		res = strconv.FormatInt(hours, 10) + "小时" + res
	}
	if days > 0 {
		res = strconv.FormatInt(days, 10) + "天" + res
	}
	return res
}

type CPUInfo struct {
	VendorID string `json:"vendorID"` // CPU制造商ID
	CPUModel string `json:"cpuModel"` // CPU具体型号
	CoreNum  string `json:"coreNum"`  // 核心数
	Percent  string `json:"percent"`  // 百分比
}

// GetCPUInfo 获取CPU信息
func GetCPUInfo() (*CPUInfo, error) {
	result := &CPUInfo{
		VendorID: "",
		CPUModel: "",
		CoreNum:  "",
		Percent:  "",
	}
	coreNum, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}
	result.CoreNum = strconv.Itoa(coreNum)
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	if len(cpuInfo) > 0 {
		result.VendorID = cpuInfo[0].VendorID
		result.CPUModel = cpuInfo[0].ModelName
	}
	percent, _ := cpu.Percent(time.Second, false)
	if len(percent) > 0 {
		result.Percent = strconv.FormatFloat(percent[0], 'f', 2, 64) + "%"
	}
	return result, nil
}

type MemoryInfo struct {
	Total string `json:"total"` // 总占用
	Used  string `json:"used"`  // 已使用
	Free  string `json:"free"`  // 可用的
	Usage string `json:"usage"` // 使用占比
}

// GetMemInfo 获取内存信息
func GetMemInfo() (*MemoryInfo, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return &MemoryInfo{
		Total: GetFmtSize(memory.Total),
		Used:  GetFmtSize(memory.Used),
		Free:  GetFmtSize(memory.Available),
		Usage: strconv.Itoa(conv.Int(memory.UsedPercent)) + "%",
	}, nil
}

type DiskInfo struct {
	DirName     string `json:"dirName"`     // 盘符路径
	SysTypeName string `json:"sysTypeName"` // 文件系统
	TypeName    string `json:"typeName"`    // 盘符类型
	Total       string `json:"total"`       // 总占用
	Used        string `json:"used"`        // 已使用
	Free        string `json:"free"`        // 可用的
	Usage       string `json:"usage"`       // 使用占比
}

// GetDiskInfo 获取磁盘信息
func GetDiskInfo() ([]DiskInfo, error) {
	result := make([]DiskInfo, 0)
	partStats, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(partStats); i++ {
		part := partStats[i]
		usage, uErr := disk.Usage(part.Mountpoint)
		if uErr != nil {
			continue
		}
		result = append(result, DiskInfo{
			DirName:     part.Mountpoint,
			SysTypeName: part.Fstype,
			TypeName:    part.Device,
			Total:       GetFmtSize(usage.Total),
			Used:        GetFmtSize(usage.Used),
			Free:        GetFmtSize(usage.Free),
			Usage:       strconv.Itoa(conv.Int(usage.UsedPercent)) + "%",
		})
	}
	return result, nil
}

type SysInfo struct {
	ComputerName string `json:"computerName"` // 服务器名称
	LocalIP      string `json:"localIP"`      // 内网ip
	PublicIP     string `json:"publicIP"`     // 外网ip
	Os           string `json:"os"`           // 系统类型
	Arch         string `json:"arch"`         // 系统架构

	GoVersion        string `json:"goVersion"`        // golang 版本
	NpmVersion       string `json:"npmVersion"`       // npm 版本
	NodeVersion      string `json:"nodeVersion"`      // node 版本
	ProjectPath      string `json:"projectPath"`      // 项目地址
	ProjectStartTime string `json:"projectStartTime"` // 项目启动时间
	ProjectRunTime   string `json:"projectRunTime"`   // 项目运行时间

}

// GetSysInfo 获取服务器信息
func GetSysInfo() (*SysInfo, error) {
	infoStat, err := host.Info()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	npmVersion, _, err := cmdutil.ExecCommand("npm -v")
	if err != nil {
		return nil, err
	}
	nodeVersion, _, err := cmdutil.ExecCommand("node -v")
	if err != nil {
		return nil, err
	}
	curProc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}
	startTime, err := curProc.CreateTime()
	if err != nil {
		return nil, err
	}
	publicIP, err := iputil.GetPublicIPByHTTP()
	if err != nil {
		return nil, err
	}
	return &SysInfo{
		ComputerName:     infoStat.Hostname,
		LocalIP:          iputil.GetLocalIP(),
		PublicIP:         publicIP,
		Os:               infoStat.OS,
		Arch:             infoStat.KernelArch,
		GoVersion:        runtime.Version(),
		NpmVersion:       npmVersion,
		NodeVersion:      nodeVersion,
		ProjectPath:      pwd,
		ProjectStartTime: time.UnixMilli(startTime).Format("2006-01-02 15:04:05"),
		ProjectRunTime:   GetFmtTime(time.Now().UnixMilli() - startTime),
	}, nil
}
