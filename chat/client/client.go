package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	Start(os.Args[1])
}
func Start(tcpAddrStr string) {
	//1.根据输入的IP加端口生成TCP的地址信息
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addrr failed:%v\n", err)
		return
	}
	//2.向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed:%v\n", err)
		return
	}
	//向服务器发送消息
	go SendMsg(conn)
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			os.Exit(0)
			break
		}
		fmt.Println(string(buf[0:length]))
	}
}
func SendMsg(conn net.Conn) {
	//username := conn.LocalAddr().String()
	username := "欧阳:"
	for {
		var input string
		fmt.Scanln(&input)
		if input == "/q" || input == "/quit" {
			fmt.Println("Byebye...")
			conn.Close()
			os.Exit(0)
		}
		if len(input) > 0 {
			msg := username + input
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
