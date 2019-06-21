package configuration_constant

type NodeConf struct {
	Appid     string            `json:"appid"`
	Namespace []NodeConfNamespace `json:"namespace"`
}

type NodeConfNamespace struct {
	Name string `json:"name"`
	Path string `json:"path"`
}