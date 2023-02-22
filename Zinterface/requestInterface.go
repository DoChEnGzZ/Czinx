package Zinterface

type RequestInterface interface {
	//获取连接
	GetConnection()ConnectionInterface
	//获取数据
	GetData()[]byte

}
