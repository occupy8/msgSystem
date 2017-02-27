package main

import (
	"encoding/json"
	"fmt"
	"msgSystem/protocol"
	"net"
)

type Location struct {
	Info string
	Conn net.Conn
}

type MsgServer struct {
	Cfg           *MsgServerConfig
	Sessions      map[string]Location
	Server        *net.TCPListener
	Client        net.Conn
	ClientSession string
}

func NewMsgServer(cfg *MsgServerConfig) *MsgServer {
	return &MsgServer{
		Cfg:      cfg,
		Sessions: make(map[string]Location),
	}
}

func (self *MsgServer) StartServer(cfg *MsgServerConfig) {
	hostAndPort := cfg.Listen
	listener := self.InitServer(hostAndPort)
	self.Server = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept error")
		}

		go self.ConnectionHandler(conn)
	}
}

func (self *MsgServer) InitServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	checkError(err, "Resolving address:port failed: '"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	checkError(err, "ListenTCP: ")
	println("Listening to: ", listener.Addr().String())
	return listener
}

func (self *MsgServer) ConnectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	println("Connection from: ", connFrom)

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
	for key, value := range self.Sessions {
		if value.Conn == conn {
			delete(self.Sessions, key)
			break
		}
	}

	err := conn.Close()
	println("Closed connection:", connFrom)
	checkError(err, "Close:")
}

func (self *MsgServer) Reader(conn net.Conn, readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			cmd, message := protocol.GetCmdAndMsg(data)
			fmt.Printf("Reader: cmd:%04x, message:%s\n", cmd, message)
			switch cmd {
			case protocol.LOGIN_MESSAGE:
				var msg protocol.LoginMsg

				json.Unmarshal(message, &msg)

				loc := Location{
					Info: msg.Session,
					Conn: conn,
				}

				fmt.Println("login:" + msg.Session)
				self.Sessions[msg.Session] = loc

			case protocol.P2P_SEND_MESSAGE:
				var msg protocol.Send_message
				var response protocol.Send_message_ack

				json.Unmarshal(message, &msg)
				locat, exist := self.Sessions[msg.Session_to]
				if exist {
					send := protocol.Packet(message, protocol.P2P_SEND_MESSAGE)
					locat.Conn.Write(send)
				} else {
					response.Session_from = msg.Session_from
					response.Ack = "not exist"
					response.Session_to = ""

					encode, err := json.Marshal(response)
					checkError(err, "json")

					send := protocol.Packet(encode, protocol.P2P_SEND_MESSAGE_ACK)
					conn.Write(send)
				}
				//
			case protocol.LOGOUT_MANAGER:
				var msg protocol.Logout

				err := json.Unmarshal(message, &msg)
				if err != nil {
					fmt.Printf("json error")
				}

				//response to client
				var logout_ack protocol.Logout_ack

				logout_ack.Session = msg.Session
				logout_ack.Ack = "ok"

				buf, err := json.Marshal(logout_ack)

				fmt.Println(string(buf))

				send := protocol.Packet(buf, protocol.LOGOUT_MANAGER_ACK)
				conn.Write(send)

				fmt.Println("receiv logout cmd 12")

				conn.Close()
				fmt.Println("receiv logout cmd 3")

				delete(self.Sessions, msg.Session)
				fmt.Println("receiv logout cmd 4")
			case protocol.P2P_SEND_MESSAGE_ACK:
				var msg protocol.Send_message_ack

				json.Unmarshal(message, &msg)
				locat, exist := self.Sessions[msg.Session_to]
				if exist {
					send := protocol.Packet(message, protocol.P2P_SEND_MESSAGE_ACK)
					locat.Conn.Write(send)
				}
				//
			case protocol.KEEP_ALIVE:

				send := protocol.Packet(message, protocol.KEEP_ALIVE)
				conn.Write(send)

			case protocol.MESSAGE_GROUP:
				var msg protocol.Message_group

				json.Unmarshal(message, &msg)

				send := protocol.Packet(message, protocol.MESSAGE_GROUP)

				for key, value := range self.Sessions {
					fmt.Printf("send to %s-> %s\n", key, value.Conn.RemoteAddr().String())
					value.Conn.Write(send)
				}
			case protocol.MESSAGE_GROUP_ACK:
				//
			}
		}
	}
}

func (self *MsgServer) StartClient(cfg *MsgServerConfig) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.Server)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("login error")
		return
	}

	self.Client = conn

	self.SendLogin(conn)

	go self.ClientConnectionHandler(conn)
}

func (self *MsgServer) ClientConnectionHandler(conn net.Conn) {
	ClientReaderChannel := make(chan []byte, 16)
	tmpBuffer := make([]byte, 0)

	go self.ClientReader(conn, ClientReaderChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)

		switch err {
		case nil:
			tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), ClientReaderChannel)
		default:
			goto DISCONNECT
		}
	}

DISCONNECT:
	err := conn.Close()
	checkError(err, "Close:")
}

func (self *MsgServer) ClientReader(conn net.Conn, readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			cmd, message := protocol.GetCmdAndMsg(data)
			fmt.Printf("ClientReader: cmd:%04x, message:%s\n", cmd, message)
			switch cmd {
			case protocol.LOGIN_MANAGER_ACK:
				var msg protocol.Login_ack

				json.Unmarshal(message, &msg)

				self.ClientSession = msg.Session

			case protocol.REPORT_LOCATION:
				var msg protocol.Report_location

				json.Unmarshal(message, &msg)

				//broacast to all user

				send := protocol.Packet(message, protocol.REPORT_LOCATION)

				fmt.Println(self.Sessions)

				for key, value := range self.Sessions {
					fmt.Printf("send to %s-> %s\n", key, value.Conn.RemoteAddr().String())
					value.Conn.Write(send)
				}
			}
		}
	}
}

func (self *MsgServer) SendLogin(conn net.Conn) {
	var l protocol.Login

	l.Username = "msgserver"
	l.Password = "msgserver"
	l.Usertype = "3"

	message, err := json.Marshal(l)

	checkError(err, "json marshal")
	fmt.Println("------------------")
	fmt.Println(string(message))

	encode := protocol.Packet(message, protocol.LOGIN_MANAGER)
	fmt.Println(encode)

	conn.Write(encode)
}

func checkError(error error, info string) {
	if error != nil {
		fmt.Println("ERROR: " + info + " " + error.Error()) // terminate program
	}
}
