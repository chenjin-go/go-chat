package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       string //当前模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       "999",
	}

	//链接server
	dial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = dial

	//返回对象
	return client
}

func (this *Client) DealResponse() {
	//有数据接收，就打印
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) menu() bool {
	fmt.Println("1.1vn")
	fmt.Println("2.1v1")
	fmt.Println("3.rename")
	fmt.Println("0.exit")

	fmt.Scanln(&this.flag)
	int, err := strconv.Atoi(this.flag)

	if err != nil {
		fmt.Println(" input 0-3 ")
		int = 999
	}

	if int >= 0 && int <= 3 {
		return true
	} else {
		fmt.Println("flag err")
		return false
	}
}

func (this *Client) Rename() bool {
	fmt.Println(">>>>>>new name>>>>>>")
	fmt.Scanln(&this.Name)

	sendMsg := "rename|" + this.Name + "\n"
	err := this.sendMsg(sendMsg)
	if !err {
		fmt.Println("rename err")
		return false
	}
	fmt.Println("rename success")
	return true
}

//私聊
func (this *Client) privateChat() {
	var remoteName string
	var msg string

	this.selectUsets()
	fmt.Println(">>>>>>input to user name>>>>>>")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>input msg>>>")
		fmt.Scanln(&msg)

		for msg != "exit" {
			if len(msg) > 0 {
				sendMsg := "to|" + remoteName + "|" + msg + "\n"
				err := this.sendMsg(sendMsg)
				if !err {
					fmt.Println("to "+remoteName+" msg err:", err)
					break
				}
			}

			msg = ""
			fmt.Scanln(&msg)
		}
		this.selectUsets()
		fmt.Println(">>>>>>input to user name>>>>>>")
		fmt.Scanln(&remoteName)
	}
}

//公聊
func (this *Client) PublicChat() {
	var msg string
	fmt.Println(">>>>>send msg, send 'exit' out>>>>>>")
	fmt.Scanln(&msg)

	for msg != "exit" {

		if len(msg) != 0 {
			sendMsg := msg + "\n"
			err := this.sendMsg(sendMsg)
			if !err {
				break
			}
		}
		msg = ""
		fmt.Scanln(&msg)
	}
}

func (this *Client) Run() {

	for this.flag != "0" {

		for this.menu() != true {
		}
		//根据flag处理业务
		switch this.flag {
		case "1":
			//1vn
			fmt.Println("1vN...")
			this.PublicChat()
		case "2":
			//1v1
			fmt.Println("1v1...")
			this.privateChat()
		case "3":
			//rename
			fmt.Println("rename...")
			this.Rename()
		default:

		}
	}
}

func (this *Client) sendMsg(sendMsg string) bool {
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (this *Client) selectUsets() {
	sendMsg := "who\n"
	err := this.sendMsg(sendMsg)
	if !err {
		fmt.Println("select users err:", err)
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set ip")
	flag.IntVar(&serverPort, "port", 8888, "set port")
}

func main() {
	//解析flag
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>error>>>>>>>>>")
		return
	}

	fmt.Println(">>>>>>>>>>success>>>>>>>>>>")

	// 监听消息
	go client.DealResponse()

	//启动业务
	client.Run()
}
