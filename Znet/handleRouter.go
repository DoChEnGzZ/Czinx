package Znet

import (
	"github.com/DoChEnGzZ/Czinx/Zinterface"
)

type HandleRouter struct{}

func (HandleRouter) PreHandle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}

func (HandleRouter) Handle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}

func (HandleRouter) PostHandle(requestInterface Zinterface.RequestI) {
	panic("implement me")
}

