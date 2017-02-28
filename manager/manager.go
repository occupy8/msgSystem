package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"msgSystem/protocol"
	"net"

	slog "github.com/antigloss/go/logger"
)

type Location struct {
	Info string
	Conn net.Conn
}

type ManagerServer struct {
	Cfg         *ManagerServerConfig
	Sessions    map[string]Location
	MsgSessions map[string]Location
	Server      *net.TCPListener
	DbM         *DbManager
}

func NewManagerServer(cfg *ManagerServerConfig) *ManagerServer {
	return &ManagerServer{
		Cfg:         cfg,
		Sessions:    make(map[string]Location),
		MsgSessions: make(map[string]Location),
	}
}

func initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	if err != nil {
		slog.Info("Resolving address:port failed error:%s", err.Error())
		return nil
	}

	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		slog.Info("listen Tcp failed error:%s", err.Error())
		return nil
	}

	slog.Info("ListenTCP: %s", listener.Addr().String())

	return listener
}

func (self *ManagerServer) StartServer(cfg *ManagerServerConfig) {
	//init db
	self.DbM = NewDbManager(cfg)
	self.DbM.ConnectDb(cfg)

	//init server
	hostAndPort := cfg.Listen
	listener := initServer(hostAndPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Info("server Accept error")

			return
		}

		go self.ConnectionHandler(conn)
	}
}

func (self *ManagerServer) ConnectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	slog.Info("Connection from: %s", connFrom)

	readerChannel := make(chan []byte, 16)
	tmpBuffer := make([]byte, 0)

	go self.Reader(conn, readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		slog.Info("read buffer:%s", buffer[:n])
		switch err {
		case nil:
			tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
		default:
			goto DISCONNECT
		}
	}

DISCONNECT:
	conn.Close()
	slog.Info("Closed connection: %s", connFrom)
}

func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		slog.Info("encode cmd error")
		return nil, err
	}
	return buf.Bytes(), nil
}

