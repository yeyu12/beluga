<template>
    <el-card class="box-card" id="config-box-card">
        <div class="config-top">
            <div class="config-top-left">
                <span>修改任务</span>
            </div>
            <div class="config-top-right">
                <el-button size="mini" type="primary" @click="history_go">返回任务页</el-button>
            </div>
        </div>

        <!--        <div class="cut-off-rule"></div>-->

        <el-scrollbar style="height: calc(100vh - 230px)">
            <el-form ref="form" :model="form" label-width="170px" style="padding: 0 40px;">
                <el-divider content-position="left">基础设置</el-divider>

                <el-form-item label="任务名" :rules="[{ required: true, message: '请输入任务名', trigger: 'blur' }]"
                              prop="name">
                    <el-input v-model="form.name" placeholder="请输入任务名"></el-input>
                </el-form-item>

                <el-form :inline="true" :model="form" label-width="170px">
                    <el-form-item label="任务类型" :rules="[{ required: true, message: '请选择任务类型', trigger: 'blur' }]"
                                  prop="task_type">
                        <el-select v-model="form.task_type" placeholder="请选择任务类型" disabled>
                            <el-option label="主任务" value="1"></el-option>
<!--                            <el-option label="子任务" value="2"></el-option>-->
                        </el-select>
                    </el-form-item>

<!--                    <el-form-item label="依赖关系" v-if="displayDialog.task_type_dialog">-->
<!--                        <el-select v-model="form.rely" placeholder="请选择依赖关系">-->
<!--                            <el-option label="强依赖" value="1"></el-option>-->
<!--                            <el-option label="弱依赖" value="2"></el-option>-->
<!--                        </el-select>-->
<!--                    </el-form-item>-->

                    <!--<el-form-item label="子任务ID" v-if="displayDialog.task_type_dialog">
                        <el-select v-model="form.subtasks_id" name="subtasks" placeholder="请选子任务" clearable multiple
                                   filterable remote :remote-method="searchSubtasksNode" reserve-keyword>
                            <el-option v-for="(item, index) in subtasks_data" :key="index" :label="item.name" :value="item.id"></el-option>
                        </el-select>
                    </el-form-item>-->
                </el-form>

                <el-form :inline="true" :model="form" label-width="170px">
                    <el-form-item label="执行方式" :rules="[{ required: true, message: '请选择执行方式', trigger: 'blur' }]"
                                  prop="task_exec_type">
                        <el-select v-model="form.task_exec_type" placeholder="请选择执行方式">
                            <el-option label="shell" value="shell"></el-option>
                            <el-option label="http" value="http"></el-option>
                        </el-select>
                    </el-form-item>

                    <el-form-item label="请求方式" v-if="displayDialog.http_type_dialog">
                        <el-select v-model="form.http_type" placeholder="请选择请求方式">
                            <el-option label="post" value="post"></el-option>
                            <el-option label="get" value="get"></el-option>
                        </el-select>
                    </el-form-item>

                    <el-form-item label="节点选择">
                        <el-select v-model="form.exec_task_node_id" name="exec_task_node" placeholder="请选择节点" clearable multiple
                                   filterable remote :remote-method="searchNode">
                            <el-option label="随机" value="0" selected></el-option>
<!--                            <el-option v-for="(item, index) in task_node_data" :key="index" :label="item.ip" :value="item.id"></el-option>-->
                        </el-select>
                    </el-form-item>
                </el-form>

                <el-divider content-position="left">任务执行设置</el-divider>

                <el-form-item label="cron表达式" :rules="[{ required: true, message: '请输入cron表达式', trigger: 'blur' }]"
                              prop="cron">
                    <el-input v-model="form.cron" placeholder="请输入cron表达式" class="input-with-select">
                        <el-select v-model="select_cron" slot="prepend" placeholder="快速选择表达式">
                            <el-option label="每秒钟" value="* * * * * ? *"></el-option>
                            <el-option label="每分钟" value="0 * * * * ? *"></el-option>
                            <el-option label="每小时" value="0 0 * * * ? *"></el-option>
                            <el-option label="每天凌晨0点" value="0 0 0 * * ? *"></el-option>
                        </el-select>
                        <el-button @click="showCronDialog" slot="append">生成 cron</el-button>
                    </el-input>
                </el-form-item>

                <el-form-item label="命令" :rules="[{ required: true, message: '请输入命令或url地址', trigger: 'blur' }]"
                              prop="cmd">
                    <el-input type="textarea" v-model="form.cmd" placeholder="请输入命令或url地址" :rows="5"></el-input>
                </el-form-item>

                <el-form :inline="true" :model="form" label-width="170px">
                    <el-form-item label="任务超时时间" :rules="[{ required: true, message: '请输入任务超时时间', trigger: 'blur' }]"
                                  prop="overtime">
                        <el-input v-model="form.overtime" placeholder="请输入任务超时时间"></el-input>
                    </el-form-item>

                    <el-form-item label="任务失败重试次数"
                                  :rules="[{ required: true, message: '请输入任务失败重试次数', trigger: 'blur' }]"
                                  prop="task_fail_num">
                        <el-input v-model="form.task_fail_num" placeholder="请输入任务失败重试次数"></el-input>
                    </el-form-item>

                    <el-form-item label="任务失败重试间隔时间"
                                  :rules="[{ required: true, message: '请输入任务失败重试间隔时间', trigger: 'blur' }]"
                                  prop="task_fail_retry_time">
                        <el-input v-model="form.task_fail_retry_time" placeholder="请输入任务失败重试间隔时间"></el-input>
                    </el-form-item>
                </el-form>

                <el-form :inline="true" :model="form" label-width="170px">
                    <el-form-item label="任务通知">
                        <el-select v-model="form.task_notice" placeholder="请选择任务通知">
                            <el-option label="不通知" value="0"></el-option>
                            <el-option label="失败通知" value="1"></el-option>
                            <el-option label="总是通知" value="2"></el-option>
                            <el-option label="关键字匹配通知" value="3"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="通知类型" v-if="displayDialog.notice_type_dialog">
                        <el-select v-model="form.notice_type" placeholder="请选择通知类型">
                            <el-option label="邮件" value="1"></el-option>
                            <el-option label="webhook" value="2"></el-option>
                        </el-select>
                    </el-form-item>
                </el-form>
                <el-form-item label="任务执行输出关键字" :rules="[{ required: true, message: '请输入任务执行输出关键字', trigger: 'blur' }]"
                              prop="keyword_notice" v-if="displayDialog.keyword_notice_dialog">
                    <el-input v-model="form.keyword_notice" placeholder="请输入任务执行输出关键字"></el-input>
                </el-form-item>

                <el-form-item label="备注">
                    <el-input type="textarea" v-model="form.remake" placeholder="请输入备注"></el-input>
                </el-form-item>

                <el-form-item style="text-align: center">
                    <el-button type="primary" @click="saveTask">保存</el-button>
                    <el-button @click="history_go">取消</el-button>
                </el-form-item>
            </el-form>
        </el-scrollbar>

        <el-dialog title="生成 cron" :visible.sync="showCron" style="padding: 0;">
            <vcrontab @hide="showCron=false" @fill="crontabFill" :expression="expression"></vcrontab>
        </el-dialog>
    </el-card>
