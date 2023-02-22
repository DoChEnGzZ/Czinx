package Znet

import "Czinx/Zinterface"

type Request struct {
	conn Zinterface.ConnectionInterface
	data []byte
}

func (r *Request) GetConnection()Zinterface.ConnectionInterface  {
	return r.conn
}

func (r *Request) GetData()[]byte  {
	return r.data
}
