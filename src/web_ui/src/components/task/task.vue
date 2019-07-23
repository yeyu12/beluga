<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addTask">添加任务</el-button>
            </div>
            <div class="beluga-config-search" style="width:30vw">
                <el-input placeholder="请输入任务ID" suffix-icon="el-icon-search" clearable v-model="task_id"
                          @keyup.enter.native="getSearch"
                          style="width:14vw"></el-input>
                <el-input placeholder="请输入任务名" suffix-icon="el-icon-search" clearable v-model="search"
                          @keyup.enter.native="getSearch"
                          style="width:14vw;float: right;"></el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="tasktData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="name" label="任务名"></el-table-column>
                <el-table-column prop="cron" label="cron表达式"></el-table-column>
                <el-table-column prop="next_exec_time" label="下次执行时间"></el-table-column>
                <el-table-column prop="task_exec_type" label="执行方式"></el-table-column>
                <el-table-column label="状态">
                    <template slot-scope="scope">
                        <span v-if="scope.row.status == 1"> run</span>
                        <span v-else style="color:red;">stop</span>
                    </template>
                </el-table-column>
                <el-table-column prop="nickname" label="创建人"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作" fixed="right" width="120">
                    <template slot-scope="scope">
                        <el-button v-if="scope.row.status == 1" type="text" style="font-size:14px;"
                                   @click="taskStop(scope.$index, scope.row)">停止
                        </el-button>
                        <el-button v-else type="text" style="font-size:14px; color:red"
                                   @click="taskRun(scope.$index, scope.row)">启动
                        </el-button>

                        <el-dropdown trigger="click">
                            <span class="el-dropdown-link" style="font-size: 14px;color: #409EFF">
                                操作<i class="el-icon-arrow-down el-icon--right"></i>
                            </span>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item @click.native="editTask(scope.row)">修改任务</el-dropdown-item>
                                <el-dropdown-item @click.native="delTask(scope.row)">删除任务</el-dropdown-item>
<!--                                <el-dropdown-item @click.native="execTask(scope.row)">手动执行</el-dropdown-item>-->
                                <el-dropdown-item @click.native="killTask(scope.row)">强杀任务</el-dropdown-item>
                                <el-dropdown-item @click.native="logTask(scope.row)">任务日志</el-dropdown-item>
                            </el-dropdown-menu>
                        </el-dropdown>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize"
                           layout="total, prev, pager, next, jumper"
                           :total="total" background>
            </el-pagination>
        </div>
    </div>
</template>

<script>

    export default {
        name: "task",
        data() {
            return {
                page: 1,
                pageSize: 20,
                total: 0,
                total_page: 0,
                search: "",
                tasktData: [],
                task_id: ""
            };
        },
        methods: {
            editTask(row) {
                this.$router.push({
                    name: "taskEdit",
                    params: {
                        task_id: row.id,
                    }
                });
            },
            delTask(row) {
                this.$confirm(
                    "是否删除该任务?确认后将在执行完当前任务后删除",
                    "提示",
                    {
                        confirmButtonText: "确定",
                        cancelButtonText: "取消",
                        type: "warning"
                    }
                )
                    .then(() => {
                        let _this = this;

                        this.$axios
                            .post(
                                _this.$server.task_del,
                                {task_id: parseInt(row.id), status: 1},
                                {}
                            )
                            .then(function (response) {
                                let data = response.data;

                                if (data.status) {
                                    _this.$message({
                                        type: "success",
                                        message: data.msg
                                    });
                                    _this.reload();
                                } else {
                                    _this.$message({
                                        type: "warning",
                                        message: data.msg
                                    });
                                }
                            })
                            .catch(function (err) {
                                _this.$message("请检查您的网络");
                            });
                    })
                    .catch(() => {
                        this.$message({
                            type: "info",
                            message: "取消删除"
                        });
                    });
            },
            execTask(row) {
            },
            killTask(row) {
                let _this = this;

                this.$confirm(
                    "是否要强杀该任务?确认后将正在执行的任务将被中断执行。",
                    "提示",
                    {
                        confirmButtonText: "确定",
                        cancelButtonText: "取消",
                        type: "warning"
                    }
                ).then(() => {
                    this.$axios
                        .post(
                            _this.$server.task_kill,
                            {task_id: parseInt(row.id)},
                            {}
                        )
                        .then(function (response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });

                                _this.reload();
                            } else {
                                _this.$message({
                                    type: "warning",
                                    message: data.msg
                                });
                            }
                        })
                        .catch(function (err) {
                            _this.$message("请检查您的网络");
                        });
                })
                    .catch(() => {
                        this.$message({
                            type: "info",
                            message: "取消删除"
                        });
                    });
            },
            logTask(row) {
                this.$router.push({
                    name: "taskLog",
                    params: {
                        task_id: row.id,
                    }
                });
            },
            handleCurrentChange(val) {
                this.page = val;
            },
            taskStop(index, row) {
                this.$confirm(
                    "是否停止该任务?确认后将在执行完当前任务后停止",
                    "提示",
                    {
                        confirmButtonText: "确定",
                        cancelButtonText: "取消",
                        type: "warning"
                    }
                )
                    .then(() => {
                        let _this = this;

                        this.$axios
                            .post(
                                _this.$server.task_run_stop,
                                {task_id: parseInt(row.id), status: 2},
                                {}
                            )
                            .then(function (response) {
                                let data = response.data;

                                if (data.status) {
                                    _this.$message({
                                        type: "success",
                                        message: data.msg
                                    });

                                    _this.reload();
                                } else {
                                    _this.$message({
                                        type: "warning",
                                        message: data.msg
                                    });
                                }
                            })
                            .catch(function (err) {
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
            taskRun(index, row) {
                let _this = this;

                this.$axios
                    .post(
                        _this.$server.task_run_stop,
                        {task_id: parseInt(row.id), status: 1},
                        {}
                    )
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.$message({
                                type: "success",
                                message: data.msg
                            });
                            _this.reload();
                        } else {
                            _this.$message({
                                type: "warning",
                                message: data.msg
                            });
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                    });
            },
            addTask() {
                this.$router.push({
                    name: "taskAdd"
                });
            },
            getTaskList() {
                let _this = this;
                let request_data = {
                    page: this.page.toString()
                };

                if (this.search != "") {
                    request_data["search"] = this.search;
                }

                if (this.task_id != "") {
                    request_data["task_id"] = this.task_id;
                }

                this.$axios
                    .post(_this.$server.task_list, request_data, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.pageSize = data.data.page_size;
                            _this.total = data.data.total;
                            _this.total_page = data.data.total_page;
                            _this.tasktData = data.data.list;
                        } else {
                            _this.$message({
                                type: "warning",
                                message: data.msg
                            });
                        }
                    })
                    .catch(function (err) {
                        _this.$message("请检查您的网络");
                        _this.loading = false;
                    });
            },
            getSearch() {
                this.getTaskList();
            }
        },
        components: {},
        mounted() {
            this.getTaskList();
        },
        watch: {
            page(val, oldVal) {
                this.getTaskList();
            }
        },
        inject: ["reload"]
    };
</script>