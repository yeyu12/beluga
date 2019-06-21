<template>
    <el-dialog title="修改命名空间" :visible.sync="openEditNamespaceDialog">
        <el-form :model="form">
            <el-form-item label="命名空间名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenEditNamespace">取 消</el-button>
            <el-button type="primary" @click="editNamespace" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "editPriject",
    data() {
        return {
            openEditNamespaceDialog: false,
            form: {
                id: "",
                name: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            namespaceData: {}
        };
    },
    methods: {
        hiddenEditNamespace() {
            this.form.name = "";
            this.form.id = "";
            this.openEditNamespaceDialog = false;
        },
        editNamespace() {
            this.isLoading = true;
            let _this = this;

            this.$axios
                .post(
                    _this.$server.config_namespace_edit,
                    {
                        namespace_id: _this.form.id.toString(),
                        namespace_name: _this.form.name
                    },
                    {}
                )
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.$message({
                            message: "命名空间修改成功",
                            type: "success"
                        });

                        _this.hiddenEditNamespace()
                    } else {
                        _this.$message(data.msg);
                    }
                    _this.openEditNamespaceDialog = false;
                    _this.isLoading = false;
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        }
    },
    props: ["editNamespaceDialog", "editNamespaceDatas"],
    watch: {
        editNamespaceDialog(val, oldVal) {
            this.openEditNamespaceDialog = val;
        },
        openEditNamespaceDialog(val, oldVal) {
            if (val == false) {
                this.$emit("hidden-edit-namespace-dialog", false);
            }
        },
        editNamespaceDatas(val, oldVal)  {
            this.form.id = val.id
            this.form.name = val.namespace_name
        }
    }
};
</script>

<style>
</style>
