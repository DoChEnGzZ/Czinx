package Znet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
	"io"
	"net"
	"strconv"
	"sync"
)

type Connection struct {
	TcpServer Zinterface.ServerI
	//Client Zinterface.ClientI
	Conn *net.TCPConn
	ConnID uint64
	IsClosed bool
	//HandleApi Zinterface.HandleFunc
	Handler Zinterface.MsgHandleI
	//连接的配置和读写锁
	Property map[string]interface{}
	PropertyMutex sync.RWMutex
	//告知当前连接以及停止
	StopChan chan bool
	WriteChan chan []byte
	//有缓发送区
	WriteBufChan chan []byte
}

func NewConnection(server Zinterface.ServerI,conn *net.TCPConn,coonId uint64,handler Zinterface.MsgHandleI)*Connection{
	c:= &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    coonId,
		IsClosed:  false,
		Handler: handler,
		Property: make(map[string]interface{}),
		PropertyMutex: sync.RWMutex{},
		StopChan:  make(chan bool,1),
		WriteChan: make(chan []byte),
		WriteBufChan: make(chan []byte,utils.GlobalConfig.MaxPackageSize),
	}
	//if c.TcpServer.GetManager()!=nil{
	//	c.TcpServer.GetManager().Add(c)
	//}
	return c
}

//func NewClientConnection(client Zinterface.ClientI,conn *net.TCPConn,coonId uint64,handler Zinterface.MsgHandleI)*Connection{
//	c:= &Connection{
//		Client: client,
//		Conn:      conn,
//		ConnID:    coonId,
//		IsClosed:  false,
//		Handler: handler,
//		Property: make(map[string]interface{}),
//		PropertyMutex: sync.RWMutex{},
//		StopChan:  make(chan bool,1),
//		WriteChan: make(chan []byte),
//		WriteBufChan: make(chan []byte,utils.GlobalConfig.MaxPackageSize),
//	}
//	return c
//}

//启动读写业务
func (c *Connection) StartReader(){
	//defer c.Stop()
	//defer log.Printf("[Connection]ConnID=%d RemoteAddr=%s stop reading",c.GetConnID(),c.GetRemoteAddr())

	for{
		//获取包头，根据包头设计缓冲区
		head:=make([]byte,DefaultDataPack.GetHeadLen())
		_, err := io.ReadFull(c.GetTcpConnection(),head)
		if err != nil {
			//log.Printf("[Connection]Read error=%s,no msg since connected",err.Error())
			//zap.L().Debug(fmt.Sprintf("Read error=%s,no msg since connected",err.Error()))
			//c.StopChan<-true
			continue
		}
		//log.Printf("%d",head[4])
		//log.Printf("Recieve msg head is %d,%d",binary.LittleEndian.Uint32(head[:4]),binary.LittleEndian.Uint32(head[4:]))
		//获取数据长度
		dataLen:=binary.LittleEndian.Uint32(head[:4])
		data:=make([]byte,dataLen)
		_, err = io.ReadFull(c.GetTcpConnection(), data)
		if err != nil {
			zap.L().Error(fmt.Sprintf("[Connection]Read error=%s",err.Error()))
			//c.StopChan<-true
			continue
		}
		//将包头和数据合并
		buf:=append(head,data...)
		msg, err :=DefaultDataPack.UnPack(buf)
		if err!=nil{
			zap.L().Error(fmt.Sprintf("[Connection]Read error=%s",err.Error()))
			//c.StopChan<-true
			continue
		}
		//根据获取的数据构造request
		//cnt, err := c.Conn.Read(buf)
		req:=&Request{
			conn: c,
			message: msg,
		}
		if utils.GlobalConfig.MaxWorkPoolSize>0{
			_ = c.GetManager().UseLru(c.GetConnID())
			c.Handler.SendMessage(req)
		}
	}
}

