package Znet

import (
	"context"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
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
	var ctx context.Context
	ctx,s.cancel =context.WithCancel(context.Background())
	go func() {
		//start pool
		s.Handler.StartWorkerPool(ctx)
	//1 获取本服务器的ip地址
	addr, err := net.ResolveTCPAddr(s.ipVersion,fmt.Sprintf("%s:%d",s.ipAddress,s.Port))
	if err != nil {
		zap.L().Error(fmt.Sprintf("server %s init tcp,ip addr error:%s",s.Name,err.Error()))
		return 
	}
	//2 监听给出的ip地址和端口
	listener, err := net.ListenTCP(s.ipVersion,addr)
	if err != nil {
		zap.L().Error(fmt.Sprintf("%s listen error:%s",s.Name,err.Error()))
		return
	}
		zap.L().Debug(fmt.Sprintf("%s start finished,now is listening......",s.Name))
	//3 阻塞等待客户端连接，处理客户端连接业务
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%s connection failed:%s",s.Name,err.Error()))
			continue
		}
		//if s.manager.Size()>utils.GlobalConfig.MaxConn{
		//	zap.L().Error(fmt.Sprintf("manager connPool is full",err))
		//	err := conn.Close()
		//	if err != nil {
		//		zap.L().Error(err.Error())
		//	}
		//	continue
		//}
		zap.L().Info(fmt.Sprintf("server connection established with %s",conn.RemoteAddr().String()))
		//完成connection的注册，此时将方法传入此
		snowId,err:=utils.GetId()
		if err!=nil{
			zap.L().Error(err.Error())
			continue
		}
		c:=NewConnection(s,conn, snowId,s.Handler)
		if err:=s.manager.Add(c);err!=nil{
			zap.L().Error(err.Error())
			c.Stop()
			continue
		}
		go c.Start()
	}
	}()
}

func (s *Server) Stop()  {
	err := s.manager.Clear()
	s.cancel()
	if err != nil {
		zap.L().Error("server stop error:"+err.Error())
	}
	zap.L().Error("server stop error:"+err.Error())
}

func (s *Server) Serve()  {
	s.Start()
	select {
	}
}

func (s *Server) AddRouter(msgId uint32,router Zinterface.RouterI)  {
	err := s.Handler.AddRouter(msgId, router)
	if err != nil {
		zap.L().Error("add router error:"+err.Error())
		return
	}
}

func NewServer(name string) Zinterface.ServerI {
	utils.InitLogger()
	utils.Init(1)
	zap.L().Info("server "+utils.GlobalConfig.Name+"is creating")
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
		zap.L().Info("[Connect]:call hook func after connect")
		s.afterStart(c)
	}
}
func (s *Server)CallBeforeConnect(c Zinterface.ConnectionI){
	if s.beforeStart!=nil{
		zap.L().Info("[Connect]:call hook func after connect")
		s.beforeStart(c)
	}
}
func (s *Server)CallBeforeStop(c Zinterface.ConnectionI){
	if s.beforeStop!=nil{
		zap.L().Info("[Connect]:call hook func after connect")
		s.beforeStop(c)
	}
}
func (s *Server)SetAfterConnect(f Zinterface.ConnectionFunc){
	s.afterStart=f
}
func (s *Server)SetBeforeConnect(f Zinterface.ConnectionFunc){
	//zap.L().Debug("set before connect")
	s.beforeStart=f
}
func (s *Server)SetBeforeStop(f Zinterface.ConnectionFunc){
	s.beforeStop=f
}
