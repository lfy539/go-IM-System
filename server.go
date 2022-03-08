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

	// 在线用户的列表
	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

// ListenMessage 监听message广播消息channel的goruntine ，一旦有消息就发送给全部的在线User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		// 将message 发给users
		s.mapLock.Lock()
		for _, cli := range s.OnLineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}
func (s *Server) Handler(conn net.Conn) {
	// 当前链接业务
	// fmt.Println("链接建立成功。。。。")
	// 用户上线，将用户加入到onlineMap中
	user := NewUser(conn, s)
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)
	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			// 提取用户的消息 （去除'\n'）
			msg := string(buf[:n-1])
			// 用户针对msg进行
			user.DoMessage(msg)
			// 用户的任意消息，代表当前用户是一个活跃的
			isLive <- true
		}
	}()
	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前是活跃的，应该重置定时器
			// 不做任何事情，更新定时器
		case <-time.After(time.Second * 120):
			// 已经超时
			//将当前的user强制关闭
			user.SendMsg("你被踢了")
			close(user.C)
			conn.Close()
			//退出当前handler
			return

		}
	}

}

// Start 启动服务器的接口
func (s *Server) Start() {
	// socket listen

	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net Listen err :", err)
		return
	}

	defer listener.Close()
	// 启动监听Message的goruntine
	go s.ListenMessage()
	for {
		// accept
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener accept err :", err)
			continue
		}
		fmt.Println(conn.RemoteAddr().String())
		//do handler
		go s.Handler(conn)
	}

	// close socket
}
