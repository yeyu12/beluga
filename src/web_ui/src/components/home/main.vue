<template>
    <el-container>
        <el-aside width="250px">
            <div class="main-nav-top" @click="backToHome">
                <img :src="beluga_logo"/>
                <div>
                    <span>{{ beluga_name }}</span>
                </div>
            </div>
            <div class="main-nav">
                <Menus @menu-selection="menu_selection_path"></Menus>
            </div>
        </el-aside>
        <el-container>
            <el-header>
                <template v-for="(val, key) in route_top_info">
                    <span :key="key" class="route_top_info">{{
                        val.name
                        }}</span>
                    <span :key="key + '-'" class="route_top_info_separator"
                          v-if="showSeparator(key)">&nbsp;-&nbsp;</span>
                </template>
                <div id="user-info">
                    <el-dropdown :hide-on-click="false" @command="handleCommand" placement="top">
                        <span class="el-dropdown-link">
                            {{user_info.nickname}}
                            <i class="el-icon-arrow-down el-icon--right"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item command="userSetting">
                                <i class="iconfont el-beluga-icon-test1"></i>
                                用户设置
                            </el-dropdown-item>
                            <el-dropdown-item divided command="logout">
                                <i class="iconfont el-beluga-Logout"></i>
                                退出登录
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </div>
            </el-header>
            <el-main>
                <!-- 需要变更滚动条 -->
                <router-view v-if="isRouterAlive"></router-view>
            </el-main>
            <el-footer style="color: #888;text-align: center">
                ©&nbsp;2019&nbsp;&nbsp;&nbsp;beluga&nbsp;当前版本：0.0.1&nbsp;&nbsp;&nbsp;golang版本：&nbsp;1.10
            </el-footer>
        </el-container>
    </el-container>
</template>

<script>
    import Menus from "@/components/home/menu";
    import Configuration from "@/components/configuration/configuration";

    export default {
        name: "home",
        data() {
            return {
                isRouterAlive: true,
                beluga_logo: require("@/assets/icon/logo.png"),
                beluga_name: "白鲸管理平台",
                index: "",
                beluag_top_info: "",
                route_top_info: "",
                route_top_info_length: 0,
                user_info: {},
            };
        },
        provide() {
            return {
                reload: this.reload
            };
        },
        methods: {
            backToHome() {
                this.$router.replace("/");
            },
            menu_selection_path(data) {
                this.index = data.path;
                this.$router.replace(this.index);
            },
            showSeparator(key) {
                if (key < this.route_top_info_length) {
                    return true;
                } else {
                    return false;
                }
            },
            reload() {
                this.isRouterAlive = false;
                this.$nextTick(() => {
                    this.isRouterAlive = true;
                });
            },
            logout() {
                this.$cookie.delCookie("beluga_token");
                this.$router.replace("/login");
            },
            userSetting(){
                this.$router.replace("/account")
            },
            handleCommand(command) {
                if (typeof (this[command]) == "function") {
                    this[command]()
                }
            },
            getUserInfo() {
                let _this = this;

                this.$axios
                    .post(_this.$server.user_info, {}, {})
                    .then(function (response) {
                        var data = response.data;

                        if (data.status) {
                            _this.user_info = data.data
                        } else {
                            _this.$message("获取用户信息失败");
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            }
        },
        mounted() {
            this.getUserInfo();
        },
        components: {
            Menus,
            Configuration
        },
        computed: {
            routeToTopInfo() {
                return this.$store.state.topInfo.top_info;
            }
        },
        watch: {
            routeToTopInfo(val) {
                this.route_top_info = val;
                this.route_top_info_length = val.length - 1;
            }
        }
    };
</script>
