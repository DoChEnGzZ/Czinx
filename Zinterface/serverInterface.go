package Zinterface

type ServerInterface interface {
	Start()
	Stop()
	Serve()
	AddRouter(routerInterface2 RouterInterface)
}

