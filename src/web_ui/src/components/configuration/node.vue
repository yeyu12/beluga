<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-search" style="padding-bottom: 15px;">
                <el-input placeholder="请输入节点IP" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch">
                </el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="nodeData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="ip" label="IP"></el-table-column>
                <el-table-column label="状态">
                    <template slot-scope="scope">
                        <span v-if="scope.row.is_delete == 1"> 在线</span>
                        <span v-else style="color:#F56C6C;font-weight: bold;">未在线</span>
                    </template>
                </el-table-column>
                <el-table-column prop="remake" label="备注"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" @click="editNodeDialog(scope.$index, scope.row)">配置</el-button>
                        <el-button size="mini" type="danger" @click="nodeDel(scope.$index, scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>
        <node-conf :edit-node-dialog="open_edit_node_dialog" :edit-node-datas="selectd_edit_node_data"
            @hidden-edit-node-dialog="hiddenEditNodeDialog" @update-edit_node="editNodeData"></node-conf>
    </div>
</template>

<script>
import NodeConf from "@/components/configuration/node/conf";

export default {
    name: "node",
    data(){
        return{
            nodeData: [],
            page: 1,
            pageSize: 20,
            total: 0,
            total_page: 0,
            search: "",
            open_add_node_dialog: false,
            open_edit_node_dialog: false,
            selectd_edit_node_data: {}
        }
    },
    methods:{
        handleCurrentChange(val) {
            this.page = val;
        },
        nodeDel(index, row) {
            this.$confirm(
                "此操作将永久删除该节点数据, 是否继续?",
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
                            _this.$server.config_node_del,
                            { id: row.id },
                            {}
                        )
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });
                                _this.nodeData.splice(index, 1);
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
        addNodeDialog() {
            this.open_add_node_dialog = true;
        },
        editNodeDialog(index, row) {
            this.selectd_edit_node_data = row;
            this.open_edit_node_dialog = true;
        },
        editNodeData(val) {
            for (let index = 0; index < this.nodeData.length; index++) {
                if (
                    parseInt(this.nodeData[index].id) === parseInt(val.id)
                ) {
                    let data = this.nodeData[index];
                    this.$set(this.nodeData[index], "name", val.name);

                    break;
                }
            }
        },
        hiddenAddNodeDialog(val) {
            this.open_add_node_dialog = val;
        },
        addNodeData(val) {
            this.nodeData.unshift(val);
        },
        hiddenEditNodeDialog(val) {
            this.open_edit_node_dialog = val;
        },
        getNodeList() {
            let _this = this;
            let request_data = {
                page: this.page.toString()
            };

            if (this.search != "") {
                request_data["search"] = this.search;
            }

            this.$axios
                .post(_this.$server.config_node_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page;
                        _this.nodeData = data.data.list;
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
            this.getNodeList();
        }
    },
    mounted() {
        this.getNodeList();
    },
    watch: {
        page(val, oldVal) {
            this.getNodeList();
        }
    },
    components: {
        NodeConf
    }
};
</script>

<style>
</style>
