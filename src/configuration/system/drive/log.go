package drive

import (
	log "github.com/sirupsen/logrus"
	"os"
	"github.com/pkg/errors"
)

func InitLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	cfg := G_conf.Cfg

	debug, err := cfg.Section("configuration_node").Key("debug").Bool()
	if err != nil {
		Notices(log.Fields{}, errors.Wrap(err, "缺少Debug配置"))
		debug = false
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log_path := cfg.Section("configuration_node").Key("log_path").String()
	log_file_name := "configuration_node.log"

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
		file, _ := os.OpenFile(log_path + log_file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		log.SetOutput(file)
	}
}

func Notices(fields map[string]interface{}, msg interface{})  {
	log.WithFields(fields).Info(msg)
}

func Waring(fields map[string]interface{}, msg interface{}) {
	log.WithFields(fields).Warn(msg)
}

func Debug(fields map[string]interface{}, msg interface{})  {
	log.WithFields(fields).Debug(msg)
}

func Err(fields map[string]interface{}, msg interface{})  {
	log.WithFields(fields).Error(msg)
	os.Exit(0)
}