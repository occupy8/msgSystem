package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ManagerServerConfig struct {
	Configfile string
	Listen     string
	LogPath    string
}

func NewManagerServerConfig(config string) *ManagerServerConfig {
	return &ManagerServerConfig{
		Configfile: config,
	}
}

func (self *ManagerServerConfig) LoadConfig() error {
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

func (self *ManagerServerConfig) DumpConfig() {
	fmt.Printf("listen:%s\n", self.Listen)
}
