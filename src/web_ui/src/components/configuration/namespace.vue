<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addNamespaceDialog">添加命名空间</el-button>
                <add-namespace :open-namespace-dialog="open_add_namespace_dialog" @hidden-add-namespace-dialog="hiddenAddNamespaceDialog"></add-namespace>
            </div>
            <div class="beluga-config-search">
                <el-input placeholder="请输入命名空间名" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch"></el-input>
            </div>
        </div>
        <div class="beluga-config-namespace-main">
            <el-table :data="namespaceData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="project_name" label="项目名"></el-table-column>
                <el-table-column prop="namespace_name" label="命名空间名"></el-table-column>
                <el-table-column prop="nickname" label="创建人"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" @click="editNamespaceDialog(scope.$index, scope.row)">编辑</el-button>
                        <el-button size="mini" @click="namespaceConfig(scope.$index, scope.row)">配置</el-button>
                        <el-button size="mini" type="danger" @click="namespaceDel(scope.$index, scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>
        <edit-namespace :edit-namespace-dialog="open_edit_namespace_dialog" :edit-namespace-datas="selectd_edit_namespace_data"
            @hidden-edit-namespace-dialog="hiddenEditNamespaceDialog"></edit-namespace>
    </div>
</template>

<script>
import AddNamespace from "@/components/configuration/namespace/addNamespace";
import EditNamespace from "@/components/configuration/namespace/editNamespace";

export default {
    name: "namespace",
    data() {
        return {
            page: 1,
            pageSize: 20,
            total: 0,
            total_page: 0,
            search: "",
            config_project_id: 0,
            open_add_namespace_dialog: false,
            open_edit_namespace_dialog: false,
            selectd_edit_namespace_data: {},
            namespaceData: []
        };
    },
    mounted() {
        this.$common.addVuexTopInfoData(this, {
            name: "命名空间",
            path: "namespace",
            params: {
                project_id: this.$route.params.project_id
            }
        });

        this.config_project_id = this.$route.params.project_id;
        this.getNamespaceList();
    },
    methods: {
        namespaceConfig(index, row) {
            this.$router.push({
                name: "config",
                params: {
                    project_id: this.$route.params.project_id,
                    namespace_id: row.id,
                    name:row.namespace_name
                }
            });
        },
        namespaceDel(index, row) {
            this.$confirm(
                "此操作将永久删除该项目并且会删除该项目下的命名空间和配置, 是否继续?",
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
                            _this.$server.config_namespace_del,
                            { namespace_id: row.id },
                            {}
                        )
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });
                                _this.namespaceData.splice(index, 1);
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
        addNamespaceDialog() {
            this.open_add_namespace_dialog = true;
        },
        editNamespaceDialog(index, row) {
            this.selectd_edit_namespace_data = row;
            this.open_edit_namespace_dialog = true;
        },
        hiddenAddNamespaceDialog(val) {
            this.open_add_namespace_dialog = val;

            if(!val) {
                this.getNamespaceList()
            }
        },
        hiddenEditNamespaceDialog(val) {
            this.open_edit_namespace_dialog = val;

            if(!val) {
                this.getNamespaceList()
            }
        },
        getNamespaceList() {
            let _this = this;
            let request_data = {
                page: this.page.toString(),
                project_id: this.config_project_id.toString()
            };

            if (this.search != "") {
                request_data["search"] = this.search;
            }

            this.$axios
                .post(_this.$server.config_namespace_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page;
                        _this.namespaceData = data.data.list;
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
            this.getNamespaceList()
        }
    },
    components: {
        AddNamespace,
        EditNamespace
    }
};
</script>

<style>
</style>
