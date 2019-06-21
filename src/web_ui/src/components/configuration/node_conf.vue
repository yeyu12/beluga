<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addNodeConfDialog">添加节点配置</el-button>
                <add-node-conf :open-dialog="open_add_node_conf_dialog" @hidden-add-node-conf-dialog="hiddenAddNodeConfDialog"
                    @add-node-conf-data="addNodeConfData"></add-node-conf>
            </div>
            <div class="beluga-config-search">
                <el-input placeholder="请输入配置名" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch">
                </el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="nodeConfData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="name" label="配置名"></el-table-column>
                <el-table-column prop="nickname" label="创建人"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" @click="editNodeConfDialog(scope.$index, scope.row)">编辑</el-button>
                        <el-button size="mini" type="primary" @click="copyNodeConfDialog(scope.$index, scope.row)">复制</el-button>
                        <el-button size="mini" type="danger" @click="nodeConfDel(scope.$index, scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>
        <edit-node-conf :edit-node-conf-dialog="open_edit_node_conf_dialog" :edit-node-conf-datas="selectd_edit_node_conf_data"
            @hidden-edit-node-conf-dialog="hiddenEditNodeConfDialog" @update-edit_node_conf="editNodeConfData"></edit-node-conf>

        <copy-node-conf :copy-node-conf-dialog="open_copy_node_conf_dialog" :copy-node-conf-datas="copy_node_data"
            @hidden-copy-node-conf-dialog="hiddenCopyNodeConfDialog" @copy_node_conf="copyNodeConfData"></copy-node-conf>

    </div>
</template>

<script>
import AddNodeConf from "@/components/configuration/node_conf/add";
import EditNodeConf from "@/components/configuration/node_conf/edit";
import CopyNodeConf from "@/components/configuration/node_conf/copy";

export default {
    name: "node_conf",
    data() {
        return {
            nodeConfData: [],
            page: 1,
            pageSize: 20,
            total: 0,
            total_page: 0,
            search: "",
            open_add_node_conf_dialog: false,
            open_edit_node_conf_dialog: false,
            open_copy_node_conf_dialog: false,
            selectd_edit_node_conf_data: {},
            copy_node_data: {}
        };
    },
    methods: {
        handleCurrentChange(val) {
            this.page = val;
        },
        nodeConfDel(index, row) {
            this.$confirm(
                "此操作将永久删除该配置和清除该节点下的配置, 是否继续?",
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
                        .post(_this.$server.node_conf_del, { id: row.id }, {})
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });
                                _this.nodeConfData.splice(index, 1);
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
        },
        addNodeConfDialog() {
            this.open_add_node_conf_dialog = true;
        },
        editNodeConfDialog(index, row) {
            this.selectd_edit_node_conf_data = row;
            this.open_edit_node_conf_dialog = true;
        },
        editNodeConfData(val) {
            for (let index = 0; index < this.nodeConfData.length; index++) {
                if (
                    parseInt(this.nodeConfData[index].id) === parseInt(val.id)
                ) {
                    let data = this.nodeConfData[index];
                    this.$set(this.nodeConfData[index], "name", val.name);
                    this.$set(this.nodeConfData[index], "conf", val.conf);

                    break;
                }
            }
        },
        hiddenAddNodeConfDialog(val) {
            this.open_add_node_conf_dialog = val;
        },
        addNodeConfData(val) {
            this.nodeConfData.unshift(val);
        },
        copyNodeConfData(val) {
            this.nodeConfData.unshift(val);
        },
        hiddenEditNodeConfDialog(val) {
            this.open_edit_node_conf_dialog = val;
        },
        hiddenCopyNodeConfDialog(val) {
            this.open_copy_node_conf_dialog = val;
        },
        getNodeConfList() {
            let _this = this;
            let request_data = {
                page: this.page.toString()
            };

            if (this.search != "") {
                request_data["search"] = this.search;
            }

            this.$axios
                .post(_this.$server.node_conf_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page;
                        _this.nodeConfData = data.data.list;
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
        getSearch() {
            this.getNodeConfList();
        },
        copyNodeConfDialog(index, row) {
            this.copy_node_data = row;
            this.open_copy_node_conf_dialog = true;
        }
    },
    mounted() {
        this.getNodeConfList();
    },
    watch: {
        page(val, oldVal) {
            this.getNodeConfList();
        }
    },
    components: {
        AddNodeConf,
        EditNodeConf,
        CopyNodeConf
    }
};
</script>

<style>
</style>
