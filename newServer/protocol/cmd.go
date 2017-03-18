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

type Register_req struct {
	UserName string
	Password string
	UserType string
}

type Register_ack struct {
	Ack string
}

type Login_req struct {
	UserName string
	Password string
	UserType string
}

type Login_ack struct {
	Ack        string
	KeepHeat   string // 0 - 300
	IsLocation string // 0:不上报 1：上报位置
	IsNew      string // 检查更新否
	Money      string // 钱数
}

type Logout_req struct {
	UserName string
}

type Logout_ack struct {
	Ack string
}

type Message_req struct {
	UserName string
}

type Message_ack struct {
	MessageType string
	Message     string
}

type Report_location struct {
	UserName  string
	Latitude  string
	Longitude string
}

type Report_location_ack struct {
	Ack string
}

type Alive struct {
	Ping string
}

type Pkg_info struct {
	Id             string
	Sender         string
	Sender_addr    string
	Sender_phone   string
	Receiver       string
	Receiver_addr  string
	Receiver_phone string
}

type Task_ack struct {
	Deliver_id string
	Time_      string
	Pkg_list   []Pkg_info
}

type Get_Task struct {
	Deliver_id string
	Time_      string
}
