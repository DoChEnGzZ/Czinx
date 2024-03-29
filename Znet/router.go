package Znet

import (
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"go.uber.org/zap"
)


//实现路由时，此为基础中间件
type BaseRouter struct {
	text string
}

func NewBaseRouter(text string)*BaseRouter{
	return &BaseRouter{text: text}
}

func (r BaseRouter) PreHandle(request Zinterface.RequestI)  {

}
func (r BaseRouter) Handle(request Zinterface.RequestI)  {
	zap.L().Info("[Server]receive"+string(request.GetData()))
	err := request.GetConnection().Send(3, []byte(r.text))
	if err != nil {
		zap.L().Error("[Basic Router]handle error:"+err.Error())
		return
	}
}

func (r BaseRouter) PostHandle(request Zinterface.RequestI)  {

}


