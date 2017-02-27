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
	cfg := NewClientConfig("config.json")

	err := cfg.LoadConfig()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	cfg.DumpConfig()

	MsgServer := NewMsgServerClient(cfg)
	ManagerServer := NewManagerServerClient(cfg)

	go ManagerServer.StartClient(cfg)
	go MsgServer.StartClient(cfg)

	var userName string

	fmt.Println("user:")
	fmt.Scanf("%s", &userName)

	time.Sleep(2 * time.Second)
	ManagerServer.SendRegisterCmd(userName, "123456", 1)
	time.Sleep(2 * time.Second)
	ManagerServer.SendLoginManagerCmd(userName, "123456", 1)
	time.Sleep(2 * time.Second)
	MsgServer.SendLoginMsgServer(userName, ManagerServer.Session)
	//ManagerServer.SendLogoutManagerCmd("zozo", ManagerServer.Session)

	//for {
	//	anounce := "zozo " + strconv.FormatInt(time.Now().UnixNano(), 10)
	//	ManagerServer.SendReportLocationCmd("(100,99)", anounce)
	//	time.Sleep(10 * time.Second)
	//}
	time.Sleep(3 * time.Second)
	var msg string
	//to = "co"
	msg = "hello, there!!!"

	for {
		//fmt.Println("msg:")
		//fmt.Scanf("%s", &msg)
		time.Sleep(2 * time.Second)
		//MsgServer.SendP2pMessage(userName, msg, to)
		MsgServer.SendGroupMessage(userName, msg)
	}
}
