package main

import (
	"fmt"

	"github.com/kardianos/service"

	slog "github.com/antigloss/go/logger"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		slog.Info("Running in terminal.")
	} else {
		slog.Info("Running under service manager.")
	}

	go p.run()

	return nil
}

func (p *program) run() {
	err := NewMasterConfig(DefaultConfPath, &Config)
	if err != nil {
		slog.Info("config init err")
		return
	}

	fmt.Println(Config)

	StartServer()
}

func (p *program) Stop(s service.Service) error {
	return nil
}
