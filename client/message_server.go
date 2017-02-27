package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"msgSystem/protocol"
	"net"
)

type MsgServerClient struct {
	Cfg    *ClientConfig
	Client net.Conn
}

func NewMsgServerClient(cfg *ClientConfig) *MsgServerClient {
	return &MsgServerClient{
		Cfg: cfg,
	}
}

func (self *MsgServerClient) StartClient(cfg *ClientConfig) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", cfg.Message)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "dialTcp")

	self.Client = conn

	go self.ConnectionHandler(conn)

}

func (self *MsgServerClient) ConnectionHandler(conn net.Conn) {
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

func (self *MsgServerClient) Reader(conn net.Conn, readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			cmd, message := protocol.GetCmdAndMsg(data)
			fmt.Printf("MsgServerClient: cmd:%04x, message:%s\n", cmd, message)
			switch cmd {
			case protocol.P2P_SEND_MESSAGE:
				var msg protocol.Send_message

				json.Unmarshal(message, &msg)

			case protocol.P2P_SEND_MESSAGE_ACK:
				var msg protocol.Send_message_ack

				json.Unmarshal(message, &msg)

			case protocol.MESSAGE_GROUP:
				var msg protocol.Message_group

				json.Unmarshal(message, &msg)

			case protocol.MESSAGE_GROUP_ACK:
				var msg protocol.Message_group_ack

				json.Unmarshal(message, &msg)
			case protocol.REPORT_LOCATION:
				var msg protocol.Report_location

				json.Unmarshal(message, &msg)
			}
		}
	}
}

func (self *MsgServerClient) SendToServer(data interface{}, cmd int32) {

	message, err := json.Marshal(data)

	checkError(err, "json marshal")
	fmt.Println("------------------")
	fmt.Println(string(message))

	encode := protocol.Packet(message, cmd)

	self.Client.Write(encode)
}

func (self *MsgServerClient) SendP2pMessage(from string, msg string, to string) {
	var send protocol.Send_message

	send.Message = msg
	send.Session_from = from78u7,.
	send.Session_to = to

	self.SendToServer(send, protocol.P2P_SEND_MESSAGE)
}

func (self *MsgServerClient) SendGroupMessage(from string, msg string) {
	var send protocol.Message_group

	send.Message = msg
	send.From = from

	self.SendToServer(send, protocol.MESSAGE_GROUP)
}

func (self *MsgServerClient) SendLoginMsgServer(user string, session string) {
	var send protocol.LoginMsg

	send.Username = user
	send.Session = session

	self.SendToServer(send, protocol.LOGIN_MESSAGE)
}

func (self *MsgServerClient) Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MsgcheckError(error error, info string) {
	if error != nil {
		panic("ERROR: " + info + " " + error.Error()) // terminate program
	}
}
