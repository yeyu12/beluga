<template>
    <div id="app">
        <router-view />
    </div>
</template>

<script>
export default {
    name: "App",
    mounted() {
        let _this = this;
        // 监听浏览器返回操作
        window.addEventListener(
            "popstate",
            function(e) {
                //根据自己的需求实现自己的功能
                _this.$store.commit("topInfo/delTopInfo");
            },
            false
        );
    }
};
</script>