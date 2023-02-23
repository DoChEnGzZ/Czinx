package Zinterface

type MessageI interface {
	GetData()[]byte
	GetMessageId() uint32
	GetDataLen() uint32
	SetData([]byte)
	SetMessageId(uint32)
	SetDataLen(uint32)
}
