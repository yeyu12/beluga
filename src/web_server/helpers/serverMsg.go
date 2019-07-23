package helpers

const (
	SUCCESS_MSG = "成功"
	FAIL_MSG    = "失败"

	ILLEGAL_FAIL_MSG = "非法请求"
	FAIL_PARAMS_MSG  = "缺少参数"
	CAPTCHA_FAIL_MSG = "验证码错误或超时"

	// 登录相关
	LOGIN_FATL_MSG              = "您还未登录，请先登录"
	USERNAME_PASSWD_ERROR_MSG   = "账号或密码错误"
	REPEAT_LOGIN_FAIL_MSG       = "你已登录，请勿重复登录"
	TOKEN_FATL_MSG              = "token非法"
	TOKEN_FATL_OVERDUE_MSG      = "token过期"
	LOGIN_TOKEN_UPDATE_FATL_MSG = "登录失败，token修改失败"
	USER_PROHIBIT_MSG           = "账号被禁止..."

	// 配置中心相关
	CONFIGURATION_PROJECT_NAME_FAIL_MSG          = "项目名称不能为空"
	CONFIGURATION_ADD_PROJECT_FAIL_MSG           = "项目添加失败"
	CONFIGURATION_EDIT_PROJECT_FAIL_MSG          = "项目修改失败"
	CONFIGURATION_DEL_PROJECT_FAIL_MSG           = "项目删除失败"
	CONFIGURATION_NAMESPACE_NAME_FAIL_MSG        = "命名空间名称不能为空"
	CONFIGURATION_NAMESPACE_EDIT_FAIL_MSG        = "命名空间修改失败"
	CONFIGURATION_DEL_NAMESPACE_FAIL_MSG         = "命名空间删除失败"
	CONFIGURATION_PROJECT_NAME_REPEAT_FAIL_MSG   = "项目名称不可重复"
	CONFIGURATION_NAMEPSACE_NAME_REPEAT_FAIL_MSG = "命名空间名称不可重复"

	ETCD_SERVER_EXISTENCE_MSG      = "ETCD服务地址存在"
	ETCD_SERVER_NOTE_EXISTENCE_MSG = "ETCD服务地址不存在"
	ETCD_SERVER_INIT_MSG           = "ETCD服务初始化失败,请检查服务"
	ETCD_SERVER_NOTE_MSG           = "ETCD服务不存在"

	// 配置中心节点配置
	CONFIGURATION_NODE_CONF_NAME_FAIL_MSG   = "节点配置名称不能为空"
	CONFIGURATION_DEL_NODE_CONF_FAIL_MSG    = "节点配置删除失败"
	CONFIGURATION_EDIT_NODE_CONF_FAIL_MSG   = "节点配置修改失败"
	CONFIGURATION_NODE_CONF_ONLINE_NODE_MSG = "节点配置有在线配置在使用，禁止删除"

	// 配置中心节点
	CONFIGURATION_NODE_ONLINE_DEL_FATL_MSG = "节点在线禁止删除"

	// 修改密码
	CHANGE_PASSWD_IDENTICAL_FAIL_MSG = "原始密码和新密码相同"
	CHANGE_PASSWD_EMPTY_FAIL_MSG     = "密码不能为空"
	CHANGE_PASSWD_FAIL_MSG           = "原始密码错误"

	// 用户相关
	USERNAME_EXISTENCE_MSG      = "用户名已存在"
	USERNAME_TO_PASSWD_FAIL_MSG = "用户名或密码不能为空"

	// 任务相关
	TASK_CRONTAB_ANALYSIS_FAIL_MSG = "cron表达式解析失败"
	TASK_NAME_FAIL_MSG             = "任务名不能为空"
	TASK_CMD_FAIL_MSG              = "任务命令不能为空"
	TASK_OVERTIME_FAIL_MSG         = "任务超时时间非法，不能小于0与不能大于86400秒"
	TASK_FAIL_NUM_MSG              = "任务失败重试次数不能小于0"
	TASK_FAIL_RETRY_TIME_MSG       = "任务失败重试间隔时间不能小于0"
	TASK_NOTICE_FAIL_MSG           = "任务通知类型非法"
	TASK_KEYWORD_NOTICE_FAIL_MSG   = "任务执行输出关键字不能为空"
	TASK_NODE_ONLINE_DEL_FATL_MSG  = "节点在线禁止删除"
)