</template>

<script>
    import vcrontab from "vcrontab";

    export default {
        name: "editTask",
        data() {
            return {
                form: {
                    name: "",
                    task_type: "1",
                    rely: "1",
                    subtasks_id: [],
                    task_exec_type: "shell",
                    http_type: "post",
                    exec_task_node_id: [],
                    cmd: "",
                    cron: "",
                    overtime: 0,
                    task_fail_num: 0,
                    task_fail_retry_time: 0,
                    task_notice: "0",
                    notice_type: "1",
                    keyword_notice: "",
                    remake: "",
                },
                displayDialog: {
                    task_type_dialog: true,
                    http_type_dialog: false,
                    notice_type_dialog: false,
                    keyword_notice_dialog: false,
                },
                expression: "",
                showCron: false,
                select_cron: "",
                task_node_data: [],
                subtasks_data: [],
            };
        },
        methods: {
            history_go() {
                this.$router.go(-1);
            },
            saveTask() {
                let _this = this;

                this.$axios
                    .post(_this.$server.task_edit, this.form, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.$message({
                                message: data.msg,
                                type: "success"
                            });

                            _this.history_go();
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            crontabFill(value) {
                this.form.cron = value;
            },
            showCronDialog() {
                this.expression = this.form.cron;
                this.showCron = true;
            },
            getTaskInfo() {
                let _this = this;

                this.$axios
                    .post(_this.$server.task_id_info, {task_id: parseInt(this.$route.params.task_id)}, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            data.data.task_type = data.data.task_type.toString();
                            data.data.rely = data.data.rely.toString();
                            data.data.task_notice = data.data.task_notice.toString();
                            data.data.notice_type = data.data.notice_type.toString();
                            data.data.exec_task_node_id = data.data.exec_task_node_id.split(',');

                            if (data.data.subtasks_id) {
                                data.data.subtasks_id = data.data.subtasks_id.split(',');
                            }

                            _this.form = data.data;
                            _this.select_cron = _this.form.cron;
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            searchNode(search) {
                let _this = this;

                this.$axios
                    .post(_this.$server.task_node_list, {"search": search}, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.task_node_data = data.data.list
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            getTaskNode() {
                let _this = this;

                this.$axios
                    .post(_this.$server.task_node_list, {}, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.task_node_data = data.data.list
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            getSubtasksList() {
                let _this = this;

                this.$axios
                    .post(_this.$server.subtasks_list, {}, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.subtasks_data = data.data.list;
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            searchSubtasksNode(search) {
                let _this = this;

                this.$axios
                    .post(_this.$server.subtasks_list, {"search": search}, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.subtasks_data = data.data.list
                        } else {
                            _this.$message(data.msg);
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            }
        },
        mounted() {
            this.$common.addVuexTopInfoData(this, {
                name: "修改任务",
                path: "taskEdit"
            });

            this.getTaskInfo();
            this.getTaskNode();
            this.getSubtasksList();
        },
        components: {
            vcrontab
        },
        watch: {
            form: {
                handler(val, oldVal) {
                    if (val.task_exec_type == "shell") {
                        this.displayDialog.http_type_dialog = false
                    } else {
                        this.displayDialog.http_type_dialog = true
                    }

                    if (val.task_type == "1") {
                        this.displayDialog.task_type_dialog = true
                    } else {
                        this.displayDialog.task_type_dialog = false
                    }

                    if (val.task_notice == "0") {
                        this.displayDialog.notice_type_dialog = false
                    } else {
                        this.displayDialog.notice_type_dialog = true
                    }

                    if (val.task_notice == "3") {
                        this.displayDialog.keyword_notice_dialog = true
                    } else {
                        this.displayDialog.keyword_notice_dialog = false
                    }

                    let exec_task_node_id_len = val.exec_task_node_id.length;
                    for (let i = 0; i < exec_task_node_id_len; i++) {
                        if (exec_task_node_id_len == 1) break;
                        if (val.exec_task_node_id[i] == "0") {
                            this.form.exec_task_node_id.splice(0, exec_task_node_id_len);
                            this.form.exec_task_node_id = ["0"];
                            break
                        }
                    }
                },
                deep: true,
            },
            select_cron(val, oldVal) {
                this.form.cron = val;
            }
        }
    };
</script>