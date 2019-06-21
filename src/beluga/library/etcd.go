package library

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"beluga/src/beluga/helpers"
	"context"
)

type EtcdClient struct {
	Client  *clientv3.Client
	Kv      clientv3.KV
	Lease   clientv3.Lease
	LocalIp string
}

var (
	G_conf_etcd_client *EtcdClient
)

// 初始化
func InitRegister(etcd_ip []string, timeoute int) (err error) {
	//etcd 初始化配置
	etcd_config := clientv3.Config{
		Endpoints:   etcd_ip,                                    // 集群地址
		DialTimeout: time.Duration(timeoute) * time.Millisecond, // 超时时间,毫秒为单位
	}

	// 建立链接
	etcd_client, err := clientv3.New(etcd_config)
	if err != nil {
		return
	}

	// 获取本机IP
	loca_ip, err := helpers.GetLocalIp()
	if err != nil {
		return
	}

	G_conf_etcd_client = &EtcdClient{
		Client:  etcd_client,
		Kv:      clientv3.NewKV(etcd_client),
		Lease:   clientv3.NewLease(etcd_client),
		LocalIp: loca_ip,
	}

	return
}

// 服务注册,并且自动续租
func ServerRegister(dir, val string) {
	key := dir + G_conf_etcd_client.LocalIp

	for {
		// 创建租约
		lease_grant_resp, err := G_conf_etcd_client.Lease.Grant(context.TODO(), 10)
		if err != nil {
			continue
		}

		// 自动续租
		keep_alive_chan, err := G_conf_etcd_client.Lease.KeepAlive(context.TODO(), lease_grant_resp.ID)
		if err != nil {
			continue
		}

		// 数据写入
		if _, err = G_conf_etcd_client.Kv.Put(context.TODO(), key, val, clientv3.WithLease(lease_grant_resp.ID)); err != nil {
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
}