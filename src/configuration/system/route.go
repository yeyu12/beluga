package system

import (
	drive2 "beluga/src/beluga/drive"
	"beluga/src/configuration/controller"
	"net/http"
)

func syncConfig(writer http.ResponseWriter, request *http.Request) {
	con := new(controller.ConfigController)

	con.Conf = drive2.G_conf
	con.Monitor = drive2.G_monitor
	con.Redis = drive2.G_redis
	con.Mysql = drive2.G_mysql

	defer func() {
		if err := recover(); err != nil {
			request.Body.Close()
		}
	}()

	con.Requests(writer, request)
	con.SyncConfig()

	request.Body.Close()
}

func getConfig(writer http.ResponseWriter, request *http.Request) {
	con := new(controller.ConfigController)

	con.Conf = drive2.G_conf
	con.Monitor = drive2.G_monitor
	con.Redis = drive2.G_redis
	con.Mysql = drive2.G_mysql
	//con.Etcd = drive.G_etcd

	defer func() {
		if err := recover(); err != nil {
			request.Body.Close()
		}
	}()

	con.Requests(writer, request)
	con.GetConf()

	request.Body.Close()
}