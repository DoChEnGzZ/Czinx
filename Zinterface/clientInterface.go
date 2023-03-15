package Zinterface

type ClientI interface {
	Start()
	Stop()
	Conn()ConnectionI
	GetMsgHandle()MsgHandleI
	AddRouter(uint32,RouterI)
	//SendMessage(uint32,string)
	//GetManager()
}