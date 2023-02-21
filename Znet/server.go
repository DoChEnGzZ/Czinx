package Znet

import (
	"Czinx/Zinterface"
	"fmt"
	"log"
	"net"
)

type Server struct {
	//服务器名
	Name string
	//ip协议版本
	ipVersion string
	//ip地址
	ipAddress string
	//端口号
	Port int
}

func (s *Server) Start()  {
	log.SetPrefix("[start]")
	log.Printf("%s is starting on %s:%d",s.Name,s.ipAddress,s.Port)
	go func() {
	//1 获取本服务器的ip地址
	addr, err := net.ResolveTCPAddr(s.ipVersion,fmt.Sprintf("%s:%d",s.ipAddress,s.Port))
	if err != nil {
		log.Printf("%s get ip addr error:%s",s.Name,err.Error())
		return 
	}
	//2 监听给出的ip地址和端口
	listener, err := net.ListenTCP(s.ipVersion,addr)
	if err != nil {
		log.Printf("%s listen error:%s",s.Name,err.Error())
		return
	}
	log.Printf("%s start finished,now is listening......",s.Name)
	//3 阻塞等待客户端连接，处理客户端连接业务
	for connID:=0;;connID++{
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("%s connection failed:%s",s.Name,err.Error())
			continue
		}
		log.Printf("server connection established with %s",conn.RemoteAddr().String())
		//完成connection的注册，此时将方法传入此
		c:=NewConnection(conn, uint32(connID),CallBackFunc)
		go c.Start()
	}
	}()
}

func (s *Server) Stop()  {

}

func (s *Server) Serve()  {
	s.Start()
	select {
	}
}

func NewServer(name string) Zinterface.ServerInterface {
	s:=&Server{
		Name:      name,
		ipVersion: "tcp4",
		ipAddress: "0.0.0.0",
		Port:      8080,
	}
	return s
}

//模拟一个512字节的回写功能,即将发送来的功能回送回去
func CallBackFunc(conn *net.TCPConn,buf []byte,cnt int)error{
	log.SetPrefix("[HandleApi:CallBackFunc]")
	log.Printf("HandleApi start")
	if _,err:=conn.Write(buf[:cnt]);err!=nil{
		return err
	}
	return nil
}