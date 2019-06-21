<template>
    <el-menu class="el-menu-vertical-demo" background-color="#182132" text-color="#979ba5" active-text-color="#fff" unique-opened @select="selectMenu"
        :default-active="menu_path">
        <template v-for="(val, key) in menu_data">
            <template v-if="val.data">
                <el-submenu :index="key + '-' + val.name">
                    <template slot="title">
                        <i :class="val.icon"></i>
                        <span>{{ val.name }}</span>
                    </template>
                    <template v-for="v in val.data">
                        <el-menu-item :index="v.index + '-' + v.name">{{
                            v.name
                            }}</el-menu-item>
                    </template>
                </el-submenu>
            </template>
            <template v-else>
                <el-menu-item :index="val.index + '-' + val.name">
                    <i :class="val.icon"></i>
                    <span slot="title">{{ val.name }}</span>
                </el-menu-item>
            </template>
        </template>
    </el-menu>
</template>

<script>
export default {
    name: "menus",
    data() {
        return {
            menu_data: [],
            menu_path: "",
            menu_name: ""
        };
    },
    methods: {
        selectMenu(index, indexPath) {
            let route_data = {};
            let indexPathLength = 0;
            const default_index_path_length = 2;

            indexPath != undefined && (indexPathLength = indexPath.length);

            if (indexPathLength >= default_index_path_length) {
                let nav_data = indexPath[0].split("-");
                route_data = {
                    path: "",
                    name: nav_data[1]
                };

                this.$store.commit("topInfo/emptyToAddTopinfo", route_data);
            }

            let route_name = index.split("-");
            this.menu_path = route_name[0];
            this.menu_name = route_name[1];

            route_data = {
                path: this.menu_path,
                name: this.menu_name
            };

            this.$emit("menu-selection", route_data);

            if (indexPathLength >= default_index_path_length) {
                this.$store.commit("topInfo/addTopInfo", route_data);
            } else {
                this.$store.commit("topInfo/emptyToAddTopinfo", route_data);
            }
        }
    },
    created() {
        this.menu_data = require("./../../assets/data/menu.json");
        this.menu_path = "/";
        this.menu_name = "首页";
        this.selectMenu(this.menu_path + "-" + this.menu_name);

        let default_route_data = {
            default_path: this.menu_path,
            default_name: this.menu_name,
            default_params: []
        };

        let route_data = {
            path: this.menu_path,
            name: this.menu_name,
            params: []
        };

        this.$store.commit("menu/setDefaultRoute", default_route_data);
        this.$store.commit("topInfo/emptyToAddTopinfo", route_data);
    }
};
</script>
