package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"msgSystem/newServer/protocol"
	"net/http"

	slog "github.com/antigloss/go/logger"
)

var DbM *DbManager

func AnswerTaskPush(resp_byte []byte, w http.ResponseWriter) {

	fmt.Println("answer: " + string(resp_byte))
	w.Write(resp_byte)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var reg protocol.Register_req
	var reg_ack protocol.Register_ack

	if r.Method == "GET" {
		slog.Info("method: %s\n", r.Method)
	} else if r.Method == "POST" {

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		slog.Info("register %s\n", result)

		err := json.Unmarshal(result, &reg)
		if err != nil {
			slog.Info("json unmarshal error:%s", err.Error())
			return
		}

		err = DbM.CheckUser(reg.UserName)
		if err == nil {
			reg_ack.Ack = "exist"
		} else {
			//write to db
			DbM.InsertUser(reg.UserName, reg.Password, reg.UserType)

			reg_ack.Ack = "ok"
		}

		buf, err := json.Marshal(reg_ack)
		if err != nil {
			slog.Info("json marshal error:%s", err.Error())

			return
		}
		slog.Info(string(buf))

		AnswerTaskPush(buf, w)
	} else {
		AnswerTaskPush([]byte("FAIL"), w)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var login protocol.Login_req
	var login_ack protocol.Login_ack

	if r.Method == "GET" {
		slog.Info("method: %s\n", r.Method)
	} else if r.Method == "POST" {

		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		slog.Info("register %s\n", result)

		err := json.Unmarshal(result, &login)
		if err != nil {
			slog.Info("json unmarshal error:%s", err.Error())
			return
		}

		ret := DbM.CheckUserPassword(login.UserName, login.Password)
		if ret == true {
			login_ack.Ack = "ok"
		} else {
			login_ack.Ack = "error"
		}

		buf, err := json.Marshal(login_ack)
		if err != nil {
			slog.Info("json marshal error:%s", err.Error())

			return
		}
		slog.Info(string(buf))

		AnswerTaskPush(buf, w)
	} else {
		AnswerTaskPush([]byte("FAIL"), w)
	}
}

func StartServer() {
	//init db
	DbM = NewDbManager()
	DbM.ConnectDb()

	http.HandleFunc("/v1/register", registerHandler)
	http.HandleFunc("/v1/login", loginHandler)

	slog.Info("http server listening :%s", Config.ServerAddress.IpAddress)
	http.ListenAndServe(Config.ServerAddress.IpAddress, nil)
}