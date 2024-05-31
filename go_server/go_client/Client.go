package main

import(
	"fmt"
	"io"		
	"flag" //解析命令行
	"os"
	"net"
)

type Client struct{
	ServerIP string
	ServerPort int
	name string
	conn net.Conn
	flag int
}

func NewClient(ServerIP_ string,ServerPort_ int) *Client{
	client := &Client{
		ServerIP : ServerIP_,
		ServerPort : ServerPort_,
		flag : 10,
	}
	conn, err := net.Dial("tcp",fmt.Sprintf("%s:%d",ServerIP_,ServerPort_))
	if err != nil{
		fmt.Println("Dial error:",err)
		return nil
	}
	client.conn = conn
	return client
}

func (this *Client) Respond(){
	//一旦conn缓冲区有数据，就将数据发送到标准输出上
	io.Copy(os.Stdout,this.conn)
}

func (this * Client) menu() bool{
	fmt.Println("请根据一下功能选择序号..........")
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 修改名称")
	fmt.Println("0. 退出")

	var flag int
	fmt.Scanln(&flag)
	if flag < 0 && flag > 3 {
		fmt.Println("请输入正确的序号！")
		return false
	}
	this.flag = flag
	return true
}
func (this * Client) PushlicMode(){
	var msg string
	fmt.Println("请输入你需要发送的消息,如若退出请输入exit..........")

	fmt.Scanln(&msg)
	for msg != "exit"{
		_, err :=this.conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("conn Write err:", err)
			return
		}		

		msg = ""
		fmt.Scanln(&msg)
	}
}

func (this * Client)SelectUsers(){
	sendMsg := "who\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (this * Client) PrivateMode(){
	var msg string
	var towho string

	this.SelectUsers()
	fmt.Println("请输入你私聊的对象,如若退出请输入exit..........")
	fmt.Scanln(&towho)
	for towho != "exit"{

		fmt.Println("请输入你需要发送的消息,如若退出请输入exit..........")
		fmt.Scanln(&msg)
		for msg != "exit"{
			if len(msg) != 0{
				sendmsg := "to|" + towho + "|" + msg + "\n\n"
				this.conn.Write([]byte(sendmsg))
			}
			msg = ""
			fmt.Println("请输入你需要发送的消息,如若退出请输入exit..........")
			fmt.Scanln(&msg)
		}

		this.SelectUsers()
		fmt.Println("请输入你私聊的对象,如若退出请输入exit..........")
		fmt.Scanln(&towho)		
	}

}
func (this * Client) Rename()bool{
	fmt.Println(">>>>请输入用户名:")
	fmt.Scanln(&this.name)

	sendMsg := "rename|" + this.name + "\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (this *Client) Run(){
	for this.flag != 0{
		for this.menu() != true{
		}
		switch this.flag {
		case 1:
			this.PushlicMode()
			break
		case 2:
			this.PrivateMode()
			break
		case 3:
			this.Rename()
			break
		}
	}
}

var IP string
var PORT int

func init(){
	flag.StringVar(&IP,"i","127.0.0.1","设置目标IP(默认127.0.0.1)")
	flag.IntVar(&PORT,"p",9999,"设置目标端口(默认9999)")
}

func main(){
	flag.Parse()

	client := NewClient(IP,PORT)

	if client == nil{
		fmt.Println("NewClient error")
		return
	}

	go client.Respond()

	fmt.Println(">>>>>链接服务器成功...")
	client.Run()
}
