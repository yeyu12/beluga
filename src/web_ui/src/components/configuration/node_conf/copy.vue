<template>
    <el-dialog title="添加中心配置节点配置" :visible.sync="openCopyNodeConfDialog">
        <el-form :model="form">
            <el-form-item label="配置名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off" :placeholder="name_placeholder"></el-input>
            </el-form-item>
            <el-form-item label="配置" :label-width="formLabelWidth">
                <vue-json-editor :mode="'code'" :modes="['code']" v-model="conf_json" @has-error="jsonErr" @json-change="onJsonChange"></vue-json-editor>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenCopyNodeConf">取 消</el-button>
            <el-button type="primary" @click="addNodeConf" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
import vueJsonEditor from "vue-json-editor";

export default {
    name: "copy",
    data() {
        return {
            openCopyNodeConfDialog: false,
            form: {
                name: "",
                conf: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            nodeConfData: {},
            conf_json: {},
            name_placeholder: ""
        };
    },
    methods: {
        hiddenCopyNodeConf() {
            this.form.name = "";
            this.form.conf = "";
            this.openCopyNodeConfDialog = false;
            this.isLoading = false;
            this.conf_json = {};
            this.name_placeholder = "";
        },
        addNodeConf() {
            this.isLoading = true;

            let _this = this;

            this.$axios
                .post(_this.$server.node_conf_add, _this.form, {})
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.$message({
                            message: "节点配置添加成功",
                            type: "success"
                        });

                        // 返回新添加的数据到复层
                        _this.$emit("copy_node_conf", data.data);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenCopyNodeConf();
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        },
        onJsonChange(val) {
            this.form.conf = JSON.stringify(val);
        },
        jsonErr(val) {}
    },
    props: ["copyNodeConfDialog", "copyNodeConfDatas"],
    watch: {
        copyNodeConfDialog(val, oldVal) {
            this.openCopyNodeConfDialog = val;
        },
        openCopyNodeConfDialog(val, oldVal) {
            if (val == false) {
                this.$emit("hidden-copy-node-conf-dialog", false);
                this.hiddenCopyNodeConf();
            }
        },
        copyNodeConfDatas(val, oldVal) {
            try {
                this.form.conf = val.conf;
                this.conf_json = JSON.parse(val.conf);
            } catch (e) {
                this.$notify.error({
                    title: "错误",
                    message:
                        '"' + val.name + '"' + "的配置不是json格式，请修改。"
                });
            }

            this.name_placeholder = val.name;
        }
    },
    components: {
        vueJsonEditor
    }
};
</script>

<style scoped>
</style>
