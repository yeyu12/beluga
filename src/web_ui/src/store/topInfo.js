const state = {
    /**
     * {
     *      name:"",
     *      path:"",
     *      params:{}
     * }
     */
    top_info: []
};

const actions = {};

const mutations = {
    // 添加信息
    addTopInfo(state, top_info) {
        if (!state.top_info.length) {
            state.top_info.push(top_info)
        } else {
            let menu_info_type = false;

            for (let i = 0; i < state.top_info.length; i++) {
                if (state.top_info[i]['path'] == top_info['path']) {
                    menu_info_type = true;
                    break;
                }
            }

            if (!menu_info_type) {
                state.top_info.push(top_info)
            }
        }

    },
    // 清空信息
    emptyTopInfo(state) {
        state.top_info = []
    },
    // 清空并且添加数据
    emptyToAddTopinfo(state, top_info) {
        mutations.emptyTopInfo(state);
        mutations.addTopInfo(state, top_info);
    },
    // 清除数组里面最后一个数据
    delTopInfo(state) {
        state.top_info.splice(state.top_info.length - 1, 1)
    }
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}