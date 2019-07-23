<template>
    <div>
        <div class="beluga-config-project-main">
            <el-table :data="taskLogData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="task_id" label="任务ID"></el-table-column>
                <el-table-column prop="task_name" label="任务名称"></el-table-column>
                <el-table-column prop="node_ip" label="执行节点IP"></el-table-column>
                <el-table-column prop="cmd" label="执行命令"></el-table-column>
                <el-table-column prop="create_time" label="开始时间"></el-table-column>
                <el-table-column prop="end_time" label="结束时间"></el-table-column>
                <el-table-column prop="consume_time" label="执行时间"></el-table-column>
                <el-table-column label="执行状态">
                    <template slot-scope="scope">
                        <span v-if="scope.row.task_exec_type == 1"> 成功</span>
                        <span v-else style="color:#F56C6C;font-weight: bold;">失败</span>
                    </template>
                </el-table-column>
                <el-table-column prop="err" label="执行错误"></el-table-column>
                <el-table-column prop="output" label="执行输出"></el-table-column>
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
        name: "taskLog",
        data() {
            return {
                taskLogData: [],
                page: 1,
                pageSize: 20,
                total: 0,
                total_page: 0,
            }
        },
        methods: {
            handleCurrentChange(val) {
                this.page = val;
            },
            getNodeList() {
                let _this = this;
                let request_data = {
                    page: this.page.toString(),
                    task_id: this.$route.params.task_id.toString()
                };

                this.$axios
                    .post(_this.$server.task_log_list, request_data, {})
                    .then(function (response) {
                        let data = response.data;

                        if (data.status) {
                            _this.pageSize = data.data.page_size;
                            _this.total = data.data.total;
                            _this.total_page = data.data.total_page;
                            _this.taskLogData = data.data.list;
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
        },
        watch: {
            page(val, oldVal) {
                this.getNodeList();
            }
        },
        mounted() {
            this.getNodeList();
        },
    }
</script>