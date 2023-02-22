package Znet

import (
	"Czinx/Zinterface"
	"log"
	"net"
)

//实现路由时，此为基础中间件
type BaseRouter struct {

}

func (r BaseRouter) PreHandle(request Zinterface.RequestInterface)  {

}
func (r BaseRouter) Handle(request Zinterface.RequestInterface)  {
	err := CallBackFunc(request.GetConnection().GetTcpConnection(),request.GetData())
	if err != nil {
		return
	}
}

func (r BaseRouter) PostHandle(request Zinterface.RequestInterface)  {

}

func CallBackFunc(conn *net.TCPConn,buf []byte)error{
	log.SetPrefix("[HandleApi:CallBackFunc]")
	log.Printf("HandleApi start")
	if _,err:=conn.Write(buf);err!=nil{
		return err
	}
	return nil
}
