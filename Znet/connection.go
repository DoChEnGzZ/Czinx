package Znet

import (
	"Czinx/Zinterface"
	"Czinx/utils"
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
	Handler Zinterface.MsgHandleI
	//告知当前连接以及停止
	StopChan chan bool
	WriteChan chan []byte
}

func NewConnection(conn *net.TCPConn,coonId uint32,handler Zinterface.MsgHandleI)*Connection{
	return &Connection{
		Conn:      conn,
		ConnID:    coonId,
		IsClosed:  false,
		Handler: handler,
		StopChan:  make(chan bool,1),
		WriteChan: make(chan []byte),
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
		if utils.GlobalConfig.MaxWorkPoolSize>0{
			c.Handler.SendMessage(req)
		}
	}
}

func (c *Connection) StartWriter(){
	log.Println("[Writer Goroutine] is running")
	defer log.Println("[Writer Goroutine] is closing")
	for{
		select {
		//从WriteChan中读到数据并发送出去
		case data:=<-c.WriteChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				log.Printf("[Writer Goroutine] write data error:%s"+err.Error())
				return
			}
		//从stopChan中收到信号，关闭WriteRoutine
		case <-c.StopChan:
			return
		}
	}
}

func (c *Connection) Start()  {
	log.SetPrefix("[Server Start]")
	if c.IsClosed{
		log.Printf("%d connection is closed",c.ConnID)
		return
	}
	go c.StartReader()
	go c.StartWriter()
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
		log.Printf("Message Pack error:%s",err.Error())
		return err
	}
	c.WriteChan<-bytes
	return nil
}

func (c *Connection) GetConnID()uint32  {
	return c.ConnID
}



