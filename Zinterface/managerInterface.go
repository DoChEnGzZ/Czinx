package Zinterface

type ManagerI interface {
	Add(i ConnectionI)error
	Get(id uint64)(ConnectionI,error)
	GetCount()int
	UseLru(id uint64)error
	Remove(id uint64)error
	Size()int
	Clear()error
}
