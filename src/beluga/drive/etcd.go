package drive

import (
	"beluga/src/beluga/library"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func InitEtcd() {
	cfg := G_conf.Cfg
	ip := cfg.Section("edit_server").Key("ip").String()
	timeout, _ := cfg.Section("edit_server").Key("timeout").Int()

	etcd_host := strings.Split(ip, ";")

	if err := library.InitRegister(etcd_host, timeout); err != nil {
		Notices(logrus.Fields{}, errors.Wrap(err, "etcd链接失败"))
		os.Exit(0)
		return
	}
}