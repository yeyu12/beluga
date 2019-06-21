<template>
    <el-dialog title="修改项目" :visible.sync="openEditProjectDialog">
        <el-form :model="form">
            <el-form-item label="项目名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenEditProject">取 消</el-button>
            <el-button type="primary" @click="editProject" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "editPriject",
    data() {
        return {
            openEditProjectDialog: false,
            form: {
                id: "",
                name: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            projectData: {}
        };
    },
    methods: {
        hiddenEditProject() {
            this.form.name = "";
            this.openEditProjectDialog = false;
        },
        editProject() {
            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_project_edit,
                    {
                        project_name: _this.form.name,
                        project_id: _this.form.id
                    },
                    {}
                )
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: "项目修改成功",
                            type: "success"
                        });
                        // 返回新添加的数据到复层
                        _this.$emit("update-edit_project", _this.form);
                    } else {
                        _this.$message(data.msg);
                    }
                    _this.openEditProjectDialog = false;
                    _this.isLoading = false;
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        }
    },
    props: ["editProjectDialog", "editProjectDatas"],
    watch: {
        editProjectDialog(val, oldVal) {
            this.openEditProjectDialog = val;
        },
        openEditProjectDialog(val, oldVal) {
            if (val == false) {
                this.$emit("hidden-edit-project-dialog", false);
            }
        },
        editProjectDatas(val, oldVal) {
            this.form.id = val.id;
            this.form.name = val.project_name;
        }
    }
};
</script>

<style>
</style>
