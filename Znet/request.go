package Znet

import "Czinx/Zinterface"

type Request struct {
	conn Zinterface.ConnectionI
	message Zinterface.MessageI
}

func (r *Request) GetConnection()Zinterface.ConnectionI {
	return r.conn
}

func (r *Request) GetData()[]byte {
	return r.message.GetData()
}

func (r *Request) GetMessageID()uint32  {
	return r.message.GetMessageId()
}
