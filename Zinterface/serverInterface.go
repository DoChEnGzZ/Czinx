package Zinterface

type ConnectionFunc func(i ConnectionI)

type ServerI interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32,routerInterface2 RouterI)
	GetManager()ManagerI
	CallAfterConnect(ConnectionI)
	CallBeforeConnect(ConnectionI)
	CallBeforeStop(ConnectionI)
	SetAfterConnect(ConnectionFunc)
	SetBeforeConnect(ConnectionFunc)
	SetBeforeStop(ConnectionFunc)
}

