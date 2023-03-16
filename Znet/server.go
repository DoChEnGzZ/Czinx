package Znet

import (
	"context"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
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
	//路由与消息
	Handler *MsgHandler
	manager *Manager
	beforeStart Zinterface.ConnectionFunc
	afterStart Zinterface.ConnectionFunc
	beforeStop Zinterface.ConnectionFunc
	cancel     context.CancelFunc
}

func (s *Server) Start()  {
	//log.Printf("%s is starting on %s:%d",s.Name,s.ipAddress,s.Port)
	log.SetFlags(log.Ldate|log.Ltime|log.Llongfile)
	var ctx context.Context
	ctx,s.cancel =context.WithCancel(context.Background())
	go func() {
		//start pool
		s.Handler.StartWorkerPool(ctx)
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
		if s.manager.Size()>utils.GlobalConfig.MaxConn{
			err := conn.Close()
			if err != nil {
				log.Println("[Server]manager connPool is full",err)
			}
			continue
		}
		log.Printf("server connection established with %s",conn.RemoteAddr().String())
		//完成connection的注册，此时将方法传入此
		c:=NewConnection(s,conn, uint32(connID),s.Handler)
		go c.Start()
	}
	}()
}

func (s *Server) Stop()  {
	err := s.manager.Clear()
	s.cancel()
	if err != nil {
		log.Println("[Server stop]error:",err)
	}
	log.Println("[server stop]")
}

func (s *Server) Serve()  {
	s.Start()
	select {
	}
}

func (s *Server) AddRouter(msgId uint32,router Zinterface.RouterI)  {
	log.Println("Add router")
	err := s.Handler.AddRouter(msgId, router)
	if err != nil {
		fmt.Println("add router error:",err)
		return
	}
}

func NewServer(name string) Zinterface.ServerI {
	s:=&Server{
		Name:      utils.GlobalConfig.Name,
		ipVersion: "tcp4",
		ipAddress: utils.GlobalConfig.Host,
		Port:      utils.GlobalConfig.Port,
		Handler: NewMsgHandlerByConfig(),
		manager: NewManager(),
	}
	return s
}

func (s *Server) GetManager()Zinterface.ManagerI  {
	return s.manager
}

func (s *Server)CallAfterConnect(c Zinterface.ConnectionI){
	if s.afterStart!=nil{
		log.Println("[Connect]:call user hook func after connect")
		s.afterStart(c)
	}
}
func (s *Server)CallBeforeConnect(c Zinterface.ConnectionI){
	if s.beforeStart!=nil{
		log.Println("[Connect]:call user hook func after connect")
		s.beforeStart(c)
	}
}
func (s *Server)CallBeforeStop(c Zinterface.ConnectionI){
	if s.beforeStop!=nil{
		log.Println("[Connect]:call user hook func after connect")
		s.beforeStop(c)
	}
}
func (s *Server)SetAfterConnect(f Zinterface.ConnectionFunc){
	s.afterStart=f
}
func (s *Server)SetBeforeConnect(f Zinterface.ConnectionFunc){
	s.beforeStart=f
}
func (s *Server)SetBeforeStop(f Zinterface.ConnectionFunc){
	s.beforeStop=f
}
