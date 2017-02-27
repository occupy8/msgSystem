package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ClientConfig struct {
	CfgFile string
	Message string
	Manager string
}

func NewClientConfig(configfile string) *ClientConfig {
	return &ClientConfig{
		CfgFile: configfile,
	}
}

func (self *ClientConfig) LoadConfig() error {
	file, err := os.Open(self.CfgFile)
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

func (self *ClientConfig) DumpConfig() {
	fmt.Printf("messageServer:%s,managerServer:%s\n",
		self.Message, self.Manager)
}
