package drive

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Mysql struct {
	Db *gorm.DB
}

func InitMysql() {
	cfg := G_conf.Cfg
	mysql_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		cfg.Section("mysql").Key("username").String(),
		cfg.Section("mysql").Key("password").String(),
		cfg.Section("mysql").Key("host").String(),
		cfg.Section("mysql").Key("port").String(),
		cfg.Section("mysql").Key("database").String(),
		cfg.Section("mysql").Key("charset").String(),
	)
	mysql_url += "&loc=Local"

	db, err := gorm.Open("mysql", mysql_url)
	if err != nil {
		Err(logrus.Fields{}, errors.Wrap(err, "数据库链接失败！"))
	}

	max_count, err := cfg.Section("mysql").Key("max_count").Int()
	if err != nil {
		max_count = 1
	}
	db.DB().SetMaxOpenConns(max_count)
	db.DB().SetMaxIdleConns(10)

	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	G_mysql = db

	//go dbHear()
}

func dbHear() {
	for {
		if G_mysql.DB().Ping() != nil {
			InitMysql()
		} else {
			time.Sleep(2 * time.Minute)
		}
	}
}
