package main

import (
	"github.com/bbangert/toml"
)

type General struct {
	MaxTasks int    `toml:"max_run_tasks"`
	LogPath  string `toml:"log_path"`
	ConfFile string `toml:"conf_file"`
	ConfPath string `toml:"conf_path"`
}
type ServerAddr struct {
	IpAddress string `toml:"ip_address"`
}

type MasterConfig struct {
	ServerAddress ServerAddr `toml:"server_addr"`
	Gen           General    `toml:"general"`
}

type ServerConf struct {
	path string
	conf MasterConfig
}

func decodeFile(path string, obj *MasterConfig) error {
	var err error

	_, err = toml.DecodeFile(path, obj)

	return err
}

func NewMasterConfig(path string, conf *MasterConfig) error {

	err := decodeFile(path, conf)
	if err != nil {
		return err
	}

	return err
}
