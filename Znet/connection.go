package Znet

import (
	"Czinx/Zinterface"
	_ "Czinx/utils"
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Connection struct {
	Conn *net.TCPConn
	ConnID uint32
	IsClosed bool
	//HandleApi Zinterface.HandleFunc
	Router Zinterface.RouterI
	//告知当前连接以及停止
	StopChan chan bool
}

func NewConnection(conn *net.TCPConn,coonId uint32,router Zinterface.RouterI)*Connection{
	return &Connection{
		Conn:      conn,
		ConnID:    coonId,
		IsClosed:  false,
		Router: router,
		StopChan:  make(chan bool,1),
	}
}

//启动读写业务
func (c *Connection) StartReader(){
	log.Printf("Reader GoRoutine is running...")
	defer c.Stop()
	defer log.Printf("ConnID=%d RemoteAddr=%s stop reading",c.GetConnID(),c.GetRemoteAddr())

	for{
		//获取包头，根据包头设计缓冲区
		head:=make([]byte,DefaultDataPack.GetHeadLen())
		_, err := io.ReadFull(c.GetTcpConnection(),head)
		if err != nil {
			log.Printf("Read error=%s",err.Error())
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
			log.Printf("Read error=%s",err.Error())
			c.StopChan<-true
			continue
		}
		//将包头和数据合并
		buf:=append(head,data...)
		msg, err :=DefaultDataPack.UnPack(buf)
		if err!=nil{
			log.Printf("Read error=%s",err.Error())
			c.StopChan<-true
			continue
		}
		//根据获取的数据构造request
		//cnt, err := c.Conn.Read(buf)
		req:=&Request{
			conn: c,
			message: msg,
		}
		//理由路由绑定的handler执行
		go func(requestInterface Zinterface.RequestI) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)
	}
}

func (c *Connection) Start()  {
	log.SetPrefix("[Server Start]")
	if c.IsClosed{
		log.Printf("%d connection is closed",c.ConnID)
		return
	}
	c.StartReader()
	for{
		select {
		case <-c.StopChan:
			log.Printf("recieve stop signal from chan")
			return
		}
	}
}

func (c *Connection) Stop()  {
	log.SetPrefix("[Stop]")
	log.Printf("stop ConnID=%d",c.ConnID)
	if c.IsClosed{
		return
	}
	c.IsClosed=true
	err := c.Conn.Close()
	if err != nil {
		log.Printf("Stop error=%s",err.Error())
		return
	}
	c.StopChan<-true
	close(c.StopChan)
}

func (c *Connection) GetTcpConnection()*net.TCPConn  {
	return c.Conn
}

func (c *Connection) GetRemoteAddr()net.Addr  {
	return c.Conn.RemoteAddr()

}

func (c *Connection) Send(messageId uint32,data []byte)error  {
	msg:=NewMessage(data,messageId)
	bytes, err := DefaultDataPack.Pack(msg)
	if err != nil {
		return err
	}
	if _, err := c.Conn.Write(bytes);err!=nil{
		c.StopChan<-true
		return err
	}
	return nil
}

func (c *Connection) GetConnID()uint32  {
	return c.ConnID
}



