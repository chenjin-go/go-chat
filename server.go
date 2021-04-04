package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	// 用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//广播管道
	Message chan string
}

// 创建一个Server
func newServer(ip string, port int) *Server {
	return &Server{Ip: ip, Port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
}

//监听Message 广播消息
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()

		// 将消息发送给全部的在线user
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 发送广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) handler(conn net.Conn) {
	//链接后的业务
	//fmt.Println("链接成功")
	user := newUser(conn, this)

	user.Online()

	//监听是否活跃
	isLive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			read, err := conn.Read(buf)
			if read == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户消息
			msg := string(buf[:read-1])
			user.DoMessage(msg)

			// 捕捉操作，确认活跃
			isLive <- true
		}
	}()

	// 阻塞，防止死亡
	for {
		select {
		case <-isLive:
			// 激活select，更新定时器
		case <-time.After(time.Second * 60 * 30):
			// 超时
			// 将当前User强制关闭
			user.SendMsg(" to out")
			//销毁资源
			close(user.C)
			user.conn.Close()
			return // 或者 runtime.Goexit()
		}
	}

}

// 启动服务
func (this *Server) start() {
	// socket listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	// close listen socket
	defer listen.Close()

	//启动监听Message的管道
	go this.ListenMessage()

	for {
		// accept
		accept, err := listen.Accept()
		if err != nil {
			fmt.Println("listen accept err:", err)
			continue
		}
		// do handler
		go this.handler(accept)
	}
}
