<template>
    <el-dialog title="编辑用户账号信息" :visible.sync="openEditUserDialog">
        <el-form :model="form">
            <el-form-item label="用户名" :label-width="formLabelWidth">
                <el-input v-model="form.username" autocomplete="off" disabled></el-input>
            </el-form-item>
            <el-form-item label="昵称" :label-width="formLabelWidth">
                <el-input v-model="form.nickname" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="密码" :label-width="formLabelWidth">
                <el-input v-model="form.passwd" type="password" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="激活" :label-width="formLabelWidth">
                <el-checkbox v-model="form.status">该用户激活</el-checkbox>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenEditUser">取 消</el-button>
            <el-button type="primary" @click="editUser" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "edit",
    data() {
        return {
            openEditUserDialog: false,
            form: {
                username: "",
                passwd: "",
                nickname: "",
                status: false
            },
            formLabelWidth: "120px",
            isLoading: false,
            UserData: {}
        };
    },
    methods: {
        hiddenEditUser() {
            this.form.username = "";
            this.form.passwd = "";
            this.openEditUserDialog = false;
            this.isLoading = false;
        },
        editUser() {
            this.isLoading = true;

            let _this = this;

            let data = this.form;
            data.status = this.form.status ? "1":"0";

            this.$axios
                .post(_this.$server.edit_user, _this.form, {})
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.openEditUserDialog = false;
                        _this.$message({
                            message: "用户信息修改成功",
                            type: "success"
                        });

                        // 返回新添加的数据到复层
                        _this.$emit("edit-user-data", data.data);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenEditUser();
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        }
    },
    props: ["openDialog", "editUserDatas"],
    watch: {
        openDialog(val, oldVal) {
            this.openEditUserDialog = val;
        },
        openEditUserDialog(val, oldVal) {
            val == false && this.$emit("hidden-edit-user-dialog", false);
        },
        editUserDatas(val, oldVal) {
            this.form.username = val.username;
            this.form.nickname = val.nickname;
            this.form.status = val.status == 1 ? true : false;
        }
    }
};
</script>

<style>
</style>
