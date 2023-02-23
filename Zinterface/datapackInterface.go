package Zinterface

type DataPackI interface {
	GetHeadLen()uint32
	Pack(MessageI)([]byte,error)
	UnPack([]byte)(MessageI,error)
}