func (self *ManagerServer) Reader(conn net.Conn, readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			cmd, message := protocol.GetCmdAndMsg(data)
			slog.Info("receive cmd:%04x, message:%s\n", cmd, message)
			switch cmd {
			case protocol.REGISTER_MANAGER:
				var reg protocol.Register
				var user string
				var reg_ack protocol.Register_ack

				err := json.Unmarshal(message, &reg)
				if err != nil {
					slog.Info("json unmarshal error:%s", err.Error())
					break
				}
				//find user from db
				cond := fmt.Sprintf("SELECT name FROM user where name='%s'", reg.Username)

				err = self.DbM.Db.QueryRow(cond).Scan(&user)

				if err == nil {
					reg_ack.Ack = "exist"
				} else {
					//write to db
					stmt, err := self.DbM.Db.Prepare("INSERT user SET name=?,password=?,type=?")
					if err != nil {
						slog.Info("db prepare error:%s", err.Error())
						break
					}

					stmt.Exec(reg.Username, reg.Password, reg.Usertype)

					reg_ack.Ack = "ok"
				}

				buf, err := json.Marshal(reg_ack)
				if err != nil {
					slog.Info("json marshal error:%s", err.Error())

					break
				}
				slog.Info(string(buf))

				send := protocol.Packet(buf, protocol.REGISTER_MANAGER_ACK)
				conn.Write(send)
				break

			case protocol.LOGIN_MANAGER:
				var msg protocol.Login
				var login_ack protocol.Login_ack
				var loc Location

				err := json.Unmarshal(message, &msg)
				if err != nil {
					slog.Info("json unmarshal. error:%s", err.Error())
					break
				}
				//find user from mysql
				var oneUser User
				find := fmt.Sprintf("SELECT id,name,password,type FROM user where name='%s'", msg.Username)
				err = self.DbM.Db.QueryRow(find).Scan(&oneUser.Id, &oneUser.Username, &oneUser.Password, &oneUser.UserType)
				if err != nil {
					login_ack.Ack = "err"
					login_ack.Session = "err"
				} else {
					loc.Conn = conn
					loc.Info = oneUser.Username

					if oneUser.UserType == 3 {
						self.MsgSessions[oneUser.Username] = loc
					} else {
						self.Sessions[oneUser.Username] = loc
					}

					login_ack.Ack = "ok"
					login_ack.Session = oneUser.Username
				}

				//response to client
				buf, err := json.Marshal(login_ack)
				if err != nil {
					slog.Info("json marshal error:%s", err.Error())
					break
				}

				send := protocol.Packet(buf, protocol.LOGIN_MANAGER_ACK)
				conn.Write(send)
				break

			case protocol.LOGOUT_MANAGER:
				var msg protocol.Logout

				err := json.Unmarshal(message, &msg)
				if err != nil {
					slog.Info("json unmarshal error:%s", err.Error())
					break
				}

				//response to client
				var logout_ack protocol.Logout_ack

				logout_ack.Session = msg.Session
				logout_ack.Ack = "ok"

				buf, err := json.Marshal(logout_ack)
				if err != nil {
					slog.Info("json marshal error:%s", err.Error())
					break
				}

				slog.Info(string(buf))

				send := protocol.Packet(buf, protocol.LOGOUT_MANAGER_ACK)
				conn.Write(send)

				conn.Close()

				delete(self.Sessions, msg.Session)
				break

			case protocol.KEEP_ALIVE:
				{
					send := protocol.Packet(message, protocol.KEEP_ALIVE)
					conn.Write(send)
				}
				break
			case protocol.CREATE_GROUP:
				break
			case protocol.GET_TASK:
				{
					var msg protocol.Get_Task

					err := json.Unmarshal(message, &msg)
					if err != nil {
						slog.Info("get task json unmarshal error")
						break
					}

					var lis protocol.Task_ack
					lis.Deliver_id = msg.Deliver_id
					var condition string
					condition = fmt.Sprintf("SELECT Pkg_id,Sender,Sender_addr,Sender_phone,Receiver,Receiver_addr,Receiver_phone FROM pkg_list where Deliver_id='%s'", msg.Deliver_id)
					rows, err := self.DbM.Db.Query(condition)
					if err == nil {
						slog.Info("find pkg info failed")
						break
					}

					for rows.Next() {
						var info protocol.Pkg_info
						err = rows.Scan(&info.Id, &info.Sender, &info.Sender_addr, &info.Sender_phone, &info.Receiver, &info.Receiver_addr, &info.Receiver_phone)
						if err != nil {
							slog.Info("scan error")
							continue
						}

						lis.Pkg_list = append(lis.Pkg_list, info)
					}
					//response to client
					buf, err := json.Marshal(lis)
					if err != nil {
						slog.Info("json marshal error:%s", err.Error())
						break
					}

					send := protocol.Packet(buf, protocol.GET_TASK_ACK)
					conn.Write(send)
				}
				break
			case protocol.REPORT_LOCATION:
				var msg protocol.Report_location

				err := json.Unmarshal(message, &msg)

				//insert sql
				var rloc protocol.Report_location
				var condition string
				condition = fmt.Sprintf("SELECT user,location FROM location where user='%s'", msg.User)
				err = self.DbM.Db.QueryRow(condition).Scan(&rloc.User, &rloc.Gps)
				if err == nil {
					stmt, err := self.DbM.Db.Prepare(`UPDATE location SET location=? WHERE user=?`)
					if err != nil {
						slog.Info("update location failed")
						break
					}
					stmt.Exec(msg.Gps, msg.User)

				} else {
					stmt, err := self.DbM.Db.Prepare(`INSERT location (user, location) values(?,?)`)
					if err != nil {
						slog.Info("prepare location failed")
						break
					}
					stmt.Exec(msg.User, msg.Gps)
				}

				//send to msgserver
				send := protocol.Packet(message, protocol.REPORT_LOCATION)

				for key, value := range self.MsgSessions {
					slog.Info("send to msgServer:%s,key:%s", value.Conn.RemoteAddr().String(), key)
					value.Conn.Write(send)
				}

				break
			}
		}
	}
}
