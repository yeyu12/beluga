const state = {
    default_path: "",
    default_name: "",
    default_params: []
};

const actions = {};

const mutations = {
    // 更新默认路由
    setDefaultRoute(state, route_data) {
        state.default_path = route_data.default_path;
        state.default_name = route_data.default_name;
        state.default_params = route_data.default_params;
    }
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}