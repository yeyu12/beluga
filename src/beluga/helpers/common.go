package helpers

import (
	"path/filepath"
	"os"
	"log"
	"strings"
	"strconv"
	"io/ioutil"
	"net"
	"runtime"
	"github.com/urfave/cli"
	"time"
)

// 获取当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 获取配置文件目录
func GetIniPath() string {
	return GetCurrentDirectory() + "/config"
}

// 获取pid文件目录
func GetPidPath() string {
	return GetCurrentDirectory() + "/runtime/pid"
}

// pid写入文件
func SavePidFile(filePath string, filename string) bool {
	if filePath == ""{
		filePath = GetPidPath() + "/"
	}
	pid_path := filePath +filename
	pid := strconv.Itoa(os.Getpid())

	if err := ioutil.WriteFile(pid_path, []byte(pid), 0644); err != nil {
		return false
	}

	return true
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
