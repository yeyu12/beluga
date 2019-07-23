package drive

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
)

func InitRedis() {
	cfg := G_conf.Cfg
	host := cfg.Section("redis").Key("host").String()
	port := cfg.Section("redis").Key("port").String()
	passwd := cfg.Section("redis").Key("password").String()
	db, _ := cfg.Section("redis").Key("select").Int()

	redis_cli := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: passwd,
		DB:       db,
	})

	_, err := redis_cli.Ping().Result()
	if err != nil {
		Err(logrus.Fields{}, errors.Wrap(err, "redis链接失败"))
		os.Exit(0)
	}

	G_redis = redis_cli
}
