package routers

import (
	"beluga/src/web_server/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 登录
	beego.Router("/login", &controllers.LoginController{}, "*:Login")
	beego.Router("/logout", &controllers.LoginController{}, "*:Logout")

	// 用户信息
	beego.Router("/getUserInfo", &controllers.AccountController{}, "*:GetUserInfo")
	beego.Router("/setUserInfo", &controllers.AccountController{}, "*:SetUserInfo")
	beego.Router("/setPasswd", &controllers.AccountController{}, "*:ChangePasswd")
	beego.Router("/userProjectList", &controllers.ConfigProjectController{}, "*:GetUserProjectList")

	// 配置
	beego.Router("/getProjectList", &controllers.ConfigProjectController{}, "*:GetProjectList")
	beego.Router("/addProject", &controllers.ConfigProjectController{}, "*:AddProject")
	beego.Router("/delProject", &controllers.ConfigProjectController{}, "*:DelProject")
	beego.Router("/editProject", &controllers.ConfigProjectController{}, "*:EditProject")
	beego.Router("/getNamespaceList", &controllers.NamespaceController{}, "*:GetNamespaceList")
	beego.Router("/addNamespace", &controllers.NamespaceController{}, "*:AddNamespace")
	beego.Router("/editNamespace", &controllers.NamespaceController{}, "*:EditNamespace")
	beego.Router("/delNamespace", &controllers.NamespaceController{}, "*:DelNamespace")
	beego.Router("/getConfigLogList", &controllers.ConfigController{}, "*:GetConfigLogList")
	beego.Router("/getConfigList", &controllers.ConfigController{}, "*:GetConfigList")
	beego.Router("/addConfig", &controllers.ConfigController{}, "*:AddConfig")
	beego.Router("/delConfig", &controllers.ConfigController{}, "*:DelConfig")
	beego.Router("/editConfig", &controllers.ConfigController{}, "*:EditConfig")
	beego.Router("/releaseConfig", &controllers.ConfigController{}, "*:ReleaseConfig")
	beego.Router("/getRollbackLast", &controllers.ConfigController{}, "*:GetRollbackLast")
	beego.Router("/rollback", &controllers.ConfigController{}, "*:RollbackConfig")
	beego.Router("/getReleaseHistory", &controllers.ConfigController{}, "*:ReleaseHistory")
	beego.Router("/getProjectNameToNamespaceName", &controllers.ConfigController{}, "*:GetProjectNameToNamespaceName")
	beego.Router("/getConfigVersionReleaseList", &controllers.ConfigController{}, "*:GetConfigVersionReleaseList")
	beego.Router("/getConfigAllVersionReleaseList", &controllers.ConfigController{}, "*:GetConfigAllVersionReleaseList")
	beego.Router("/saveConfigText", &controllers.ConfigController{}, "*:SaveConfigText")
	beego.Router("/syncConfig", &controllers.ConfigController{}, "*:SyncConfig")

	// 中心配置,节点配置
	beego.Router("/getNodeConfList", &controllers.ConfigNodeConfController{}, "*:NodeConfList")
	beego.Router("/editNodeConf", &controllers.ConfigNodeConfController{}, "*:EditNodeConf")
	beego.Router("/addNodeConf", &controllers.ConfigNodeConfController{}, "*:AddNodeConf")
	beego.Router("/delNodeConf", &controllers.ConfigNodeConfController{}, "*:DelNodeConf")
	beego.Router("/getIdsNodeConf", &controllers.ConfigNodeConfController{}, "*:GetIdsToNodeConfList")

	// 配置节点管理
	beego.Router("/getNodeList", &controllers.ConfigNodeController{}, "*:GetNodeList")
	beego.Router("/delNode", &controllers.ConfigNodeController{}, "*:DelNode")
	beego.Router("/editNode", &controllers.ConfigNodeController{}, "*:EditNode")

	// 定时任务
	beego.Router("/addTask", &controllers.TaskController{}, "*:Add")
	beego.Router("/editTask", &controllers.TaskController{}, "*:Edit")
	beego.Router("/getTaskList", &controllers.TaskController{}, "*:List")
	beego.Router("/taskRunOrStop", &controllers.TaskController{}, "*:RunOrStop")
	beego.Router("/taskDel", &controllers.TaskController{}, "*:Del")
	beego.Router("/getTaskIdToInfo", &controllers.TaskController{}, "*:GetTaskIdToInfo")
	beego.Router("/taskKill", &controllers.TaskController{}, "*:Kill")
	beego.Router("/taskNodeList", &controllers.TaskNodeController{}, "*:GetNodeList")
	beego.Router("/taskNodeDel", &controllers.TaskNodeController{}, "*:Del")
	beego.Router("/subtasksList", &controllers.TaskController{}, "*:SubtasksList")
	beego.Router("/taskLogList", &controllers.TaskLogController{}, "*:List")

	// etcd
	beego.Router("/getEtcdIpList", &controllers.EtcdController{}, "*:GetEtcdNode")
	beego.Router("/addEtcdNodeConf", &controllers.EtcdController{}, "*:AddEtcdNodeConf")
	beego.Router("/delEtcdNode", &controllers.EtcdController{}, "*:DelEtcdNode")

	// 操作记录
	beego.Router("/getOperationLogList", &controllers.OperationLogController{}, "*:List")

	// 系统
	// 系统用户
	beego.Router("/getUserList", &controllers.UserController{}, "*:List")
	beego.Router("/addUser", &controllers.UserController{}, "*:Add")
	beego.Router("/editUser", &controllers.UserController{}, "*:Edit")

	// 定时任务

	beego.NSNamespace("/*",
		beego.NSRouter("/*", &controllers.BaseController{}, "OPTIONS:Options"),
	)
}
