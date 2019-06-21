<template>
    <el-card class="box-card" id="config-box-card">
        <!-- 顶部按钮 -->
        <div class="config-top">
            <div class="config-top-left">
                <span>{{project_name}}</span>-<span>{{namespace_name}}</span></span>
            </div>
            <div class="config-top-right">
                <el-button size="mini" icon="iconfont el-beluga-fabu" @click="releaseConfig">发布</el-button>
                <el-button size="mini" icon="iconfont el-beluga-fanhui" @click="rollbackDialogVisible">回滚
                </el-button>
                <el-button size="mini" icon="iconfont el-beluga-lishi" @click="releaseHistory">发布历史</el-button>
            </div>
        </div>
        <!-- end 顶部按钮 -->
        <div class="cut-off-rule"></div>

        <div>
            <div class="config-tabs-menu" v-show="menu_type.config_table">
                <el-button size="mini" icon="iconfont el-beluga-guolv">过滤配置</el-button>
                <el-button size="mini" icon="iconfont el-beluga-tongbu" @click="syncConfig">同步配置</el-button>
                <el-button size="mini" type="primary" icon="iconfont el-beluga-icon-test" @click="addConfigDialogVisible = true">新增配置</el-button>
            </div>
            <div class="config-tabs-menu" v-show="menu_type.conifg_text" style="padding-top: 5px;">
                <template v-if="menu_type_text">
                    <el-tooltip class="item" effect="dark" content="修改" placement="bottom">
                        <i class="iconfont el-beluga-xiugaiziliao" style="font-size: 30px;cursor: pointer;color:#555" @click="menuTextModify"></i>
                    </el-tooltip>
                </template>
                <template v-else>
                    <el-tooltip class="item" effect="dark" content="取消" placement="bottom">
                        <i class="iconfont el-beluga-quxiao" style="font-size: 30px;cursor: pointer;color:#555; font-weight: bold;" @click="menuTextcancel"></i>
                    </el-tooltip>
                    <el-tooltip class="item" effect="dark" content="提交" placement="bottom">
                        <i class="iconfont el-beluga-tijiao" style="font-size: 30px;cursor: pointer;color:#555; font-weight: bold;" @click="menuTextSubmit"></i>
                    </el-tooltip>
                </template>
            </div>
            <!-- <div class="config-tabs-menu" v-show="menu_type.json_text">第三个按钮</div> -->
            <!-- <div class="config-tabs-menu" v-show="menu_type.change_history">第四个按钮</div> -->
        </div>

        <el-tabs v-model="activeName" @tab-click="handleClick">
            <!-- 配置表格 -->
            <el-tab-pane name="configTable">
                <span slot="label">
                    <i class="iconfont el-beluga-biaoge"></i>
                    表格
                </span>
                <el-table :data="config_key_val_table" border style="width: 100%;" height="calc(100vh - 320px)" stripe>
                    <el-table-column label="发布状态">
                        <template slot-scope="scope">
                            <span v-if="scope.row.release_status == 1">
                                <el-button size="mini" type="info">已发布</el-button>
                            </span>
                            <span v-else-if="scope.row.release_status == 0">
                                <el-button size="mini" type="warning">未发布</el-button>
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column label="key" sortable>
                        <template slot-scope="scope">
                            <span v-if="scope.row.type == -1">{{scope.row.key}} <el-button size="mini" type="danger">删</el-button></span>
                            <span v-else-if="scope.row.type == 1">{{scope.row.key}} <el-button size="mini" type="success">新</el-button></span>
                            <span v-else-if="scope.row.type == 2">{{scope.row.key}} <el-button size="mini" type="primary">改</el-button></span>
                            <span v-else>{{scope.row.key}}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="val" label="value"></el-table-column>
                    <el-table-column prop="remake" label="备注"></el-table-column>
                    <el-table-column prop="username" label="最后修改人" sortable></el-table-column>
                    <el-table-column prop="update_time" label="最后修改时间" sortable></el-table-column>
                    <el-table-column label="操作">
                        <template slot-scope="scope">
                            <el-button size="mini" @click="editConfigDialog(scope.$index, scope.row)">
                                编辑
                            </el-button>
                            <el-button size="mini" type="danger" @click="configDel(scope.$index, scope.row)">
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>

            </el-tab-pane>
            <!-- end 配置表格 -->

            <!-- 文本 -->
            <el-tab-pane name="conifgText">
                <span slot="label">
                    <i class="iconfont el-beluga-wenben"></i>
                    文本
                </span>
                <el-input type="textarea" :rows="2" placeholder="请输入内容" v-model="config_text" :readonly="config_text_isread" class="config-text"></el-input>
            </el-tab-pane>
            <!-- end 文本 -->

            <!-- json文本 -->
            <el-tab-pane name="jsonText" v-if="false">
                <span slot="label">
                    <i class="iconfont el-beluga-jsongeshihua"></i>
                    json文本
                </span>
                <el-input type="textarea" :rows="2" placeholder="请输入内容" v-model="config_json_text" class="config-text" disabled></el-input>
            </el-tab-pane>
            <!-- end json文本 -->

            <!-- 更改历史 -->
            <el-tab-pane name="changeHistory" v-if="false">
                <span slot="label">
                    <i class="iconfont el-beluga-lishi"></i>
                    更改历史
                </span>
                <el-scrollbar class="config-change_history">
                    <div v-for="i in 4">
                        <div>
                            <span class="config-change_histore-project_name">项目名</span>
                            <span class="config-change_histore-date">2019-04-03 19:44:27</span>
                        </div>
                        <div>
                            <el-table :data="tableData" border style="width:100%">
                                <el-table-column label="Type">
                                    <template slot-scope="scope">
                                        <span v-if="scope.row.type == 1">正常状态</span>
                                        <span v-else-if="scope.row.type == 2">非正常状态</span>
                                    </template>
                                </el-table-column>
                                <el-table-column prop="key" label="Key">
                                </el-table-column>
                                <el-table-column prop="old_val" label="Old Value"></el-table-column>
                                <el-table-column prop="new_val" label="New Value"></el-table-column>
                                <el-table-column prop="remake" label="Remake"></el-table-column>
                            </el-table>
                        </div>
                        <div class="cut-off-rule" style="margin-top: 15px;margin-bottom: 10px;"></div>
                    </div>
                </el-scrollbar>
            </el-tab-pane>
            <!-- end 更改历史 -->
        </el-tabs>

        <!-- 回滚弹窗 -->
        <el-dialog title="回滚" :visible.sync="centerDialogVisible" center class="config-rollback">
            <p class="config-waring">
                此操作将会回滚到上一个发布版本，且当前版本作废，但不影响正在修改的配置。可在发布历史页面查看当前生效的版本
                <a href="javascript:;" @click="releaseHistory" style="color: #337ab7;text-decoration: none;">点击查看</a>
            </p>
            <el-table :data="rollbackData" border style="width: 100%">
                <el-table-column label="Type">
                    <template slot-scope="scope">
                        <span v-if="scope.row.type == -1">删除</span>
                        <span v-else-if="scope.row.type == 1">新增</span>
                        <span v-else-if="scope.row.type == 2">修改</span>
                        <span v-else>不知道什么鬼</span>
                    </template>
                </el-table-column>
                <el-table-column prop="key" label="Key"></el-table-column>
                <el-table-column prop="val" label="回滚前">
                </el-table-column>
                <el-table-column prop="old_val" label="回滚后">
                </el-table-column>
            </el-table>
            <span slot="footer" class="dialog-footer">
                <el-button @click="centerDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="setRollback">回 滚</el-button>
            </span>
        </el-dialog>
        <!-- end 回滚弹窗 -->

        <!-- 添加配置 -->
        <el-dialog title="添加配置" :visible.sync="addConfigDialogVisible" center>
            <el-form ref="form" :model="form" label-width="80px">
                <el-form-item label="Key" prop="key" :rules="{required: true, message: '键值不能为空'}">
                    <el-input v-model="form.key" placeholder="请输入键值"></el-input>
                </el-form-item>
                <el-form-item label="Value" prop="val" :rules="{required: true, message: '值不能为空'}">
                    <el-input type="textarea" v-model="form.val" :rows="5" placeholder="请输入值"></el-input>
                </el-form-item>
                <el-form-item label="Remake" prop="remake">
                    <el-input type="textarea" v-model="form.remake" placeholder="请输入备注，默认为值。"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addConfigDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="addConfig" :loading="isLoading">提 交</el-button>
            </span>
        </el-dialog>
        <!-- end 添加配置 -->

        <!-- 修改配置 -->
        <el-dialog title="修改配置" :visible.sync="editConfigDialogVisible" center>
            <el-form ref="form" :model="edit" label-width="80px">
                <el-form-item label="Value" prop="val" :rules="{required: true, message: '值不能为空'}">
                    <el-input type="textarea" v-model="edit.val" :rows="5" placeholder="请输入值"></el-input>
                </el-form-item>
                <el-form-item label="Remake" prop="remake">
                    <el-input type="textarea" v-model="edit.remake" placeholder="请输入备注，默认为值。"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editConfigDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="editConfig" :loading="isLoading">提 交</el-button>
            </span>
        </el-dialog>
        <!-- end 修改配置 -->
    </el-card>
