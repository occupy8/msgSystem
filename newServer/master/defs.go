package main

const (
	MaxRunningTasks = 20 // 最大并发运行切片任务的数量
	SERVERADDRESS   = ":9999"

	LogFilePath       = "/opt/project/master/log/"
	DefaultConfPath   = "/opt/project/master/conf/master.toml"
	DefaultConfigPath = "/opt/project/master/conf/"
	DefaultHlsCreater = "/opt/project/master/hls/bin/hls_creater"
)

var (
	Config     MasterConfig
	ConfigFile string
)
