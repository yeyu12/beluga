package drive

import (
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	// 配置中心对外配置类型
	NODE_CONF_HTTP_TYPE = "http"
	NODE_CONF_JSON_TYPE = "json"
	NODE_CONF_PHP_TYPE  = "php"
)

var G_conf *Config
var G_redis *redis.Client
var G_monitor *Monitor
var G_mysql *gorm.DB
var G_node_conf map[string]string
