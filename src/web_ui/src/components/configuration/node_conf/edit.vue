<template>
    <el-dialog title="修改节点配置" :visible.sync="openEditNodeConfDialog">
        <el-form :model="form">
            <el-form-item label="配置名称" :label-width="formLabelWidth">
                <el-input v-model="form.name" autocomplete="off" disabled></el-input>
            </el-form-item>
            <el-form-item label="配置" :label-width="formLabelWidth">
                <vue-json-editor :mode="'code'" :modes="['code']" v-model="conf_json" @has-error="jsonErr" @json-change="onJsonChange"></vue-json-editor>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenEditNodeConf">取 消</el-button>
            <el-button type="primary" @click="editNodeConf" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
import vueJsonEditor from 'vue-json-editor'

export default {
    name: "edit",
    data() {
        return {
            openEditNodeConfDialog: false,
            form: {
                id: "",
                conf: "",
                nmae: ""
            },
            formLabelWidth: "120px",
            isLoading: false,
            nodeConfData: {},
            conf_json:{}
        };
    },
    methods: {
        hiddenEditNodeConf() {
            this.openEditNodeConfDialog = false;
            this.isLoading = false;
        },
        editNodeConf() {
            let _this = this;

            this.$axios
                .post(_this.$server.node_conf_edit, _this.form, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: "节点配置修改成功",
                            type: "success"
                        });
                        // 返回新添加的数据到复层
                        _this.$emit("update-edit_node_conf", _this.form);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenEditNodeConf();
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
    props: ["editNodeConfDialog", "editNodeConfDatas"],
    watch: {
        editNodeConfDialog(val, oldVal) {
            this.openEditNodeConfDialog = val;
        },
        openEditNodeConfDialog(val, oldVal) {
            val == false && this.$emit("hidden-edit-node-conf-dialog", false);
        },
        editNodeConfDatas(val, oldVal) {
            this.form.id = val.id;
            this.form.name = val.name;
            this.form.conf = val.conf;
            this.conf_json = JSON.parse(this.form.conf)
        }
    },
    components:{
        vueJsonEditor
    },
};
</script>

<style>
</style>
