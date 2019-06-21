<template>
    <div>
        <div id="belugaLoginBg">
        </div>
        <div id="login">
            <div class="belugaLogo">
            </div>
            <div class="title">
                <strong>登录</strong>
                <p>使用邮箱登录</p>
            </div>
            <el-form :label-position="loginPosition" ref="loginPosition" :model="formData">
                <el-form-item prop="username">
                    <el-input type="text" v-model="formData.username" placeholder="请输入用户名" :disabled="input_disabled"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-input type="password" v-model="formData.passwd" placeholder="请输入密码" :disabled="input_disabled"></el-input>
                </el-form-item>
                <!-- <el-form-item>
                    <el-input type="text" v-model="formData.verification_code" placeholder="请输入图片验证码" :disabled="input_disabled" style="width: 70%;"></el-input>
                    <img id="verification_code" @click="getCaptcha" :src="verification_img">
                </el-form-item> -->

                <div style="position: relative">
                    <el-checkbox-group v-model="formData.remember_passwd" id="remember_passwd">
                        <el-checkbox label="记住密码" name="remember_passwd"></el-checkbox>
                    </el-checkbox-group>
                    <el-dropdown @command="handleCommand" trigger="click" id="more_operation">
                        <span class="el-dropdown-link">
                            更多选项
                            <i class="el-icon-arrow-down el-icon--right"></i>
                        </span>
                        <el-dropdown-menu slot="dropdown">
                            <!-- <el-dropdown-item command="/register">创建账号</el-dropdown-item> -->
                            <el-dropdown-item command="/forgotpassword">忘记密码</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </div>

                <el-form-item>
                    <el-button type="primary" @click="submitForm" class="login_button" :loading="loading">{{login_content}}</el-button>
                </el-form-item>
            </el-form>
        </div>
    </div>
</template>

<script>
export default {
    name: "login",
    data() {
        return {
            loginPosition: "left",
            formData: {
                username: "",
                passwd: "",
                verification_code: "",
                remember_passwd: false,
                captcha_id: ""
            },
            loading: false,
            login_content: "登录",
            input_disabled: false,
            verification_img: ""
        };
    },
    methods: {
        submitForm() {
            this.loading = true;
            this.login_content = "登录中";
            this.input_disabled = true;

            let _this = this;

            this.$axios.defaults.headers = {
                "Account-token": ""
            };

            this.$axios
                .post(_this.$server.login, _this.$data.formData, {})
                .then(function(response) {
                    var data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: data.msg,
                            type: "success"
                        });
                        _this.$cookie.setCookie(
                            "beluga_token",
                            data.data["token"]
                        );
                        _this.$axios.defaults.headers = {
                            "Account-token": data.data["token"]
                        };
                        _this.$router.replace("/");
                    } else {
                        _this.$message(data.msg);
                        _this.loading = false;
                        _this.login_content = "登录";
                        _this.input_disabled = false;
                        // _this.getCaptcha();
                    }
                })
                .catch(function(err) {
                    console.log(err);
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                    _this.login_content = "登录";
                    _this.input_disabled = false;
                });
        },
        handleCommand(command) {
            this.$message("click on item " + command);
        },
        getCaptcha() {
            var _this = this;

            this.$axios.defaults.headers = {
                "Account-token": ""
            };

            this.$axios
                .post(_this.$server.captcha, {
                    captcha_id: _this.$data.formData.captcha_id
                })
                .then(function(response) {
                    var data = response.data;

                    if (data.status) {
                        _this.$data.verification_img =
                            data["data"]["captcha_val"];
                        _this.$data.formData.captcha_id =
                            data["data"]["captcha_id"];
                    } else {
                        _this.$message(data.msg);
                    }
                })
                .catch(function(err) {
                    _this.$message("验证码获取失败");
                });
        }
    },
    created() {
        let token = this.$cookie.getCookie("beluga_token");

        if (
            token != "" ||
            token != null ||
            token != "undefined" ||
            token != undefined
        ) {
            this.$router.replace("/");
        }
        // this.getCaptcha();
    }
};
</script>


<style>
#belugaLoginBg {
    position: fixed;
    width: 100%;
    height: 100%;
    left: 0;
    top: 0;
    overflow: hidden;
}

#belugaLoginBg img {
    min-width: 100%;
    height: 100%;
}

#login {
    width: 300px;
    margin: 0 auto;
    height: 400px;
    background: #fff;
    padding: 50px 20px;
    position: absolute;
    top: 50%;
    margin-top: -245px;
    left: 50%;
    margin-left: -170px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
    z-index: 999;
}

#login .belugaLogo {
    text-align: center;
    position: relative;
    top: -20px;
}

#login .login_button {
    position: relative;
    top: 20px;
    width: 300px;
}

#login .title {
    text-align: center;
    padding-bottom: 20px;
}

#login .title strong {
    font-size: 24px;
    line-height: 34px;
    font-weight: normal;
}

#login .title p {
    line-height: 24px;
    font-size: 14px;
    color: #606266;
}

#login_footer {
    text-align: center;
    position: fixed;
    bottom: 10px;
    width: 100%;
    color: #fff;
    font-size: 12px;
    line-height: 44px;
}

#verification_code {
    width: 29%;
    border: 1px #dcdfe6 solid;
    height: 40px;
    border-radius: 4px;
    float: right;
}

#remember_passwd {
    display: inline;
    color: #606266;
}

#more_operation {
    position: relative;
    left: 140px;
    cursor: pointer;
}
</style>
