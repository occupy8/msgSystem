package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	slog "github.com/antigloss/go/logger"
	"github.com/kardianos/service"
)

const VERSION = "1.0"

func Init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := os.MkdirAll(LogFilePath, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	Init()

	err := slog.Init(LogFilePath, 20, 4, 30, false)
	if err != nil {
		log.Println("logger init failed:", err)
		return
	}

	slog.SetLogThrough(true)

	slog.Info("--------------------------------------")
	slog.Info("-----------master server--------------")
	slog.Info("--------------------------------------")

	configPath := flag.String("c", DefaultConfPath, "master server configuration file")
	deamon := flag.Bool("d", false, "deamon")
	version := flag.Bool("v", false, "version")

	flag.Parse()

	ConfigFile = *configPath
	if *version == true {
		fmt.Println("master version:", VERSION)
		return
	}

	svcConfig := &service.Config{
		Name:        "master_server",
		DisplayName: "Go Service: Master Server",
		Description: "This is a hls master server service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		slog.Info("service New error")
		return
	}

	if *deamon == true {
		err = service.Control(s, "uninstall")
		if err != nil {
			slog.Info("service control uninstall error")
		}

		slog.Info("service uninstalled")
		err = service.Control(s, "install")
		if err != nil {
			slog.Info("service control install error")
		}

		slog.Info("service installed")
		err = service.Control(s, "restart")
		if err != nil {
			log.Println(err)
		}

		slog.Info("service started")
		return
	}

	err = s.Run()
	if err != nil {
		slog.Info("master server run failed")
	}
}
