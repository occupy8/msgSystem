package main

import (
	"database/sql"
	"fmt"
	"msgSystem/newServer/protocol"

	"time"

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

func (self *DbManager) GetTaskList(deliverId string, lis *protocol.Task_ack) error {

	lis.Deliver_id = deliverId
	var condition string
	condition = fmt.Sprintf("SELECT Pkg_id,Sender,Sender_addr,Sender_phone,Receiver,Receiver_addr,Receiver_phone FROM pkg_list where Deliver_id='%s'", deliverId)
	rows, err := self.Db.Query(condition)
	if err != nil {
		return err
	}

	for rows.Next() {
		var info protocol.Pkg_info
		err = rows.Scan(&info.Id, &info.Sender, &info.Sender_addr, &info.Sender_phone, &info.Receiver, &info.Receiver_addr, &info.Receiver_phone)
		if err != nil {
			continue
		}

		lis.Pkg_list = append(lis.Pkg_list, info)
	}

	return err
}

func (self *DbManager) InsertLocation(la string, lo string, user string) error {
	stmt, err := self.Db.Prepare(`INSERT location (user, lo, la, time) values(?,?,?,?)`)
	if err != nil {

		return err
	}
	stmt.Exec(user, lo, la, string(time.Now().Unix()))

	return nil
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
