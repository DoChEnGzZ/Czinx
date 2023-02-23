package Zinterface

type RequestI interface {
	//获取连接
	GetConnection() ConnectionI
	//获取数据
	GetData()[]byte

}
