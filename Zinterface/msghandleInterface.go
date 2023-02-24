package Zinterface

type MsgHandleI interface {
	AddRouter(msgId uint32,r RouterI)error
	Handle(i RequestI)error
	SendMessage(r RequestI)
}
