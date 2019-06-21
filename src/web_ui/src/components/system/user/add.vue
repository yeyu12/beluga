<template>
    <el-dialog title="添加账号" :visible.sync="openAddUserDialog">
        <el-form :model="form">
            <el-form-item label="用户名" :label-width="formLabelWidth" :rules="[{ required: true, message: '请输入用户名', trigger: 'blur' }]">
                <el-input v-model="form.username" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="密码" :label-width="formLabelWidth" :rules="[{ required: true, message: '请输入密码', trigger: 'blur' }]">
                <el-input v-model="form.passwd" type="password" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenAddUser">取 消</el-button>
            <el-button type="primary" @click="addUser" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "add",
    data() {
        return {
            openAddUserDialog: false,
            form: {
                username: "",
                passwd: "",
            },
            formLabelWidth: "120px",
            isLoading: false,
            UserData: {}
        };
    },
    methods: {
        hiddenAddUser() {
            this.form.username = "";
            this.form.passwd = "";
            this.openAddUserDialog = false;
            this.isLoading = false;
        },
        addUser() {
            this.isLoading = true;

            let _this = this;

            this.$axios
                .post(_this.$server.add_user, _this.form, {})
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.openAddUserDialog = false;
                        _this.$message({
                            message: "用户添加成功",
                            type: "success"
                        });

                        // 返回新添加的数据到复层
                        _this.$emit("add-user-data", data.data);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenAddUser();
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
            this.openAddUserDialog = val;
        },
        openAddUserDialog(val, oldVal) {
            val == false && this.$emit("hidden-add-user-dialog", false);
        }
    }
};
</script>

<style>
</style>
