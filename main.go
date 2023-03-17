package main

import (
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/Znet"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
)

func main() {
	server:=Znet.NewServer()
	server.AddRouter(0,Znet.BaseRouter{})
	server.SetBeforeConnect(func(c Zinterface.ConnectionI) {
		zap.L().Info(fmt.Sprintf("server:%s on %s:%d,maxbufsize is %d maxconnection nums is %d," +
			"current connection id=%d ,connection nums is %d,total is %d",
			utils.GlobalConfig.Name,utils.GlobalConfig.Host,
			utils.GlobalConfig.Port,utils.GlobalConfig.MaxPackageSize,
			utils.GlobalConfig.MaxConn,c.GetConnID(),c.GetManager().Size(),
			c.GetManager().GetCount()) )
		c.SetProperty("name","ZinxV1.0")
	})
	server.SetAfterConnect(func(i Zinterface.ConnectionI) {
		i.Send(1,[]byte("hello! from Czinx"))
	})
	//var msgId uint32
	//server.AddRouter(msgId, HandleRouter{})
	go server.Serve()
	select {
	}

}
