package Zinterface

type ServerI interface {
	Start()
	Stop()
	Serve()
	AddRouter(routerInterface2 RouterI)
}

