<template>
    <el-card class="box-card" id="config-box-card">
        <div id="user-settings-top">用户设置</div>
        <div id="user-setting-menu" class="user-settings" @click="userSettingsMenuSelectd">
            <div class="user-settings-activity" style="border-top: 1px solid #eff2f7" data-components="UserInfo">个人信息</div>
            <div data-components="ChangePasswd">修改密码</div>
            <div data-components="ConfigurationProjectList">中心配置项目列表</div>
        </div>
        <div id="user-setting-info" class="user-settings">
            <component :is="currentView"></component>
        </div>
    </el-card>
</template>

<script>
import UserInfo from '@/components/account/userinfo'
import ChangePasswd from '@/components/account/changePasswd'
import ConfigurationProjectList from '@/components/account/configurationProjectList'

export default {
    name:"account",
    data(){
        return{
            currentView:"UserInfo",
        }
    },
    methods: {
        userSettingsMenuSelectd(event){
            let user_setting_dom = event.path[0];
            let compon = user_setting_dom.getAttribute("data-components");
            let user_settings_dom_class = document.getElementsByClassName("user-settings-activity");
            
            if (compon){
                this.currentView = compon;
            } else {
                return;
            }

            for (let index = 0; index < user_settings_dom_class.length; index++) {
                user_settings_dom_class[index].classList.remove("user-settings-activity");
            }

            user_setting_dom.classList.add("user-settings-activity");
        }
    },
    created() {
        let route_data = {
            name:"用户设置",
            path:"/user",
            params:{}
        };

        this.$store.commit("topInfo/emptyToAddTopinfo", route_data);
    },
    components:{
        UserInfo,
        ChangePasswd,
        ConfigurationProjectList
    }
}
</script>

<style>

</style>
