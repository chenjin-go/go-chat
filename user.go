package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建一个用户
func newUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{userAddr, userAddr, make(chan string), conn, server}
	// 启动监听，有消息就发送给用户
	go user.ListenMessage()
	return user
}

//上线
func (this *User) Online() {
	//用户上线广播，加入map
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播消息
	this.server.BroadCast(this, "log in ")
	fmt.Println(this.Name + ":log in ")
}

//下线
func (this *User) Offline() {
	//用户下线广播，删除map
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播消息
	this.server.BroadCast(this, "log out ")
	fmt.Println(this.Name + ":log out ")

}

//给自己发送消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//查询出当前所有在线人数
func (this *User) ShowLine() {
	//查询当前在线用户
	this.server.mapLock.Lock()
	for _, user := range this.server.OnlineMap {
		onlinMsg := "[" + user.Addr + "]" + user.Name + "\n"
		this.SendMsg(onlinMsg)
	}
	this.server.mapLock.Unlock()
}

//用户发送消息
func (this *User) DoMessage(msg string) {

	if msg == "who" {
		//查询当前在线用户
		this.ShowLine()
	} else if len(msg) > 7 && strings.Split(msg, "|")[0] == "rename" {
		//改名
		name := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[name]
		if ok {
			this.SendMsg("name is repeat\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[name] = this
			this.server.mapLock.Unlock()

			this.Name = name
			this.SendMsg("Successfully renamed：" + name)
		}
	} else if len(msg) > 4 && "to|" == msg[:3] {
		// 消息格式 to|user|msg
		//1 获取对方用户
		sp := strings.Split(msg, "|")
		if sp[1] == "" {
			this.SendMsg(" error:name is null ")
			return
		}
		toUser, ok := this.server.OnlineMap[sp[1]]
		if !ok {
			this.SendMsg(" error:get user error, user is null ")
			return
		}
		//2 发送消息 （通过user对象）
		if sp[2] == "" {
			this.SendMsg(" error: message is null ")
			return
		}
		toUser.SendMsg(this.Name + " to you: " + sp[2])

	} else {
		//消息广播
		this.server.BroadCast(this, msg)
	}
}

// 监听当前user chan的方法
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		write, err := this.conn.Write([]byte(msg + "\n"))
		fmt.Println(write, err)
	}
}
