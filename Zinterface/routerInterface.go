package Zinterface

type RouterInterface interface {
	//处理conn业务前中后的方法
	PreHandle(requestInterface RequestInterface)
	Handle(requestInterface RequestInterface)
	PostHandle(requestInterface RequestInterface)
}
