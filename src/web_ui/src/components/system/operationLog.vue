<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-search" style="padding-bottom: 15px;">
                <el-input placeholder="请输入操作人昵称" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch">
                </el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="operationLogList" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="c" label="操作"></el-table-column>
                <el-table-column prop="params" label="请求参数"></el-table-column>
                <el-table-column prop="nickname" label="操作人"></el-table-column>
                <el-table-column prop="ident" label="User-Agent"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>
    </div>
</template>

<script>
export default {
    name:"operationLog",
    data(){
        return{
            page: 1,
            pageSize: 20,
            total: 0,
            total_page:0,
            search:"",
            operationLogList:[],
        };
    },
    methods: {
        handleCurrentChange(val) {
            this.page = val;
        },
        getSearch(){
            this.getOperationList();
        },
        getOperationList(){
            let _this = this;
            let request_data = {
                page: this.page.toString()
            }

            if(this.search != "") {
                request_data["search"] = this.search.toString()
            }

            this.$axios
                .post(
                    _this.$server.operation_log_list,
                    request_data,
                    {}
                )
                .then(function(response) {
                    let data = response.data;
                    
                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page
                        _this.operationLogList = data.data.list;
                    } else {
                        _this.$message({
                            type: "warning",
                            message: data.msg
                        });
                    }
                })
                .catch(function(err) {
                    _this.$message("请检查您的网络");
                    _this.loading = false;
                });
        }
    },
    mounted() {
        this.getOperationList()
    },
    watch: {
        page(val, oldVal) {
            this.getOperationList();
        }
    }
}
</script>

<style>

</style>