</template>

<script>
export default {
    name: "config",
    data() {
        return {
            activeName: "configTable",
            config_text: "",
            config_text_isread: true,
            config_json_text: "待完成",
            centerDialogVisible: false,
            addConfigDialogVisible: false,
            editConfigDialogVisible: false,
            menu_type: {
                config_table: true,
                conifg_text: false,
                json_text: false,
                change_history: false
            },
            menu_type_text: true,
            form: {
                key: "",
                val: "",
                remake: ""
            },
            edit: {},
            rollbackData: [],
            tableData: [],
            config_key_val_table: [],
            isLoading:false,
            namespace_name:this.$route.params.name,
            project_name:this.$route.params.project_name,
        };
    },
    methods: {
        initData() {
            this.config_text_isread = true;
            this.centerDialogVisible = false;
            this.addConfigDialogVisible = false;
            this.menu_type_text = true;
        },
        menuDisplay() {
            this.menu_type.config_table = false;
            this.menu_type.conifg_text = false;
            this.menu_type.json_text = false;
            this.menu_type.change_history = false;
        },
        menuTextModify() {
            this.menu_type_text = false;
            this.config_text_isread = false;
        },
        menuTextcancel() {
            this.menu_type_text = true;
            this.config_text_isread = true;
        },
        menuTextSubmit() {
            this.menu_type_text = true;
            this.config_text_isread = true;

            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id,
                config_data: this.config_text
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_text_submit, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });
                        _this.reload();
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        },
        handleClick(tab, event) {
            this.menuDisplay();

            switch (tab.name) {
                case "configTable":
                    this.menu_type.config_table = true;
                    break;
                case "conifgText":
                    this.menu_type.conifg_text = true;
                    break;
                case "jsonText":
                    this.menu_type.json_text = true;
                    break;
                case "changeHistory":
                    this.menu_type.change_history = true;
                    break;
            }
            this.initData();
        },
        editConfigDialog(index, row) {
            this.edit = row;
            this.edit["project_id"] = this.$route.params.project_id;
            this.edit["namespace_id"] = this.$route.params.namespace_id;
            this.edit.release_status = this.edit.release_status.toString();

            this.editConfigDialogVisible = true;
        },
        configDel(index, row) {
            this.$confirm("此操作将永久删除该配置, 是否继续?", "提示", {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning"
            })
                .then(() => {
                    let data = {
                        project_id: this.$route.params.project_id,
                        namespace_id: this.$route.params.namespace_id,
                        key: row.key,
                        is_release: row.release_status.toString()
                    };
                    let _this = this;

                    this.$axios
                        .post(_this.$server.config_del, data, {})
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    message: data.msg,
                                    type: "success"
                                });

                                _this.reload();
                            } else {
                                _this.$message(data.msg);
                            }
                        })
                        .catch(function(err) {
                            _this.$message("请检查您的网络");
                        });
                })
                .catch(() => {
                    this.$message({
                        type: "info",
                        message: "已取消删除"
                    });
                });
        },
        releaseHistory() {
            this.$router.push({
                name: "releaseHistory",
                params: {
                    project_id: this.$route.params.project_id,
                    namespace_id: this.$route.params.namespace_id
                }
            });
        },
        addConfig() {
            this.isLoading = true;
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id,
                key: this.form.key,
                val: this.form.val,
                remake: this.form.remake
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_add, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });
                        _this.reload();
                    } else {
                        _this.$message(data.msg);
                    }

                    // 数据清空
                    _this.form.key = "";
                    _this.form.val = "";
                    _this.form.remake = "";
                    _this.isLoading = false;
                    _this.addConfigDialogVisible = false;
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        },
        editConfig() {
            let _this = this;
            this.isLoading = true;

            this.$axios
                .post(_this.$server.config_edit, _this.edit, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });
                        _this.reload();
                    } else {
                        _this.$message(data.msg);
                    }
                    _this.isLoading = false;
                    _this.editConfigDialogVisible = false;
                    _this.edit = {};
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        },
        releaseConfig() {
            this.$confirm("是否要发布?", "提示", {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning"
            })
                .then(() => {
                    this.isLoading = true;
                    let data = {
                        project_id: this.$route.params.project_id,
                        namespace_id: this.$route.params.namespace_id
                    };
                    let _this = this;

                    this.$axios
                        .post(_this.$server.config_release, data, {})
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    message: data.msg,
                                    type: "success"
                                });

                                _this.reload();
                            } else {
                                _this.$message(data.msg);
                            }
                            _this.isLoading = false;
                        })
                        .catch(function(err) {
                            _this.$message("请检查您的网络");
                        });
                })
                .catch(() => {
                    this.$message({
                        type: "info",
                        message: "已取消发布"
                    });
                });
        },
        getConfigLogList() {
            let data = {
                project_id: parseInt(this.$route.params.project_id),
                namespace_id: parseInt(this.$route.params.namespace_id)
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_log_list, data, {})
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        if (data.data) {
                            for (
                                let index = 0;
                                index < data.data.length;
                                index++
                            ) {
                                let temp = data.data[index];
                                temp.release_status = 0;
                                _this.config_key_val_table.push(temp);
                            }
                        }
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.getConfigList();
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        },
        getConfigList() {
            let data = {
                project_id: parseInt(this.$route.params.project_id),
                namespace_id: parseInt(this.$route.params.namespace_id)
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_list, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        if (data.data) {
                            for (
                                let index = 0;
                                index < data.data.length;
                                index++
                            ) {
                                let temp = data.data[index];
                                temp.release_status = 1;
                                _this.config_key_val_table.push(temp);
                            }
                        }
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.setConfigListText();
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        },
        setConfigListText() {
            this.config_text = "";
            let config_json_temp = [];
            let config_log_del = [];

            for (
                let index = 0;
                index < this.config_key_val_table.length;
                index++
            ) {
                if (this.config_key_val_table[index].release_status) {
                    config_json_temp.push(this.config_key_val_table[index]);
                } else {
                    if (this.config_key_val_table[index].type == -1) {
                        config_log_del.push(this.config_key_val_table[index]);
                        continue;
                    }

                    config_json_temp.push(this.config_key_val_table[index]);
                }
            }

            for (let index = 0; index < config_log_del.length; index++) {
                for (let j = 0; j < config_json_temp.length; j++) {
                    if (config_log_del[index].key == config_json_temp[j].key) {
                        config_json_temp.splice(j, 1);
                    }
                }
            }

            for (let index = 0; index < config_json_temp.length; index++) {
                this.config_text +=
                    config_json_temp[index].key +
                    " = " +
                    config_json_temp[index].val +
                    "\n";
            }
        },
        rollbackDialogVisible() {
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_rollback_last, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.rollbackData = data.data;
                        _this.centerDialogVisible = true;
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                });
        },
        setRollback() {
            this.$confirm(
                "此操作将回滚配置到上个版本，为不可逆操作, 是否继续?",
                "提示",
                {
                    confirmButtonText: "确定",
                    cancelButtonText: "取消",
                    type: "warning"
                }
            )
                .then(() => {
                    let data = {
                        project_id: this.$route.params.project_id,
                        namespace_id: this.$route.params.namespace_id
                    };
                    let _this = this;

                    this.$axios
                        .post(_this.$server.config_rollback, data, {})
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.centerDialogVisible = false;
                                _this.$message({
                                    message: data.msg,
                                    type: "success"
                                });
                                _this.reload();
                            } else {
                                _this.$message(data.msg);
                            }
                        })
                        .catch(function(err) {
                            console.log(err);
                            _this.$message("请检查您的网络");
                        });
                })
                .catch(() => {
                    this.$message({
                        type: "info",
                        message: "已取消删除"
                    });
                });
        },
        syncConfig() {
            let data = {
                project_id: this.$route.params.project_id,
                namespace_id: this.$route.params.namespace_id
            };
            let _this = this;

            this.$axios
                .post(_this.$server.config_sync, data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });
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
    inject: ["reload"],
    created() {
        this.$common.addVuexTopInfoData(this, {
            name: "配置",
            path: "config",
            params: {
                project_id: this.$route.params.id,
                namespace_id: this.$route.namespace_id
            }
        });
    },
    mounted() {
        this.getConfigLogList();
    }
};
</script>

<style>
.config-tabs-menu {
    position: absolute;
    right: 40px;
    top: 170px;
    z-index: 999;
}
.release-label {
    background: #e4e7ed;
    padding: 0.2em 0.6em 0.3em;
}
.release-not-label {
    background: rgb(255, 210, 142);
    padding: 0.2em 0.6em 0.3em;
}
</style>
