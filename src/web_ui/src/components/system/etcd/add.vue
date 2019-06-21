<template>
    <el-dialog title="添加ETCD服务" :visible.sync="openAddEtcdDialog">
        <el-form :model="form">
            <el-form-item label="ETCD-IP" :label-width="formLabelWidth">
                <el-input v-model="form.ip" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="ETCD-端口" :label-width="formLabelWidth">
                <el-input v-model="form.port" autocomplete="off"></el-input>
            </el-form-item>
        </el-form>
        <div slot="footer" class="dialog-footer">
            <el-button @click="hiddenAddEtcd">取 消</el-button>
            <el-button type="primary" @click="addEtcd" :loading="isLoading">确 定</el-button>
        </div>
    </el-dialog>
</template>

<script>
export default {
    name: "addEtcd",
    data() {
        return {
            openAddEtcdDialog: false,
            form: {
                ip: "",
                port:2379
            },
            formLabelWidth: "120px",
            isLoading: false,
            EtcdData: {}
        };
    },
    methods: {
        hiddenAddEtcd() {
            this.form.ip = "";
            this.openAddEtcdDialog = false;
            this.isLoading = false;
        },
        addEtcd() {
            this.isLoading = true;

            let _this = this;

            this.$axios
                .post(
                    _this.$server.etcd_ip_add,
                    _this.form,
                    {}
                )
                .then(function(response) {
                    let data = response.data;
                    if (data.status) {
                        _this.openAddEtcdDialog = false;
                        _this.$message({
                            message: "ETCD服务添加成功",
                            type: "success"
                        });
                        
                        // 返回新添加的数据到复层
                        _this.$emit("add-etcd-data", data.data);
                    } else { 
                        _this.$message(data.msg);
                    }

                    _this.hiddenAddEtcd()
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
            this.openAddEtcdDialog = val;
        },
        openAddEtcdDialog(val, oldVal) {
            val == false && this.$emit("hidden-add-etcd-dialog", false);
        }
    }
};
</script>

<style>
</style>
