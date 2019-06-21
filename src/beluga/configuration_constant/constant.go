package configuration_constant

const (
	// 节点配置目录
	NODE_CONF_DIR = "/node/conf/"

	// 配置节点发现目录
	CONFIGURATION_REGISTER_DIR = "/configuration/register/"

	// 同步中心配置
	CONFIGURATION_CONF_SYNC = "/configuration/conf/" // +appid_namespace=path

	// 节点配置更新目录
	CONFIGURATION_NODE_CONF = "/configuration/node/" // + ip + appid_namespace=path

	// 中心节点目录名，用来判断中心节点是否存在。
	BELUGA_MASTER_DIR = "/beluga_master/"
)

const (
	CONFIGURATION_REDIS_KEY = "beluga_configuration"
)