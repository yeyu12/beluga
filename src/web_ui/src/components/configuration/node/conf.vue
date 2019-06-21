<template>
    <el-dialog title="修改节点" :visible.sync="openEditNodeDialog">
        <el-form :model="form">
            <el-form-item label="节点IP" :label-width="formLabelWidth">
                <el-input v-model="form.ip" autocomplete="off" disabled></el-input>
            </el-form-item>
            <el-form-item label="备注" :label-width="formLabelWidth">
                <el-input v-model="form.remake" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="节点配置" :label-width="formLabelWidth" :rules="[{ required: true, message: '请选择配置', trigger: 'blur' }]">
                <el-select size="small" v-model="nodeConfId" multiple filterable remote reserve-keyword placeholder="请输入配置名称" :remote-method="remoteMethod" :loading="loading" clearable>
                    <el-option v-for="item in nodeConfData" :key="item.id" :label="item.name" :value="item.id"></el-option>
                </el-select>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenEditNode">取 消</el-button>
            <el-button type="primary" @click="editNode" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "conf",
    data() {
        return {
            openEditNodeDialog: false,
            form: {},
            formLabelWidth: "120px",
            isLoading: false,
            nodeData: [],
            nodeConfData:[],
            nodeConfId:[],
            loading:false,
        };
    },
    methods: {
        hiddenEditNode() {
            // this.form = {};
            this.openEditNodeDialog = false;
            this.isLoading = false;
        },
        editNode() {
            let _this = this;

            if(this.nodeConfId.length == 0) {
                this.$message("配置项不可以为空，请选择");
                return
            }

            this.$axios
                .post(
                    _this.$server.config_node_etid,
                    _this.form,
                    {}
                )
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.$message({
                            message: "节点修改成功",
                            type: "success"
                        });
                        
                        delete _this.form.node_conf_id;

                        _this.$emit("update-edit_node", _this.form);
                    } else {
                        _this.$message(data.msg);
                    }

                    _this.hiddenEditNode();
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络11111");
                    _this.loading = false;
                });
        },
        remoteMethod(search){
            if(!this.form.id) {
                this.nodeConfData = {}
                return;
            }

            this.loading = true;
            let _this = this;
            let request_data = {
                "search": search.toString()
            };

            this.$axios
                .post(_this.$server.node_conf_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.loading = false;
                        _this.nodeConfData = data.data.list;
                    } else {
                        _this.$message({
                            type: "warning",
                            message: data.msg
                        });
                        _this.loading = false;
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        },
        getNodeConfData() {
            if(this.form.node_conf_id) {
                let _this = this;
                let request_data = {
                    "id": this.form.node_conf_id.toString()
                };

                this.$axios
                    .post(_this.$server.node_conf_ids_list, request_data, {})
                    .then(function(response) {
                        let data = response.data;

                        if (data.status) {
                            _this.nodeConfData = data.data;

                            for (let index = 0; index < data.data.length; index++) {
                                _this.nodeConfId.push(data.data[index]["id"])
                            }
                        }
                    })
                    .catch(function(err) {
                        _this.$message("请检查您的网络");
                        _this.loading = false;
                    });
            }
        }
    },
    props: ["editNodeDialog", "editNodeDatas"],
    watch: {
        editNodeDialog(val, oldVal) {
            this.openEditNodeDialog = val;
        },
        openEditNodeDialog(val, oldVal) {
            (val == false) && this.$emit("hidden-edit-node-dialog", false);
        },
        editNodeDatas(val, oldVal) {
            this.form = val
            this.getNodeConfData()
        },
        nodeConfId(val, oldVal) {
            this.form.node_conf_id = val.join(",")
        }
    }
};
</script>

<style>
</style>
