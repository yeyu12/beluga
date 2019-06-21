<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addEtcdDialog">添加etcd服务</el-button>
                <add-etcd :open-dialog="open_add_etcd_dialog" @hidden-add-etcd-dialog="hiddenAddEtcdDialog" @add-etcd-data="addEtcdData"></add-etcd>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="etcdIpData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="ip" label="etcd"></el-table-column>
                <el-table-column label="状态">
                    <template slot-scope="scope">
                        <span v-if="scope.row.status"> 正常</span>
                        <span v-else style="color:#F56C6C;font-weight: bold;">非正常</span>
                    </template>
                </el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" type="danger" @click="etcdDel(scope.$index, scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>

<script>
import AddEtcd from "@/components/system/etcd/add";

export default {
    name: "etcd",
    data() {
        return {
            etcdIpData: [],
            open_add_etcd_dialog: false
        };
    },
    methods: {
        getEtcdIpList() {
            let _this = this;

            this.$axios
                .post(_this.$server.etcd_ip_list, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.etcdIpData = data.data;
                    } else {
                        _this.$message({
                            type: "warning",
                            message: data.msg
                        });
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        },
        addEtcdDialog() {
            this.open_add_etcd_dialog = true;
        },
        hiddenAddEtcdDialog(val) {
            this.open_add_etcd_dialog = val;
        },
        addEtcdData(val) {
            this.etcdIpData.unshift(val);
        },
        etcdDel(index, row) {
            this.$confirm("此操作将ETCD服务, 是否继续?", "提示", {
                confirmButtonText: "确定",
                cancelButtonText: "取消",
                type: "warning"
            })
                .then(() => {
                    let _this = this;

                    this.$axios
                        .post(_this.$server.etcd_ip_del, row, {})
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });
                                _this.etcdIpData.splice(index, 1);
                            } else {
                                _this.$message({
                                    type: "warning",
                                    message: data.msg
                                });
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
        }
    },
    mounted() {
        this.getEtcdIpList();
    },
    components: {
        AddEtcd
        // EditProject
    }
};
</script>

<style>
</style>
