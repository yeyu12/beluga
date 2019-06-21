<template>
    <el-form ref="form" :model="form" label-width="80px" id="user-info-form">
        <el-form-item label="用户名">
            <el-input v-model="form.username"></el-input>
        </el-form-item>
        <el-form-item label="昵称">
            <el-input v-model="form.nickname"></el-input>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="setUserinfo">保存</el-button>
            <el-button>取消</el-button>
        </el-form-item>
    </el-form>
</template>

<script>
export default {
    name: "userInfo",
    data() {
        return {
            form: {}
        };
    },
    methods: {
        getUserInfo() {
            let _this = this;

            this.$axios
                .post(_this.$server.user_info, {}, {})
                .then(function(response) {
                    var data = response.data;

                    if (data.status) {
                        _this.form = data.data;
                    } else {
                        _this.$message("获取用户信息失败");
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        },
        setUserinfo(){
            let _this = this;

            this.$axios
                .post(_this.$server.set_user_info, _this.form, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });

                        _this.form = data.data
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                });
        }
    },
    created() {
        this.getUserInfo();
    }
};
</script>