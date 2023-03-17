package Zinterface

import "net"

type ConnectionI interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//获取套接字
	GetTcpConnection()*net.TCPConn
	//获取连接ID
	GetConnID() uint64
	//获取远程客户端的地址
	GetRemoteAddr()net.Addr
	GetManager()ManagerI
	//发送数据
	Send(messageId uint32,data []byte)error
	SendBuff(messageId uint32,data []byte)error
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

//处理连接业务的方法
type HandleFunc func(*net.TCPConn,[]byte,int)error