package system

import (
	"net/http"
	"fmt"
	"beluga/src/configuration/system/drive"
)

func HttpServer(host, port string) {
	cfg := drive.G_conf.Cfg

	if host == "" {
		host = cfg.Section("configuration_node").Key("host").String()
	}
	if port == ""{
		port = cfg.Section("configuration_node").Key("port").String()
	}

	addr := host + ":" + port

	// 系统监控启动、当前goruntime数量，cpu占用了率、内存占用率、网络占用率
	/*go func() {
		for {
			fmt.Println(con.Monitor.Info)
			time.Sleep(time.Second * 1)
		}
	}()*/

	// 初始化其他服务

	http.HandleFunc("/getConfig", getConfig)

	http.HandleFunc("/syncConfig", syncConfig)

	fmt.Println("configuration-node,http服务启动。" + addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		//drive.Err(logrus.Fields{}, errors.Wrap(err, "http服务启动失败！"))
	}
}
