package Znet

import (
	"Czinx/Zinterface"
	"log"
)

//实现路由时，此为基础中间件
type BaseRouter struct {

}

func (r BaseRouter) PreHandle(request Zinterface.RequestI)  {

}
func (r BaseRouter) Handle(request Zinterface.RequestI)  {
	err := request.GetConnection().Send(1, request.GetData())
	if err != nil {
		log.Println("[Basic Router]handle error:",err)
		return
	}
}

func (r BaseRouter) PostHandle(request Zinterface.RequestI)  {

}


