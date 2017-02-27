package main

import (
	"encoding/json"
	"fmt"
	"msgSystem/protocol"
	"net"
)

type ManagerServerClient struct {
	Cfg     *ClientConfig
	Client  net.Conn
	Session string
}

func NewManagerServerClient(cfg *ClientConfig) *ManagerServerClient {
	return &ManagerServerClient{
		Cfg: cfg,
	}
}

func (self *ManagerServerClient) StartClient(cfg *ClientConfig) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.Manager)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("StartClient error")
		return
	}

	self.Client = conn

	go self.ConnectionHandler(conn)

}

func (self *ManagerServerClient) ConnectionHandler(conn net.Conn) {
	readerChannel := make(chan []byte, 16)
	tmpBuffer := make([]byte, 0)

	go self.Reader(conn, readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		switch err {
		case nil:
			tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
		default:
			goto DISCONNECT
		}
	}

DISCONNECT:
	err := conn.Close()
	checkError(err, "Close:")
}

func (self *ManagerServerClient) Reader(conn net.Conn, readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			cmd, message := protocol.GetCmdAndMsg(data)
			fmt.Printf("ManagerServer: cmd:%04x, message:%s\n", cmd, message)
			switch cmd {
			case protocol.LOGIN_MANAGER_ACK:
				var msg protocol.Login_ack

				json.Unmarshal(message, &msg)

				self.Session = msg.Session

				fmt.Println("receive LOGIN_MANAGER_ACK")
				fmt.Println(string(message))

			case protocol.LOGOUT_MANAGER_ACK:
				var msg protocol.Logout

				json.Unmarshal(message, &msg)

			case protocol.REGISTER_MANAGER_ACK:
				var msg protocol.Register_ack

				json.Unmarshal(message, &msg)

				if msg.Ack != "ok" {
					fmt.Println("please input another username!")
				}
			case protocol.REPORT_LOCATION_ACK:
				var msg protocol.Report_location_ack

				json.Unmarshal(message, &msg)

			}
		}
	}
}

func (self *ManagerServerClient) SendToServer(data interface{}, cmd int32) {

	message, err := json.Marshal(data)

	checkError(err, "json marshal")

	encode := protocol.Packet(message, cmd)

	self.Client.Write(encode)
}

func (self *ManagerServerClient) SendRegisterCmd(user string, password string, usertype int) {
	var send protocol.Register

	send.Username = user
	send.Usertype = usertype
	send.Password = password

	self.SendToServer(send, protocol.REGISTER_MANAGER)
}

func (self *ManagerServerClient) SendLoginManagerCmd(user string, password string, userType int) {
	var send protocol.Login

	send.Username = user
	send.Password = password
	send.Usertype = userType

	self.SendToServer(send, protocol.LOGIN_MANAGER)
}

func (self *ManagerServerClient) SendLogoutManagerCmd(user string, session string) {
	var send protocol.Logout

	send.UserName = user
	send.Session = session

	self.SendToServer(send, protocol.LOGOUT_MANAGER)
}

func (self *ManagerServerClient) SendReportLocationCmd(location string, username string) {
	var send protocol.Report_location

	send.Gps = location
	send.User = username

	self.SendToServer(send, protocol.REPORT_LOCATION)
}

func checkError(error error, info string) {
	if error != nil {
		panic("ERROR: " + info + " " + error.Error()) // terminate program
	}
}
