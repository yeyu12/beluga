// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import ElementUI from 'element-ui';
import axios from 'axios';
import VueAxios from 'vue-axios';
import Qs from 'qs';
import Common from '@/assets/js/common.js';
import Cookie from '@/assets/js/cookie.js';
import Server from '@/assets/js/server.js';

import 'element-ui/lib/theme-chalk/index.css';
import '@/assets/css/main.css';
import '@/assets/icon/iconfont.css';

Vue.use(ElementUI);
Vue.use(VueAxios, axios.create({
    transformRequest: [function(data) {
        data = Qs.stringify(data);
        return data;
    }]
}));

Vue.config.productionTip = false;
Vue.prototype.$axios = axios;
Vue.prototype.$common = Common;
Vue.prototype.$cookie = Cookie;
Vue.prototype.$server = Server;

axios.defaults.timeout = 30000;

// http请求拦截器
// var loadinginstace
axios.interceptors.request.use(
    config => {
        let token = Cookie.getCookie("beluga_token");
        if ((token == "") || (token == null) || (token == "undefined") || (token == undefined)) {
            token = ""
        }
        config.url = Server.host + config.url

        // loadinginstace = ElementUI.Loading.service({ fullscreen: true })
        if (config.url != ("http://" + Server.host + Server.login)) { // 判断是否存在token，如果存在的话，则每个http header都加上token，请求头处理
            config.headers["Account-token"] = `${token}`;
        }

        config.headers["X-Requested-With"] = "XMLHttpRequest";
        config.headers["Content-Type"] = "application/json; charset=UTF-8";

        return config;
    },
    err => {
        return Promise.reject(err);
    });

axios.interceptors.response.use(
    response => {
        // 拦截token问题
        let code = [
            5000, 5003, 5004, 5006
        ];

        if (code.indexOf(response.data.code) != -1) {
            Cookie.delCookie("beluga_token")
        }

        return response;
    },
    error => {
        if (error.response) {
            switch (error.response.status) {
                // 错误提示
            }
        }
        return Promise.reject(error.response.data) // 返回接口返回的错误信息
    });

router.beforeEach((to, from, next) => {
    if (to.meta.requireAuth) { // 判断该路由是否需要登录权限
        if (Cookie.getCookie("beluga_token")) { // 获取当前的token是否存在
            next();
        } else {
            Cookie.delCookie("beluga_token");
            next({
                path: '/login',
                // query: { redirect: to.fullPath } // 将跳转的路由path作为参数，登录成功后跳转到该路由
            })
        }
    } else {
        next();
    }
});

// TODO 在顶部需要加入加载条
/* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    store,
    components: {
        App
    },
    template: '<App/>'
});