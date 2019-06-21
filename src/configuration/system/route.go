package system

import (
	"net/http"
	"beluga/src/configuration/controller"
	"beluga/src/configuration/system/drive"
)

func syncConfig(writer http.ResponseWriter, request *http.Request) {
	con := new(controller.ConfigController)

	con.Conf = drive.G_conf
	con.Monitor = drive.G_monitor
	con.Redis = drive.G_redis
	con.Mysql = drive.G_mysql

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

	con.Conf = drive.G_conf
	con.Monitor = drive.G_monitor
	con.Redis = drive.G_redis
	con.Mysql = drive.G_mysql
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