package cmd

import (
	"beluga/src/beluga/configuration_constant"
	"beluga/src/beluga/helpers"
	"beluga/src/beluga/library"
	"beluga/src/beluga/task_constant"
	web_server_helpers "beluga/src/web_server/helpers"
	"beluga/src/web_server/models"
	_ "beluga/src/web_server/routers"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Start = cli.Command{
	Name:        "start",
	Usage:       "启动服务",
	Description: "服务启动",
	Action:      start,
	Flags: []cli.Flag{
		helpers.StringFlag("configDir, c", "conf/", "重定向配置文件路径"),
		helpers.StringFlag("host", "0.0.0.0", "监听地址"),
		helpers.StringFlag("port, p", "9410", "监听端口"),
	},
}
var config_name = "app.conf" // 配置文件名

func start(c *cli.Context) {
	if c.IsSet("configDir") {
		config_path := c.String("configDir") + config_name
		if err := beego.LoadAppConfig("ini", config_path); err != nil {
			beego.Error("配置文件不存在。", config_path)
			os.Exit(0)
		}
	} else {
		if err := beego.LoadAppConfig("ini", "conf/"+config_name); err != nil {
			beego.Error("当前路径下配置文件不存在")
			os.Exit(0)
		}
	}
	if c.IsSet("host") {
		beego.BConfig.Listen.HTTPAddr = c.String("host")
	}
	if c.IsSet("port") {
		beego.BConfig.Listen.HTTPPort, _ = strconv.Atoi(c.String("port"))
	}

	Run()
}

func Run() {
	// 初始化线程
	helpers.InitThreadNum(runtime.NumCPU())

	// 初始化etcd
	initEtcd()

	// 初始化数据库
	initDatabase()

	// 判断中心是否存在，如果存在则直接退出，不存在则继续走
	isMaster()

	// TODO 如果断线的情况的话，定时检测
	initConfNodeRegister()
	initTaskNodeRegister()

	// 监听配置节点
	watchConfNode()

	// 监听任务节点
	watchTaskNode()

	// 启动服务
	beego.Run()
}

// 初始化数据库
func initDatabase() bool {
	mysql_url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		beego.AppConfig.DefaultString("db_username", "root"),
		beego.AppConfig.DefaultString("db_password", "root"),
		beego.AppConfig.DefaultString("db_host", "127.0.0.1"),
		beego.AppConfig.DefaultInt("db_port", 3306),
		beego.AppConfig.DefaultString("db_database", "beluga"),
		beego.AppConfig.DefaultString("db_charset", "utf8"),
	)
	mysql_url += "&loc=Local"

	if err := orm.RegisterDataBase("default", "mysql", mysql_url, 30); err != nil {
		beego.Error("初始化mysql错误", err)
		os.Exit(0)
		return false
	}

	debug := beego.AppConfig.DefaultBool("debug", false)
	if debug {
		orm.Debug = true
	}

	orm.DefaultTimeLoc = time.Local
	registerModel()

	return true
}

// 注册模型
func registerModel() {
	orm.RegisterModelWithPrefix(
		web_server_helpers.GetTablePrefix(),
		new(models.Account),
		new(models.ConfigurationProject),
		new(models.ConfigurationNamespace),
		new(models.ConfigurationVersion),
		new(models.ConfigurationLog),
		new(models.Configuration),
		new(models.ConfigurationOperation),
		new(models.ConfigurationOperationLog),
		new(models.ConfigurationNodeConf),
		new(models.ConfigurationNode),
		new(models.OperationLog),
		new(models.Task),
		new(models.TaskNode),
		new(models.TaskLog),
	)
}

