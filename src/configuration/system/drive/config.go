package drive

import (
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"log"
)

type Config struct {
	Cfg *ini.File
}

const (
	// 配置文件名
	CONFIG_FILE_NAME = "configuration_node.ini"
)

// 配置目录所在位置
var CONFIG_DIR = "./config/"
var config_file_path string


func InitConfig() {
	config_file_path = CONFIG_DIR + CONFIG_FILE_NAME
	cfg, err := ini.Load(config_file_path)

	if err != nil {
		log.Fatal(errors.Wrap(err, "配置文件读取错误"))
	}

	config_obj := new(Config)
	config_obj.Cfg = cfg

	G_conf = config_obj
}