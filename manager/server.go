package main

import (
	"fmt"
	"log"
	"os"
	"time"

	slog "github.com/antigloss/go/logger"
)

const VERSION string = "0.10"

func version() {
	fmt.Printf("manager server version:%s\n", VERSION)
}

func initLog(cfg *ManagerServerConfig) {
	err := os.MkdirAll(cfg.LogPath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	err = slog.Init(cfg.LogPath, 20, 4, 30, false)
	if err != nil {
		log.Println("logger init failed:", err)
		return
	}

	slog.SetLogThrough(true)

	slog.Info("--------------------------------------")
	slog.Info("-----------manager server--------------")
	slog.Info("--------------------------------------")
}

func main() {
	version()
	cfg := NewManagerServerConfig("config.json")

	err := cfg.LoadConfig()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	cfg.DumpConfig()

	initLog(cfg)

	server := NewManagerServer(cfg)

	go server.StartServer(cfg)

	for {
		time.Sleep(3 * time.Second)
	}
}
