package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MsgServerConfig struct {
	Configfile string
	Listen     string
	Server     string
}

func NewMsgServerConfig(configfile string) *MsgServerConfig {
	return &MsgServerConfig{
		Configfile: configfile,
	}
}

func (self *MsgServerConfig) LoadConfig() error {
	file, err := os.Open(self.Configfile)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&self)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(self)

	return nil
}

func (self *MsgServerConfig) DumpConfig() {
	fmt.Printf("listen:%s,Server:%s\n",
		self.Listen, self.Server)
}
