package drive

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/net"
	"time"
	"runtime"
)

type Monitor struct {
	Info *monitorInfo
}

type monitorInfo struct {
	Cpu       float64
	Mem       float64
	DiskIo    map[string]disk.IOCountersStat
	NetIo     []net.IOCountersStat
	GoRuntime int
}



func InitMonitor() {
	resources_obj := new(Monitor)
	resources_obj.Info = new(monitorInfo)
	cfg := G_conf.Cfg
	monitor_time, err := cfg.Section("configuration_node").Key("monitor_time").Int64()

	if err != nil {
		Notices(logrus.Fields{}, errors.Wrapf(err, "监控间隔时间配置获取错误"))
		monitor_time = 1
	}

	go func() {
		for {
			resources_obj.Info.Cpu = resources_obj.getCpuPercent()
			resources_obj.Info.Mem = resources_obj.getMemPercent()
			resources_obj.Info.DiskIo = resources_obj.getDiskIo()
			resources_obj.Info.NetIo = resources_obj.getNetIo()
			resources_obj.Info.GoRuntime = resources_obj.getGoRuntimeCount()

			time.Sleep(time.Second * time.Duration(monitor_time))
		}
	}()

	G_monitor = resources_obj
}

// cpu使用率
func (c *Monitor) getCpuPercent() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		Notices(logrus.Fields{}, errors.Wrap(err, "cpu占用率获取失败"))
	}

	return percent[0]
}

// 硬盘使用情况
func (c *Monitor) getDiskIo() map[string]disk.IOCountersStat {
	disk_io_info, err := disk.IOCounters()
	if err != nil {
		Notices(logrus.Fields{}, errors.Wrap(err, "硬盘读写速度获取失败"))
	}

	return disk_io_info
}

// 内存使用率
func (v *Monitor) getMemPercent() float64 {
	mem_data, err := mem.VirtualMemory()
	if err != nil {
		Notices(logrus.Fields{}, errors.Wrap(err, "内存占用率获取失败"))
	}

	return mem_data.UsedPercent
}

// 获取网络使用情况
func (c *Monitor) getNetIo() []net.IOCountersStat {
	net_io_info, err := net.IOCounters(true)
	if err != nil {
		Notices(logrus.Fields{}, errors.Wrap(err, "网络读写速度获取失败"))
	}

	return net_io_info
}

// 获取当前goruntime数量
func (c *Monitor) getGoRuntimeCount() int {
	return runtime.NumGoroutine()
}
