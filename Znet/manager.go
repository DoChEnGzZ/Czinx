package Znet

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

type Manager struct {
	connectionPool map[uint64]Zinterface.ConnectionI
	lock sync.RWMutex
	lruLock sync.RWMutex
	length int
	MaxConn   int
	ConnCount int
	lru       *list.List
	connE map[uint64]*list.Element
}

func NewManager()*Manager{
	return &Manager{
		connectionPool: make(map[uint64]Zinterface.ConnectionI),
		lock:           sync.RWMutex{},
		length:         0,
		ConnCount:      0,
		MaxConn:        utils.GlobalConfig.MaxConn,
		lru:            list.New(),
		connE:          make(map[uint64]*list.Element),
	}
}

func (m *Manager) Add(c Zinterface.ConnectionI)error  {
	m.lock.Lock()
	m.lruLock.Lock()
	defer m.lock.Unlock()
	m.connectionPool[c.GetConnID()]=c
	//加入链表尾部
	ele:=m.lru.PushBack(c.GetConnID())
	m.connE[c.GetConnID()]=ele
	m.lruLock.Unlock()
	m.length++
	m.ConnCount++
	zap.L().Debug(fmt.Sprintf("manager add,length:%d",m.length))
	zap.L().Info(fmt.Sprintf("[Manager]connection:%d add to manager",c.GetConnID()))
	for m.length>m.MaxConn{
		m.UpdateLru()
		return nil
		//return errors.New(fmt.Sprintf("connection pool is full,max connnums is %d",m.MaxConn))
	}
	return nil
}

func (m *Manager) Get(id uint64)(Zinterface.ConnectionI,error)  {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if c,ok:=m.connectionPool[id];ok{
		err := m.UseLru(id)
		if err != nil {
			return nil, err
		}
		return c,nil
	}else {
		return nil,errors.New("[Manager]no id="+strconv.Itoa(int(id))+"connection in manager")
	}
}

func (m *Manager) Remove(id uint64)error  {
	//m.lock.RLock()
	//m.lruLock.RLock()
	//defer m.lruLock.RUnlock()
	//defer m.lock.RUnlock()
	if c,ok:=m.connectionPool[id];ok{
		c.Stop()
		delete(m.connectionPool,id)
		ele:=m.connE[id]
		delete(m.connE,id)
		m.lru.Remove(ele)
		m.length--
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
	//清空connection manager
	for connId,conn:=range m.connectionPool{
		conn.Stop()
		delete(m.connectionPool,connId)
	}
	//清空lru
	for id,ele:=range m.connE{
		delete(m.connE,id)
		m.lru.Remove(ele)
	}
	m.length=0
	return nil
}

func (m *Manager) Size()int  {
	return len(m.connectionPool)
}

func (m *Manager) UpdateLru()  {
	m.lruLock.RLock()
	defer m.lruLock.RUnlock()
	ele:=m.lru.Front()
	frontId:=ele.Value.(uint64)
	zap.L().Debug(fmt.Sprintf("remove conn id :%d",frontId))
	conn,ok:=m.connectionPool[frontId]
	if !ok{
		zap.L().Error(fmt.Sprintf("conn id: %d not in conn pool",frontId))
	}
	conn.Stop()
	zap.L().Debug("after stop")
	delete(m.connectionPool,frontId)
	delete(m.connE,frontId)
	m.lru.Remove(ele)
	m.length--
	zap.L().Debug(fmt.Sprintf("conn pool size:%d",m.length))
}

func (m *Manager) UseLru(id uint64)error  {
	m.lruLock.Lock()
	defer m.lruLock.Unlock()
	ele,ok:=m.connE[id]
	if!ok{
		return errors.New(fmt.Sprintf("no element int manager with connid:%d",id))
	}
	m.lru.MoveToBack(ele)
	return nil
}

func (m *Manager) GetCount()int  {
	return m.ConnCount
}