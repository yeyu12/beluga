package helpers

import (
	"github.com/urfave/cli"
	"net"
	"path"
	"runtime"
	"time"
)

// 获取当前目录
func GetCurrentDirectory() string {
	_, filename, _, ok := runtime.Caller(1)
	var cwdPath string
	if ok {
		cwdPath = path.Join(path.Dir(filename), "")
	}  else  {
		cwdPath = "./"
	}

	return cwdPath
}

// 获取本机IP
func GetLocalIp() (ipv4 string, err error) {
	// 获取所有网卡
	addrs, err := net.InterfaceAddrs();
	if  err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		ipNet, isIpNet := addr.(*net.IPNet);
		if  isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	return
}

// 初始化线程数
func InitThreadNum(num int) {
	runtime.GOMAXPROCS(num)
}

func StringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func BoolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func IntFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func DurationFlag(name string, value time.Duration, usage string) cli.DurationFlag {
	return cli.DurationFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
