package Zinterface

type ManagerI interface {
	Add(i ConnectionI)error
	Get(id uint32)(ConnectionI,error)
	Remove(id uint32)error
	Size()int
	Clear()error
}