func (c *Connection) StartWriter(){
	//log.Println("[Connection][Writer Goroutine] is running")
	defer zap.L().Debug("[Connection][Writer Goroutine] is closing")
	for{
		select {
		//从WriteChan中读到数据并发送出去
		case data:=<-c.WriteChan:
			_ = c.GetManager().UseLru(c.GetConnID())
			_, err := c.Conn.Write(data)
			if err != nil {
				zap.L().Error(fmt.Sprintf("[Writer Goroutine] write data error:%s"+err.Error()))
				return
			}
		case data,ok:=<-c.WriteBufChan:
			if !ok{
				zap.L().Error("[Writer Goroutine] write data error:WriteBufChan Closed")
				break
			}else {
				_ = c.GetManager().UseLru(c.GetConnID())
				_,err:=c.Conn.Write(data)
				if err != nil {
					zap.L().Error(fmt.Sprintf("[Connection][Writer Goroutine] write data error:%s"+err.Error()))
					return
			}
			}
		//从stopChan中收到信号，关闭WriteRoutine
		case <-c.StopChan:
			//zap.L().Debug("writer:receive stop signal from chan")
			return
		}
	}
}

func (c *Connection) Start()  {
	if c.IsClosed{
		zap.L().Error(fmt.Sprintf("[Connection]%d connection is closed",c.ConnID))
		return
	}
	//zap.L().Debug("call before start")
	c.TcpServer.CallBeforeConnect(c)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallAfterConnect(c)
	for{
		select {
		case <-c.StopChan:
			//zap.L().Info("[Connection]receive stop signal from chan")
			//c.Stop()
			return
		}
	}
}

func (c *Connection) Stop()  {
	//zap.L().Info(fmt.Sprintf("[Connection]ConnID=%d connection stop",c.ConnID))
	c.TcpServer.CallBeforeStop(c)
	if c.IsClosed{
		return
	}
	c.IsClosed=true
	err := c.Conn.Close()
	//c.Handler.Close()
	if err != nil {
		zap.L().Error(fmt.Sprintf("[Connection]Stop error=%s",err.Error()))
		return
	}
	c.StopChan<-true
	//if c.TcpServer.GetManager()!=nil{
	//	err = c.TcpServer.GetManager().Remove(c.ConnID)
	//}
	if err != nil {
		zap.L().Error("[Connection]stop error"+err.Error())
	}
	close(c.StopChan)
	close(c.WriteBufChan)
	close(c.WriteChan)
}

func (c *Connection) GetTcpConnection()*net.TCPConn  {
	return c.Conn
}

func (c *Connection) GetRemoteAddr()net.Addr  {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(messageId uint32,data []byte)error  {
	if c.IsClosed{
		return errors.New("[Connection]Send error,conn already closed")
	}
	msg:=NewMessage(data,messageId)
	bytes, err := DefaultDataPack.Pack(msg)
	//log.Println(string(bytes))
	zap.L().Debug("send:"+strconv.Itoa(int(messageId))+":"+string(bytes))
	if err != nil {
		zap.L().Error(fmt.Sprintf("[Connection]Message Pack error:%s",err.Error()))
		return err
	}
	c.WriteChan<-bytes
	return nil
}

func (c *Connection) SendBuff(messageId uint32,data []byte)error{
	if c.IsClosed{
		return errors.New("[Connection]Send error,conn already closed")
	}
	msg:=NewMessage(data,messageId)
	bytes,err:= DefaultDataPack.Pack(msg)
	if err != nil {
		zap.L().Error(fmt.Sprintf("[Connection]Message Pack error:%s",err.Error()))
		return err
	}
	c.WriteBufChan<-bytes
	return nil
}

func (c *Connection) GetManager()Zinterface.ManagerI  {
	return c.TcpServer.GetManager()
}

func (c *Connection) GetConnID()uint64  {
	return c.ConnID
}
//设置链接属性
func (c *Connection)SetProperty(key string, value interface{}){
	c.PropertyMutex.Lock()
	defer c.PropertyMutex.Unlock()
	c.Property[key]=value
	//log.Printf("[Connection]No.%dConnection add property key:%s string:%v",c.GetConnID(),key,value)
}
//获取链接属性
func (c *Connection)GetProperty(key string)(interface{}, error){
	c.PropertyMutex.RLock()
	defer c.PropertyMutex.RUnlock()
	if _,ok:=c.Property[key];!ok{
		return nil,errors.New("Connection]No."+strconv.Itoa(int(c.GetConnID()))+
			"Connection property key:%s not existed")
	}else {
		return c.Property[key],nil
	}
}
//移除链接属性
func (c *Connection)RemoveProperty(key string){
	c.PropertyMutex.Lock()
	defer c.PropertyMutex.Unlock()
	if _,ok:=c.Property[key];ok{
		delete(c.Property,key)
	}
}



