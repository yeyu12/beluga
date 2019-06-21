import Vue from 'vue';
import Vuex from 'vuex';
import topInfo from './topInfo'
import menu from './menu'

Vue.use(Vuex);
const store = new Vuex.Store({
    modules: {
        topInfo,
        menu
    }
});

export default store;