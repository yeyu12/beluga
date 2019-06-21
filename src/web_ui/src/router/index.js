import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/login/login'
import Home from '@/components/home/main'
import Indexs from '@/components/index/index'
import Configuration from '@/components/configuration/configuration'
import Quque from '@/components/quque/quque'
import Crontab from '@/components/crontab/crontab'
import Namespace from '@/components/configuration/namespace'
import Config from '@/components/configuration/config'
import ReleaseHistory from '@/components/configuration/config/releaseHistory';
import Etcd from '@/components/system/etcd/etcd';
import NodeConf from '@/components/configuration/node_conf';
import Node from '@/components/configuration/node';
import System from '@/components/system/system'
import OperationLog from '@/components/system/operationLog'
import Account from '@/components/account/index'
import UserList from '@/components/system/user/userlist'

Vue.use(Router);

export default new Router({
    routes: [{
            path: '/',
            name: 'home',
            component: Home,
            children: [{
                    path: "/",
                    name: "index",
                    component: Indexs,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/configuration",
                    name: "configuration",
                    component: Configuration,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/nodeConf",
                    name: "nodeConf",
                    component: NodeConf,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/node",
                    name: "node",
                    component: Node,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/quque",
                    name: "quque",
                    component: Quque,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/cron",
                    name: "cron",
                    component: Crontab,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/etcd",
                    name: "etcd",
                    component: Etcd,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/configuration/:project_id/:project_name/namesapce",
                    name: "namespace",
                    component: Namespace,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/configuration/:project_id/namesapce/:namespace_id/:name/:project_name/config",
                    name: "config",
                    component: Config,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/configuration/:project_id/namesapce/:namespace_id/releaseHistory",
                    name: "releaseHistory",
                    component: ReleaseHistory,
                    meta: {
                        requireAuth: true,
                    },
                }, {
                    path: "/sys",
                    name: "sys",
                    component: System,
                    meta: {
                        requireAuth: true,
                    },
                }, {
                    path: "/operationLog",
                    name: "operationLog",
                    component: OperationLog,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/account",
                    name: "account",
                    component: Account,
                    meta: {
                        requireAuth: true,
                    },
                },
                {
                    path: "/user",
                    name: "user",
                    component: UserList,
                    meta: {
                        requireAuth: true,
                    },
                },
            ]
        },
        {
            path: '/login',
            component: Login,
            name: 'Login',
        },
    ]
})
