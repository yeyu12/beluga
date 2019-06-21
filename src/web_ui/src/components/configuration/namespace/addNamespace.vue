<template>
    <el-dialog title="添加命名空间" :visible.sync="openAddNamespaceDialog">
        <el-form :model="form">
            <el-form-item label="命名空间名称" :label-width="formLabelWidth">
                <el-input v-model="form.namespaceName" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenAddNamespace">取 消</el-button>
            <el-button type="primary" @click="addNamespace" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "addNamespace",
    data() {
        return {
            openAddNamespaceDialog: false,
            form: {
                namespaceName: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            namespaceData: {}
        };
    },
    methods: {
        hiddenAddNamespace() {
            this.form.name = "";
            this.openAddNamespaceDialog = false;
        },
        addNamespace() {
            this.isLoading = true;

            // 添加完成后，让页面数据刷新
            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_namespace_add,
                    {
                        namespace_name: _this.form.namespaceName,
                        project_id: _this.$route.params.project_id
                    },
                    {}
                )
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.openAddProjectDialog = false;
                        _this.$message({
                            message: "命名空间添加成功",
                            type: "success"
                        });
                        _this.isLoading = false;

                        _this.openAddNamespaceDialog = false;
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
    props: ["openNamespaceDialog"],
    watch: {
        openNamespaceDialog(val, oldVal) {
            this.openAddNamespaceDialog = val;
        },
        openAddNamespaceDialog(val, oldVal) {
            val == false && this.$emit("hidden-add-namespace-dialog", false);
        }
    }
};
</script>

<style>
</style>
