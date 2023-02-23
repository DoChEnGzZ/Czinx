package Zinterface

type RouterI interface {
	//处理conn业务前中后的方法
	PreHandle(requestInterface RequestI)
	Handle(requestInterface RequestI)
	PostHandle(requestInterface RequestI)
}
