package main

import (
	"database/sql"

	slog "github.com/antigloss/go/logger"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id                 int
	Username, Password string
	UserType           int
}

var sqldata map[interface{}]interface{}

type DbManager struct {
	Cfg *ManagerServerConfig
	Db  *sql.DB
}

func NewDbManager(cfg *ManagerServerConfig) *DbManager {
	return &DbManager{
		Cfg: cfg,
	}
}

func (self *DbManager) ConnectDb(cfg *ManagerServerConfig) error {
	db, err := sql.Open("mysql", "root:123456@tcp(123.125.89.76:3306)/app?charset=utf8")
	if err != nil {
		slog.Info("connect Db error")
		return err
	}

	self.Db = db

	return err
}
