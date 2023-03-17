package Znet

import (
	"errors"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type Manager struct {
	connectionPool map[uint32]Zinterface.ConnectionI
	lock sync.RWMutex
	connNums int
	MaxConn int
}

func NewManager()*Manager{
	return &Manager{
		connectionPool: make(map[uint32]Zinterface.ConnectionI),
		lock:           sync.RWMutex{},
		connNums: 0,
		MaxConn: utils.GlobalConfig.MaxConn,
	}
}

func (m *Manager) Add(c Zinterface.ConnectionI)error  {
	if m.connNums>m.MaxConn{
		return errors.New(fmt.Sprintf("connection pool is full,max connnums is %d",m.MaxConn))
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.connectionPool[c.GetConnID()]=c
	zap.L().Info(fmt.Sprintf("[Manager]connection:%d add to manager",c.GetConnID()))
	m.connNums++
	return nil
}

func (m *Manager) Get(id uint32)(Zinterface.ConnectionI,error)  {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if c,ok:=m.connectionPool[id];ok{
		return c,nil
	}else {
		return nil,errors.New("[Manager]no id="+strconv.Itoa(int(id))+"connection in manager")
	}
}

func (m *Manager) Remove(id uint32)error  {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if c,ok:=m.connectionPool[id];ok{
		c.Stop()
		delete(m.connectionPool,id)
		return nil
	}else {
		return errors.New("[Manager]no id="+strconv.Itoa(int(id))+"connection in manager")
	}
}

func (m *Manager) Clear()error  {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.Size()==0{
		return errors.New("[Manager]ConnectionPool is empty")
	}
	for connId,conn:=range m.connectionPool{
		conn.Stop()
		delete(m.connectionPool,connId)
	}
	return nil
}

func (m *Manager) Size()int  {
	return len(m.connectionPool)
}