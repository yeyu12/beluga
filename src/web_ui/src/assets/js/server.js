let server = {
    host: 'http://127.0.0.1:9410',

    index: '/',
    login: '/login',

    // 用户信息
    user_info: '/getUserInfo',
    set_user_info: '/setUserInfo',
    set_passwd: '/setPasswd',
    user_project_list: '/userProjectList',

    config_project_list: '/getProjectList',
    config_project_add: '/addProject',
    config_project_del: '/delProject',
    config_project_edit: '/editProject',
    config_namespace_list: '/getNamespaceList',
    config_namespace_add: '/addNamespace',
    config_namespace_edit: '/editNamespace',
    config_namespace_del: '/delNamespace',
    config_log_list: '/getConfigLogList',
    config_list: '/getConfigList',
    config_add: '/addConfig',
    config_edit: '/editConfig',
    config_del: '/delConfig',
    config_release: '/releaseConfig',
    config_rollback_last: '/getRollbackLast',
    config_rollback: '/rollback',
    config_project_name_namespace_name: '/getProjectNameToNamespaceName',
    config_release_history: '/getReleaseHistory',
    config_change_release_version_list: '/getConfigVersionReleaseList',
    config_all_release_version_list: '/getConfigAllVersionReleaseList',
    config_text_submit: '/saveConfigText',
    config_sync: '/syncConfig',

    // etcd服务
    etcd_ip_list: '/getEtcdIpList',
    etcd_ip_add: '/addEtcdNodeConf',
    etcd_ip_del: '/delEtcdNode',

    // 中心配置,节点配置
    node_conf_list: '/getNodeConfList',
    node_conf_add: '/addNodeConf',
    node_conf_edit: '/editNodeConf',
    node_conf_del: '/delNodeConf',
    node_conf_ids_list: '/getIdsNodeConf',

    // 配置节点管理
    config_node_list: "/getNodeList",
    config_node_del: "/delNode",
    config_node_etid: "/editNode",

    // 定时任务
    task_list: "/getTaskList",
    task_add: "/addTask",
    task_edit: "/editTask",
    task_run_stop: "/taskRunOrStop",
    task_del: "/taskDel",
    task_id_info: "/getTaskIdToInfo",
    task_kill: "/taskKill",
    task_node_list: "/taskNodeList",
    task_node_del: "/taskNodeDel",
    subtasks_list: "/subtasksList",
    task_log_list: "/taskLogList",

    // 系统
    // 系统用户
    user_list: '/getUserList',
    add_user: '/addUser',
    edit_user: '/editUser',

    // 操作记录
    operation_log_list: "/getOperationLogList",

    get_system_info: '/get_system_info',
};

export default server;