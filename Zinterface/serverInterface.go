package Zinterface

type ServerI interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32,routerInterface2 RouterI)
}

