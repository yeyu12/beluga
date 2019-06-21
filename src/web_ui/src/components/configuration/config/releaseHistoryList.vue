<template>
    <div>
        <div :class="[{'active':index == selected_index}, 'release-history-main-left-list']" v-for="(val,index) in list" @click="selected(index,val)">
            <div class="release-history-main-left-list-name">{{val.account_name}}</div>
            <div class="release-history-main-left-list-type">
                <template v-if="val.operation_type == '0'">普通发布</template>
                <template v-else-if="val.operation_type == '-1'">回滚</template>
                <template v-else-if="val.operation_type == '3'">灰度操作</template>
                <template v-else>什么鬼啊</template>
            </div>
            <div class="release-history-main-left-list-date">{{val.create_time}}</div>
            <div class="release-history-main-left-list-mark" :style="backgroundMark(val.operation_type)"></div>
        </div>
    </div>
</template>

<script>
export default {
    name: "releaseHistoryList",
    data() {
        return {
            list: [],
            selected_index: 0,
            selected_val: ""
        };
    },
    props: ["data"],
    methods: {
        selected(index, val) {
            this.selected_index = index;
            this.selected_val = val;
        },

        backgroundMark(type) {
            let str = "background:";
            type = parseInt(type);

            switch (type) {
                case 0:
                    str += "#316510";
                    break;
                case -1:
                    str += "#997f1c";
                    break;
                case 3:
                    str += "#999999";
                    break;
                default:
                    str += "#000";
                    break;
            }

            return str;
        }
    },
    watch: {
        selected_val(val, oldVal) {
            val != undefined && this.$emit("selected", val);
        },
        data(val, oldVal) {
            if (val) {
                this.list = {};
                this.list = val;
                this.selected_val = this.list[this.selected_index];
            }
        }
    }
};
</script>

<style>
.active {
    background: #f5f5f5;
}
</style>
