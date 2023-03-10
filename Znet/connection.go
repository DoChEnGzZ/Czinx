package Znet

import (
	"Czinx/Zinterface"
	"Czinx/utils"
	_ "Czinx/utils"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

type Connection struct {
	TcpServer Zinterface.ServerI
	Conn *net.TCPConn
	ConnID uint32
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

func NewConnection(server Zinterface.ServerI,conn *net.TCPConn,coonId uint32,handler Zinterface.MsgHandleI)*Connection{
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
	c.TcpServer.GetManager().Add(c)
	return c
}

//启动读写业务
func (c *Connection) StartReader(){
	log.Printf("[Connection]Reader GoRoutine is running...")
	defer c.Stop()
	defer log.Printf("[Connection]ConnID=%d RemoteAddr=%s stop reading",c.GetConnID(),c.GetRemoteAddr())

	for{
		//获取包头，根据包头设计缓冲区
		head:=make([]byte,DefaultDataPack.GetHeadLen())
		_, err := io.ReadFull(c.GetTcpConnection(),head)
		if err != nil {
			log.Printf("[Connection]Read error=%s",err.Error())
			c.StopChan<-true
			continue
		}
		//log.Printf("%d",head[4])
		//log.Printf("Recieve msg head is %d,%d",binary.LittleEndian.Uint32(head[:4]),binary.LittleEndian.Uint32(head[4:]))
		//获取数据长度
		dataLen:=binary.LittleEndian.Uint32(head[:4])
		data:=make([]byte,dataLen)
		_, err = io.ReadFull(c.GetTcpConnection(), data)
		if err != nil {
			log.Printf("[Connection]Read error=%s",err.Error())
			c.StopChan<-true
			continue
		}
		//将包头和数据合并
		buf:=append(head,data...)
		msg, err :=DefaultDataPack.UnPack(buf)
		if err!=nil{
			log.Printf("[Connection]Read error=%s",err.Error())
			c.StopChan<-true
			continue
		}
		//根据获取的数据构造request
		//cnt, err := c.Conn.Read(buf)
		req:=&Request{
			conn: c,
			message: msg,
		}
		if utils.GlobalConfig.MaxWorkPoolSize>0{
			c.Handler.SendMessage(req)
		}
	}
}

func (c *Connection) StartWriter(){
	log.Println("[Connection][Writer Goroutine] is running")
	defer log.Println("[Connection][Writer Goroutine] is closing")
	for{
		select {
		//从WriteChan中读到数据并发送出去
		case data:=<-c.WriteChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				log.Printf("[Connection][Writer Goroutine] write data error:%s"+err.Error())
				return
			}
		case data,ok:=<-c.WriteBufChan:
			if !ok{
				log.Printf("[Connection][Writer Goroutine] write data error:WriteBufChan Closed")
				break
			}else {
				_,err:=c.Conn.Write(data)
				if err != nil {
					log.Printf("[Connection][Writer Goroutine] write data error:%s"+err.Error())
					return
			}
			}
		//从stopChan中收到信号，关闭WriteRoutine
		case <-c.StopChan:
			return
		}
	}
}

func (c *Connection) Start()  {
	if c.IsClosed{
		log.Printf("[Connection]%d connection is closed",c.ConnID)
		return
	}
	c.TcpServer.CallBeforeConnect(c)
	go c.StartReader()
	go c.StartWriter()
	c.TcpServer.CallAfterConnect(c)
	for{
		select {
		case <-c.StopChan:
			log.Printf("[Connection]recieve stop signal from chan")
			return
		}
	}
}

func (c *Connection) Stop()  {
	log.Printf("[Connection]stop ConnID=%d",c.ConnID)
	c.TcpServer.CallBeforeStop(c)
	if c.IsClosed{
		return
	}
	c.IsClosed=true
	err := c.Conn.Close()
	c.Handler.Close()
	if err != nil {
		log.Printf("[Connection]Stop error=%s",err.Error())
		return
	}
	c.StopChan<-true
	err = c.TcpServer.GetManager().Remove(c.ConnID)
	if err != nil {
		log.Println("[Connection]stop error",err)
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
	if err != nil {
		log.Printf("[Connection]Message Pack error:%s",err.Error())
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
		log.Printf("[Connection]Message Pack error:%s",err.Error())
		return err
	}
	c.WriteBufChan<-bytes
	return nil
}

func (c *Connection) GetConnID()uint32  {
	return c.ConnID
}
//设置链接属性
func (c *Connection)SetProperty(key string, value interface{}){
	c.PropertyMutex.Lock()
	defer c.PropertyMutex.Unlock()
	c.Property[key]=value
	log.Printf("[Connection]No.%dConnection add property key:%s string:%v",c.GetConnID(),key,value)
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



