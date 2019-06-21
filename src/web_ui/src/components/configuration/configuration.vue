<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addProjectDialog">添加项目</el-button>
                <add-project :open-dialog="open_add_project_dialog" @hidden-add-project-dialog="hiddenAddProjectDialog" @add-project-data="addProjectData"></add-project>
            </div>
            <div class="beluga-config-search">
                <el-input placeholder="请输入项目名" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch">
                </el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="projectData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="project_name" label="项目名"></el-table-column>
                <el-table-column prop="appid" label="appid"></el-table-column>
                <el-table-column prop="nickname" label="创建人"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" @click="editProjectDialog(scope.$index, scope.row)">编辑</el-button>
                        <el-button size="mini" @click="projectNamespace(scope.$index, scope.row)">命名空间</el-button>
                        <el-button size="mini" type="danger" @click="projectDel(scope.$index, scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>
        <edit-project :edit-project-dialog="open_edit_project_dialog" :edit-project-datas="selectd_edit_project_data"
            @hidden-edit-project-dialog="hiddenEditrojectDialog" @update-edit_project="editProjectData"></edit-project>
    </div>
</template>

<script>
import AddProject from "@/components/configuration/project/addProject";
import EditProject from "@/components/configuration/project/editProject";

export default {
    name: "configuration",
    data() {
        return {
            page: 1,
            pageSize: 20,
            total: 0,
            total_page: 0,
            search: "",
            open_add_project_dialog: false,
            open_edit_project_dialog: false,
            selectd_edit_project_data: {},
            projectData: []
        };
    },
    methods: {
        handleCurrentChange(val) {
            this.page = val;
        },
        projectNamespace(index, row) {
            this.$router.push({
                name: "namespace",
                params: {
                    project_id: row.id,
                    project_name: row.project_name
                }
            });
        },
        projectDel(index, row) {
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
                            _this.$server.config_project_del,
                            { project_id: row.id },
                            {}
                        )
                        .then(function(response) {
                            let data = response.data;

                            if (data.status) {
                                _this.$message({
                                    type: "success",
                                    message: data.msg
                                });
                                _this.projectData.splice(index, 1);
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
        addProjectDialog() {
            this.open_add_project_dialog = true;
        },
        editProjectDialog(index, row) {
            this.selectd_edit_project_data = row;
            this.open_edit_project_dialog = true;
        },
        editProjectData(val) {
            for (let index = 0; index < this.projectData.length; index++) {
                if (parseInt(this.projectData[index].id) === parseInt(val.id)) {
                    let data = this.projectData[index];
                    this.$set(
                        this.projectData[index],
                        "project_name",
                        val.name
                    );

                    break;
                }
            }
        },
        hiddenAddProjectDialog(val) {
            this.open_add_project_dialog = val;
        },
        addProjectData(val) {
            this.projectData.unshift(val);
        },
        hiddenEditrojectDialog(val) {
            this.open_edit_project_dialog = val;
        },
        getConfigList() {
            let _this = this;
            let request_data = {
                page: this.page.toString()
            };

            if (this.search != "") {
                request_data["search"] = this.search;
            }

            this.$axios
                .post(_this.$server.config_project_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page;
                        _this.projectData = data.data.list;
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
            this.getConfigList();
        }
    },
    components: {
        AddProject,
        EditProject
    },
    mounted() {
        this.getConfigList();
    },
    watch: {
        page(val, oldVal) {
            this.getConfigList();
        }
    }
};
</script>
