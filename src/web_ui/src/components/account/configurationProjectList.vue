<template>
    <div>
        <div class="beluga-config-project-main" style="padding-top:0">
            <el-table :data="projectData" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="project_name" label="项目名"></el-table-column>
                <el-table-column prop="appid" label="appid"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button type="text" size="mini" @click="seeProject(scope.$index, scope.row)">查看</el-button>
                    </template>
                </el-table-column>
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
    name:"configurationProjectList",
    data(){
        return{
            page: 1,
            pageSize: 20,
            total: 0,
            total_page:0,
            projectData: []
        }
    },
    methods: {
        handleCurrentChange(val) {
            this.page = val;
        },
        getConfigList() {
            let _this = this;
            let request_data = {
                page: this.page.toString()
            }

            if(this.search != "") {
                request_data["search"] = this.search
            }

            this.$axios
                .post(
                    _this.$server.user_project_list,
                    request_data,
                    {}
                )
                .then(function(response) {
                    let data = response.data;
                    
                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page
                        _this.projectData = data.data.list;
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
        },
        seeProject(index, row){
            this.$router.push({
                name: "namespace",
                params: {
                    project_id: row.id
                }
            });
        }
    },
    mounted() {
        this.getConfigList();
    },
    watch: {
        page(val, oldVal) {
            this.getConfigList();
        }
    }
}
</script>

<style>

</style>
