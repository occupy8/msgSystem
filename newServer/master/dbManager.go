package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id                 int
	Username, Password string
	UserType           int
}

var sqldata map[interface{}]interface{}

type DbManager struct {
	Db *sql.DB
}

func NewDbManager() *DbManager {
	return &DbManager{}
}

func (self *DbManager) ConnectDb() error {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/app?charset=utf8")
	if err != nil {
		return err
	}

	self.Db = db

	return err
}

func (self *DbManager) CheckUser(name string) error {
	var user string
	//find user from db
	cond := fmt.Sprintf("SELECT name FROM user where name='%s'", name)

	err := self.Db.QueryRow(cond).Scan(&user)

	return err
}

func (self *DbManager) CheckUserPassword(name string, pwd string) bool {
	//find user from mysql
	var oneUser User

	find := fmt.Sprintf("SELECT id,name,password,type FROM user where name='%s'", name)
	err := self.Db.QueryRow(find).Scan(&oneUser.Id, &oneUser.Username, &oneUser.Password, &oneUser.UserType)
	if err != nil {
		return false
	} else {
		if name == oneUser.Username && pwd == oneUser.Password {
			return true
		} else {
			return false
		}
	}

	return false
}

func (self *DbManager) InsertUser(name string, pwd string, utype string) error {
	//write to db
	stmt, err := self.Db.Prepare("INSERT user SET name=?,password=?,type=?")
	if err != nil {
		return err
	}

	stmt.Exec(name, pwd, utype)

	return err
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
