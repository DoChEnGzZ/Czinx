package Znet

import (
	"Czinx/Zinterface"
	"Czinx/utils"
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
	//路由
	Router Zinterface.RouterI
}

func (s *Server) Start()  {
	log.SetPrefix("[server start]")
	//log.Printf("%s is starting on %s:%d",s.Name,s.ipAddress,s.Port)
	log.Printf("server %s is starting on %s:%d,maxbufsize is %d maxconnection nums is %d",
		utils.GlobalConfig.Name,utils.GlobalConfig.Host,utils.GlobalConfig.Port,utils.GlobalConfig.MaxPackageSize,utils.GlobalConfig.MaxConn)
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
		c:=NewConnection(conn, uint32(connID),s.Router)
		go c.Start()
	}
	}()
}

func (s *Server) Stop()  {
	log.Println("[server stop]")
}

func (s *Server) Serve()  {
	s.Start()
	select {
	}
}

func (s *Server) AddRouter(routerInterface Zinterface.RouterI)  {
	log.Println("Add router")
	s.Router=routerInterface
}

func NewServer(name string) Zinterface.ServerI {
	s:=&Server{
		Name:      utils.GlobalConfig.Name,
		ipVersion: "tcp4",
		ipAddress: utils.GlobalConfig.Host,
		Port:      utils.GlobalConfig.Port,
		Router: nil,
	}
	return s
}
