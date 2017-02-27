package main

import (
	"fmt"
	"time"
)

const VERSION string = "0.10"

func version() {
	fmt.Printf("msg server version:%s\n", VERSION)
}

func main() {
	version()
	cfg := NewMsgServerConfig("config.json")

	err := cfg.LoadConfig()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	cfg.DumpConfig()

	server := NewMsgServer(cfg)

	go server.StartServer(cfg)

	go server.StartClient(cfg)

	for {
		time.Sleep(3 * time.Second)
	}
}
