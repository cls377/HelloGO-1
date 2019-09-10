package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := "10010"
	Start(port)
}

func Start(port string) {
	host := ":" + port
	//获取TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed:%v\n", err)
		return
	}
	//	监听
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed:%v\n", err)
		return
	}
	//建立连接池
	conns := make(map[string]net.Conn)
	//消息通道
	messageChan := make(chan string, 10)
	//广播消息
	go BroadMessage(&conns, messageChan)
	//启动
	for {
		fmt.Printf("listening port %s ...\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}
		//把每个客户端连接扔入连接池
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)
		go Handle(conn, &conns, messageChan)
	}
}
func Handle(conn net.Conn, conns *map[string]net.Conn, message chan string) {
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		//把消息放入信道
		recvStr := string(buf[0:length])
		message <- recvStr
	}
}

func BroadMessage(conns *map[string]net.Conn, message chan string) {
	for {
		//从信道中读取信息
		msg := <-message
		fmt.Println(msg)
		//将消息进行广播
		for k, conn := range *conns {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("borad message to %s failed:%v\n", k, err)
				delete(*conns, k)
			}

		}
	}
}
