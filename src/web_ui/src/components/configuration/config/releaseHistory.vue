<template>
    <el-card class="box-card" id="config-box-card">
        <div id="release-history-top-one">
            <p>发布历史</p>
        </div>

        <div id="release-history-top-two" style="text-align: center;">
            <div id="release-history-top-two-left" class="release-history-top-two">
                <span style="background: #316510;">主版本发布</span>
                <span style="background: #997f1c;margin-left:3px;">主版本回滚</span>
                <span style="background: #999999;margin-left:3px;">灰度操作</span>
            </div>
            <div style="color:#797979">{{project_name}}，{{namespace_name}}</div>
            <div id="release-history-top-two-right">
                <el-button type="primary" size="mini" @click="history_go">返回配置页</el-button>
            </div>
        </div>
        <div class="cut-off-rule"></div>
        <div id="release-history-main">
            <el-scrollbar id="release-history-main-left">
                <release-history-list :data="history_list" @selected="selectedHistory"></release-history-list>
            </el-scrollbar>
            <div id="release-history-main-right">
                <div id="release-history-main-right-top">
                    <div id="release-history-info">
                        <span id="release-history-version">{{releaseVersion}}</span>
                        <span id="release-history-date">{{releaseDate}}</span>
                    </div>
                    <div id="release-history-config">
                        <el-button-group>
                            <el-button size="mini" @click="switchConfigViewType(1)">变更的配置</el-button>
                            <el-button size="mini" @click="switchConfigViewType(2)">全部配置</el-button>
                        </el-button-group>
                    </div>
                </div>
                <div id="release-history-main-right-main">
                    <template v-if="isShow == 1">
                        <div id="release-history-main-right-main-msg">
                            变更的配置
                        </div>
                        <!-- 表格展示数据 -->
                        <el-table :data="configChange" border style="width: 100%" height="calc(100% - 45px)">
                            <el-table-column prop="key" label="Key"></el-table-column>
                            <el-table-column prop="old_val" label="Old Val"></el-table-column>
                            <el-table-column prop="val" label="New Val"></el-table-column>
                            <el-table-column label="Type">
                                <template slot-scope="scope">
                                    <span v-if="scope.row.type == -1">删除</span>
                                    <span v-else-if="scope.row.type == 1">新增</span>
                                    <span v-else-if="scope.row.type == 2">修改</span>
                                    <span v-else>不知道什么鬼</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </template>
                    <template v-else-if="isShow == 2">
                        <div id="release-history-main-right-main-msg">
                            全部配置
                        </div>
                        <!-- 表格展示数据 -->
                        <el-table :data="configAll" border style="width: 100%" height="calc(100% - 45px)">
                            <el-table-column prop="key" label="Key"></el-table-column>
                            <el-table-column prop="val" label="Val"></el-table-column>
                        </el-table>
                    </template>
                </div>
            </div>
        </div>
    </el-card>
</template>

<script>
import ReleaseHistoryList from "@/components/configuration/config/releaseHistoryList";

export default {
    name: "releaseHistory",
    data() {
        return {
            isShow: 1,
            project_name: "",
            namespace_name: "",
            releaseVersion: "",
            releaseDate: "",
            selectedHistoryData: {},
            history_list: "",
            configChange: [],
            configAll: []
        };
    },
    methods: {
        history_go() {
            this.$router.go(-1);
        },
        selectedHistory(val) {
            this.selectedHistoryData = val;
            this.releaseVersion = val.version;
            this.releaseDate = val.create_time;

            this.getReleaseConfigChange(val.version);
        },
        switchConfigViewType(type) {
            this.isShow = type;
            if (type == 1) {
                this.getReleaseConfigChange(this.releaseVersion);
            } else if (type == 2) {
                this.getReleaseConfigList(this.releaseVersion);
            }
        },
        getProjectNameNamespaceName() {
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id
            };
            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_project_name_namespace_name,
                    data,
                    {}
                )
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.project_name = data.data.project_name;
                        _this.namespace_name = data.data.namespace_name;
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        },
        getConfigOperationList() {
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_release_history, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.history_list = data.data;
                        if (data.data && data.data.length > 0) {
                            _this.releaseVersion = data.data[0].version;
                            _this.releaseDate = data.data[0].create_time;
                        }
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        },
        // 获取已发布的操作纪录
        getReleaseConfigChange(version) {
            // 获取纪录
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id,
                version: version
            };
            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_change_release_version_list,
                    data,
                    {}
                )
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.configChange = data.data;
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        },
        // 获取已发布的版本号下面的所有配置
        getReleaseConfigList(version) {
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id,
                version: version
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_all_release_version_list, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        console.log(data);
                        _this.configAll = data.data;
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        }
    },
    components: {
        ReleaseHistoryList
    },
    watch: {
        selectedHistoryData(val, oldVal) {
            this.isShow = 1;
        }
    },
    created() {
        this.$common.addVuexTopInfoData(this, {
            name: "发布历史",
            path: "releaseHistory",
            params: {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id
            }
        });

        // 获取项目和命名空间名称
        this.getProjectNameNamespaceName();

        this.getConfigOperationList();
    }
};
</script>