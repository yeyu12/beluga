package drive

import (
	"beluga/src/beluga/configuration_constant"
	"beluga/src/beluga/drive"
	"beluga/src/beluga/library"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 监听节点配置修改（中心管理后台中的，节点管理）
func WatchNode() {
	key := configuration_constant.NODE_CONF_DIR + library.G_conf_etcd_client.LocalIp
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "节点配置监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					watchNode(resp)
					break
				case mvccpb.DELETE:
					break

				}
			}
		}
	}()
}

// 监听配置发布，同步操作
func WatchConfigurationReleaseToSync()  {
	key := configuration_constant.CONFIGURATION_CONF_SYNC
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV(), clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "发布同步配置监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					watchConfigurationSync(resp)
					break
				case mvccpb.DELETE:
					break
				}
			}
		}
	}()
}

// 监听节点配置修改（中心管理后台，节点配置）
func WatchNodeConf()  {
	key := configuration_constant.CONFIGURATION_NODE_CONF
	watch_resp_chan := library.G_conf_etcd_client.Client.Watch(context.TODO(), key, clientv3.WithPrevKV(), clientv3.WithPrefix())

	go func() {
		for v := range watch_resp_chan {
			if v.Err() != nil {
				drive.Err(logrus.Fields{}, errors.Wrap(v.Err(), "节点配置修改监听失败"))
			}
			for _, resp := range v.Events {
				switch resp.Type {
				case mvccpb.PUT:
					watchNodeConf(resp)
					break
				case mvccpb.DELETE:
					break
				}
			}
		}
	}()
}