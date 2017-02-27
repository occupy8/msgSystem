//
// Copyright 2014 Hong Miao (miaohong@miaohong.org). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	ConstHeader         = "www.01happy.com"
	ConstHeaderLength   = 15
	ConstSaveDataLength = 4
	ConstCmdLength      = 4
)

const (
	LOGIN_MANAGER     = 0x0001
	LOGIN_MANAGER_ACK = 0xF001

	P2P_SEND_MESSAGE     = 0x0002
	P2P_SEND_MESSAGE_ACK = 0xF002

	REPORT_LOCATION     = 0x0003
	REPORT_LOCATION_ACK = 0xF003

	REGISTER_MANAGER     = 0x0004
	REGISTER_MANAGER_ACK = 0xF004

	LOGOUT_MANAGER     = 0x0005
	LOGOUT_MANAGER_ACK = 0xF005

	MESSAGE_GROUP     = 0x0006
	MESSAGE_GROUP_ACK = 0xF006

	LOGIN_MESSAGE     = 0x0007
	LOGIN_MESSAGE_ACK = 0xF007

	KEEP_ALIVE     = 0x0008
	KEEP_ALIVE_ACK = 0xF008

	CREATE_GROUP     = 0x0009
	CREATE_GROUP_ACK = 0xF009

	DEL_GROUP     = 0x0010
	DEL_GROUP_ACK = 0xF010

	ADD_FRIEND_IN_GROUP     = 0x0011
	ADD_FRIEND_IN_GROUP_ACK = 0xF011

	EXIT_GROUP     = 0x0012
	EXIT_GROUP_ACK = 0xF012
)

type AddInGroup struct {
	Groupid  string
	Username string
}

type AddInGroupAck struct {
	Ack string
}

type ExitGroup struct {
	Groupid  string
	Username string
}

type ExitGroupAck struct {
	Ack string
}

type CreateGroup struct {
	Groupid  string
	Username string
}

type CreateGroupAck struct {
	Ack string
}

type DelGroup struct {
	Groupid  string
	Username string
}

type DelGroupAck struct {
	Ack string
}

type Register struct {
	Username string
	Password string
	Usertype string
}

type Register_ack struct {
	Ack string
}

type Login struct {
	Username string
	Password string
	Usertype string
}

type Login_ack struct {
	Session string
	Ack     string
}

type LoginMsg struct {
	Username string
	Session  string
}

type LoginMsg_ack struct {
	Ack string
}

type Logout struct {
	UserName string
	Session  string
}

type Logout_ack struct {
	Session string
	Ack     string
}

type Send_message struct {
	Session_from string
	Message      string
	Session_to   string
}

type Send_message_ack struct {
	Session_from string
	Ack          string
	Session_to   string
}

type Message_group struct {
	From    string
	Message string
}

type Message_group_ack struct {
	Ack string
}

type Report_location struct {
	User string
	Gps  string
}

type Report_location_ack struct {
	Ack string
}

type Alive_ping struct {
	ping string
}

func Packet(message []byte, cmd int32) []byte {
	return append(append(append([]byte(ConstHeader), IntToBytes(int32(len(message)))...), IntToBytes(cmd)...), message...)
}

func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	length := int32(len(buffer))
	var i int32
	var tmp int32

	tmp = ConstHeaderLength + ConstSaveDataLength + ConstCmdLength

	for i = 0; i < length; i = i + 1 {
		if length < i+tmp {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+tmp+messageLength {
				break
			}

			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+tmp+messageLength]
			readerChannel <- data

			i += tmp + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int32) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int32 {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int32(x)
}

func GetCmdAndMsg(endcode []byte) (int32, []byte) {
	var cmd int32
	tmpBuffer := make([]byte, 0)

	cmdEncode := endcode[:ConstCmdLength]

	cmd = BytesToInt(cmdEncode)

	//messageEncode := endcode[ConstCmdLength:]
	tmpBuffer = append(tmpBuffer, endcode[ConstCmdLength:]...)

	//fmt.Println(cmd)
	//fmt.Println("-----------------------")
	//fmt.Println(string(tmpBuffer))

	return cmd, tmpBuffer
}
