<template>
    <el-dialog title="添加中心配置节点配置" :visible.sync="openAddNodeConfDialog">
        <el-form :model="form">
            <el-form-item label="配置名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off" placeholder="请输入配置名称"></el-input>
            </el-form-item>
            <el-form-item label="配置" :label-width="formLabelWidth">
                <vue-json-editor :mode="'code'" :modes="['code']" v-model="conf_json" @has-error="jsonErr" @json-change="onJsonChange"></vue-json-editor>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenAddNodeConf">取 消</el-button>
            <el-button type="primary" @click="addNodeConf" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
import vueJsonEditor from 'vue-json-editor'

export default {
    name: "add",
    data() {
        return {
            openAddNodeConfDialog: false,
            form: {
                name: "",
                conf: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            nodeConfData: {},
            conf_json :{}
        };
    },
    methods: {
        hiddenAddNodeConf() {
            this.form.name = "";
            this.form.conf = "";
            this.openAddNodeConfDialog = false;
            this.isLoading = false;
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
                        _this.$emit("add-node-conf-data", data.data);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenAddNodeConf();
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        },
        onJsonChange(val){
            this.form.conf = JSON.stringify(val)
        },
        jsonErr(val){}
    },
    props: ["openDialog"],
    watch: {
        openDialog(val, oldVal) {
            this.openAddNodeConfDialog = val;
        },
        openAddNodeConfDialog(val, oldVal) {
            val == false && this.$emit("hidden-add-node-conf-dialog", false);
        }
    },
    components:{
        vueJsonEditor
    },
};
</script>