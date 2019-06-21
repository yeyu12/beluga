<template>
    <el-form ref="form" :model="form" label-width="100px" id="user-info-form">
        <el-form-item label="原始密码">
            <el-input v-model="form.passwd" placeholder="请输入原始密码" type="password"></el-input>
        </el-form-item>
        <el-form-item label="新密码">
            <el-input v-model="form.new_passwd" placeholder="请输入新密码" type="password"></el-input>
        </el-form-item>
        <el-form-item label="重复新密码">
            <el-input v-model="form.repeat_new_passwd" placeholder="请重复输入新密码" type="password"></el-input>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="setPasswd">保存</el-button>
            <el-button>取消</el-button>
        </el-form-item>
    </el-form>
</template>

<script>
export default {
    name: "changePasswd",
    data() {
        return {
            form: {
                passwd: "",
                new_passwd: "",
                repeat_new_passwd: ""
            }
        };
    },
    methods: {
        setPasswd() {
            let _this = this;

            for (let index = 0; index < this.form.length; index++) {
                if(!this.form[index]) {
                    return
                }
            }

            if(this.form.passwd == this.form.new_passwd) {
                this.$message("原始密码和新密码相同");
                return
            }

            if(this.form.new_passwd != this.form.repeat_new_passwd) {
                this.$message("输入的两次新密码不同");
                return
            }

            this.$axios
                .post(_this.$server.set_passwd, _this.form, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });

                        _this.$cookie.delCookie("beluga_token");
                        _this.$router.replace("/login");
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        }
    }
};
</script>