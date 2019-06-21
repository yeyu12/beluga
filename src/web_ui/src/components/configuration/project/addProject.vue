<template>
    <el-dialog title="添加项目" :visible.sync="openAddProjectDialog">
        <el-form :model="form">
            <el-form-item label="项目名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenAddProject">取 消</el-button>
            <el-button type="primary" @click="addProject" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "addPriject",
    data() {
        return {
            openAddProjectDialog: false,
            form: {
                name: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            projectData: {}
        };
    },
    methods: {
        hiddenAddProject() {
            this.form.name = "";
            this.openAddProjectDialog = false;
        },
        addProject() {
            this.isLoading = true;

            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_project_add,
                    { project_name: _this.form.name },
                    {}
                )
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.openAddProjectDialog = false;
                        _this.$message({
                            message: "项目添加成功",
                            type: "success"
                        });
                        _this.isLoading = false;
                        // 返回新添加的数据到复层
                        _this.$emit("add-project-data", data.data);
                    } else {
                        _this.$message(data.msg);
                        _this.loading = false;
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        }
    },
    props: ["openDialog"],
    watch: {
        openDialog(val, oldVal) {
            this.openAddProjectDialog = val;
        },
        openAddProjectDialog(val, oldVal) {
            val == false && this.$emit("hidden-add-project-dialog", false);
        }
    }
};
</script>

<style>
</style>