// 判断中心服务是否存在。如果存在则禁止启动，保证中心接口服务的唯一性。
func isMaster() {
	debug := beego.AppConfig.DefaultBool("debug", false)
	if debug {
		return
	}

	key := configuration_constant.BELUGA_MASTER_DIR
	get_resp, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), key)
	if err != nil {
		beego.Error("启动失败！")
		os.Exit(0)
	}

	if len(get_resp.Kvs) != 0 {
		beego.Error("中心接口服务存在。禁止启动。")
		os.Exit(0)
	}

	go func() {
		for {
			// 创建租约
			lease_grant_resp, err := library.G_conf_etcd_client.Lease.Grant(context.TODO(), 3)
			if err != nil {
				continue
			}

			// 自动续租
			keep_alive_chan, err := library.G_conf_etcd_client.Lease.KeepAlive(context.TODO(), lease_grant_resp.ID)
			if err != nil {
				continue
			}

			// 数据写入
			if _, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), key, "", clientv3.WithLease(lease_grant_resp.ID)); err != nil {
				continue
			}

			// 续租回应
			for {
				select {
				case keep_alive_resq := <-keep_alive_chan:
					if keep_alive_resq == nil {
						break
					}
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

// 获取etcd节点配置
func initEtcd() {
	etcd_conf_ip := beego.AppConfig.DefaultStrings("etcd_host", []string{})
	etcd_conf_timeoute := beego.AppConfig.DefaultInt("etcd_timeoute", 5000)

	if err := library.InitRegister(etcd_conf_ip, etcd_conf_timeoute); err != nil {
		beego.Error("etcd链接失败", err)
		os.Exit(0)
		return
	}
}

// 监控配置节点
func watchConfNode() {
	key := configuration_constant.CONFIGURATION_REGISTER_DIR
	watch_resp_chan := clientv3.NewWatcher(library.G_conf_etcd_client.Client).Watch(context.TODO(), key, clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				beego.Error(v.Err())
			}

			for _, resp := range v.Events {
				key_spl := strings.Split(string(resp.Kv.Key), "/")
				if len(key_spl) == 3 {
					continue
				}

				node_conf_key := configuration_constant.NODE_CONF_DIR + key_spl[3]

				switch resp.Type {
				case mvccpb.PUT:
					configuration_node_model := models.NewConfigurationNode()

					configuration_node_model.Ip = key_spl[3]
					configuration_node_model.CreateTime = time.Now()
					configuration_node_model.IsDelete = 1

					id, err := configuration_node_model.Save(orm.NewOrm())

					if err != nil {
						beego.Error(err, "节点写入失败")
					}

					// 更新节点中的配置数据
					node_data := configuration_node_model.IdFind(orm.NewOrm(), id)
					var conf_data []map[string]interface{}
					var node_conf_etcd []byte

					if node_data.NodeConfId != "" {
						configuration_node_conf_model := models.NewConfigurationNodeConf()
						node_conf_data := configuration_node_conf_model.IdsToData(orm.NewOrm(), node_data.NodeConfId)

						for _, val := range node_conf_data {
							var conf_map []map[string]interface{}
							json.Unmarshal([]byte(val["conf"].(string)), &conf_map)
							conf_data = append(conf_data, conf_map...)
						}

						node_conf_etcd, _ = json.Marshal(conf_data)
					}

					_, err = library.G_conf_etcd_client.Kv.Put(context.TODO(), node_conf_key, string(node_conf_etcd))
					if err != nil {
						beego.Error(err, "etcd节点配置写入失败")
					}

					break
				case mvccpb.DELETE:
					configuration_node_model := models.NewConfigurationNode()

					configuration_node_model.Ip = key_spl[3]
					configuration_node_model.IsDelete = 0

					configuration_node_model.Edit(orm.NewOrm(), map[string]interface{}{
						"is_delete": 0,
					})

					// 删除节点中的数据
					_, err := library.G_conf_etcd_client.Kv.Delete(context.TODO(), node_conf_key)
					if err != nil {
						beego.Error(err, "etcd节点中数据删除失败")
					}

					break
				}
			}
		}
	}()
}

// 监控任务节点
func watchTaskNode() {
	key := task_constant.TASK_REGISTER_DIR
	watch_resp_chan := clientv3.NewWatcher(library.G_conf_etcd_client.Client).Watch(context.TODO(), key, clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				beego.Error(v.Err())
			}

			for _, resp := range v.Events {
				key_spl := strings.Split(string(resp.Kv.Key), "/")
				if len(key_spl) == 3 {
					continue
				}

				switch resp.Type {
				case mvccpb.PUT:
					task_node_model := models.NewTaskNode()

					task_node_model.Ip = key_spl[3]
					task_node_model.CreateTime = time.Now()
					task_node_model.IsDelete = 1

					_, err := task_node_model.Save(orm.NewOrm())
					if err != nil {
						beego.Error("节点写入失败")
					}

					break
				case mvccpb.DELETE:
					task_node_model := models.NewTaskNode()

					task_node_model.Ip = key_spl[3]
					task_node_model.IsDelete = 0

					task_node_model.Edit(orm.NewOrm(), map[string]interface{}{
						"is_delete": 0,
					})

					break
				}
			}
		}
	}()
}

// 配置节点注册初始化
func initConfNodeRegister() {
	var page int = 1
	configuration_node_model := models.NewConfigurationNode()

	for {
		config_node_list := configuration_node_model.List(orm.NewOrm(), page, 10, "")

		if len(config_node_list.List.([]orm.Params)) == 0 {
			break
		}

		for _, val := range config_node_list.List.([]orm.Params) {
			configuration_node_model.Ip = val["ip"].(string)
			key := configuration_constant.CONFIGURATION_REGISTER_DIR + val["ip"].(string)
			get_res, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), key)
			if err != nil {
				beego.Error("节点数据初始化失败", err)
			}

			if get_res.Kvs == nil {
				configuration_node_model.Edit(orm.NewOrm(), map[string]interface{}{"is_delete": 0})
			} else {
				for _ = range get_res.Kvs {
					configuration_node_model.Edit(orm.NewOrm(), map[string]interface{}{"is_delete": 1})
				}
			}
		}

		page++
	}
}

// 任务节点注册初始化
func initTaskNodeRegister() {
	var page int = 1
	task_node_model := models.NewTaskNode()

	for {
		config_node_list := task_node_model.List(orm.NewOrm(), page, 10, "")

		if len(config_node_list.List.([]orm.Params)) == 0 {
			break
		}

		for _, val := range config_node_list.List.([]orm.Params) {
			task_node_model.Ip = val["ip"].(string)
			key := task_constant.TASK_REGISTER_DIR + val["ip"].(string)
			get_res, err := library.G_conf_etcd_client.Kv.Get(context.TODO(), key)
			if err != nil {
				beego.Error("节点数据初始化失败", err)
			}

			if get_res.Kvs == nil {
				task_node_model.Edit(orm.NewOrm(), map[string]interface{}{"is_delete": 0})
			} else {
				for _ = range get_res.Kvs {
					task_node_model.Edit(orm.NewOrm(), map[string]interface{}{"is_delete": 1})
				}
			}
		}

		page++
	}
}
