package drive

import (
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"log"
)

type Config struct {
	Cfg *ini.File
}

var (
	CONFIG_DIR = ""
	CONFIG_FILENAME = ""
)


func InitConfig() {
	cfg, err := ini.Load(CONFIG_DIR + CONFIG_FILENAME)

	if err != nil {
		log.Fatal(errors.Wrap(err, "配置文件读取错误"))
	}

	config_obj := new(Config)
	config_obj.Cfg = cfg

	G_conf = config_obj
}