<template>
    <div>
        <div class="beluga-config-top">
            <div class="beluga-config-top-add">
                <el-button type="primary" @click="addUserDialog">添加账号</el-button>
                <add-user :open-dialog="open_add_user_dialog" @hidden-add-user-dialog="hiddenAddUserDialog" @add-user-data="addUserData"></add-user>
            </div>
            <div class="beluga-config-search" style="padding-bottom: 15px;">
                <el-input placeholder="请输入操作人昵称" suffix-icon="el-icon-search" v-model="search" @keyup.enter.native="getSearch">
                </el-input>
            </div>
        </div>
        <div class="beluga-config-project-main">
            <el-table :data="userList" height="70vh" border style="width: 100%" stripe>
                <el-table-column prop="id" label="id"></el-table-column>
                <el-table-column prop="username" label="用户名"></el-table-column>
                <el-table-column prop="nickname" label="用户昵称"></el-table-column>
                <el-table-column label="已激活">
                    <template slot-scope="scope">
                        <template v-if="scope.row.status == 1">激活</template>
                        <template v-else><span style="color:red">未激活</span></template>
                    </template>
                </el-table-column>
                <el-table-column prop="configuration_num" label="配置项目数"></el-table-column>
                <el-table-column prop="create_time" label="创建时间"></el-table-column>
                <el-table-column label="操作">
                    <template slot-scope="scope">
                        <el-button size="mini" @click="editUserDialog(scope.$index, scope.row)">编辑</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <div class="beluga-page" v-if="total_page > 1">
            <el-pagination @current-change="handleCurrentChange" :current-page="page" :page-size="pageSize" layout="total, prev, pager, next, jumper"
                :total="total" background>
            </el-pagination>
        </div>

        <edit-user :edit-user-datas="selectd_edit_user_data" :open-dialog="open_edit_user_dialog" @hidden-edit-user-dialog="hiddenEditUserDialog" @edit-user-data="editUserData"></edit-user>
    </div>
</template>

<script>
import AddUser from "@/components/system/user/add";
import EditUser from "@/components/system/user/edit";

export default {
    name: "userlist",
    data() {
        return {
            page: 1,
            pageSize: 20,
            total: 0,
            total_page: 0,
            search: "",
            userList: [],
            open_add_user_dialog: false,
			open_edit_user_dialog: false,
			selectd_edit_user_data:{},
        };
    },
    methods: {
        getSearch() {
            this.getUserList();
        },
        getUserList() {
            let _this = this;
            let request_data = {
                page: this.page.toString()
            };

            if (this.search != "") {
                request_data["search"] = this.search;
            }

            this.$axios
                .post(_this.$server.user_list, request_data, {})
                .then(function(response) {
                    let data = response.data;

                    if (data.status) {
                        _this.pageSize = data.data.page_size;
                        _this.total = data.data.total;
                        _this.total_page = data.data.total_page;
                        _this.userList = data.data.list;
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
		editUserDialog(index, row) {
			this.selectd_edit_user_data = row;
            this.open_edit_user_dialog = true;
        },
        addUserDialog() {
            this.open_add_user_dialog = true;
        },
        hiddenAddUserDialog(val) {
            this.open_add_user_dialog = val;
        },
        hiddenEditUserDialog(val) {
            this.open_edit_user_dialog = val;
        },
        addUserData(val) {
            this.reload();
		},
		editUserData(val){
			this.reload();
		}
    },
    watch: {
        page(val, oldVal) {
            this.getUserList();
        }
    },
    mounted() {
        this.getUserList();
    },
    components: {
		AddUser,
		EditUser,
    },
    inject: ["reload"]
};
</script>

<style scoped>
</style>
