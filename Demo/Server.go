package main

import (
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/Znet"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
)

func main() {
	server:=Znet.NewServer("ST")
	server.AddRouter(0,Znet.BaseRouter{})
	server.SetBeforeConnect(func(i Zinterface.ConnectionI) {
		zap.L().Info(fmt.Sprintf("server:%s on %s:%d,maxbufsize is %d maxconnection nums is %d," +
			"current connection nums is %d",
			utils.GlobalConfig.Name,utils.GlobalConfig.Host,
			utils.GlobalConfig.Port,utils.GlobalConfig.MaxPackageSize,
			utils.GlobalConfig.MaxConn,i.GetConnID()))
		i.SetProperty("name","ZinxV1.0")
	})
	//var msgId uint32
	//server.AddRouter(msgId, HandleRouter{})
	go server.Serve()
	select {
	}
}

type HandleRouter struct{}

func (HandleRouter) PreHandle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}

func (HandleRouter) Handle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}

func (HandleRouter) PostHandle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}
