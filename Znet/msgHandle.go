package Znet

import (
	"Czinx/Zinterface"
	"errors"
	"fmt"
)

//根据收到的不同的消息id调用不同的Router进行处理
type MsgHandler struct {
	MsgRouterMap map[uint32]Zinterface.RouterI
}

func NewMsgHandler()*MsgHandler{
	return &MsgHandler{
		MsgRouterMap: make(map[uint32]Zinterface.RouterI),
	}
}

func (h *MsgHandler) AddRouter(id uint32,r Zinterface.RouterI)error{
	if _,ok:=h.MsgRouterMap[id];ok{
		return errors.New(fmt.Sprintf("Msgid= %d,Router has been registed",id))
	}
	fmt.Printf("Msgid= %d,Router regist",id)
	h.MsgRouterMap[id]=r
	return nil
}

func (h *MsgHandler) Handle(r Zinterface.RequestI)error  {
	handle,ok:=h.MsgRouterMap[r.GetMessageID()]
	if !ok{
		return errors.New(fmt.Sprintf("Msgid= %d,Router has been registed",r.GetMessageID()))
	}
	handle.PreHandle(r)
	handle.Handle(r)
	handle.PostHandle(r)
	return nil
}