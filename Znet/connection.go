package Znet

import (
	"Czinx/Zinterface"
	"log"
	"net"
)

type Connection struct {
	Conn *net.TCPConn
	ConnID uint32
	IsClosed bool
	HandleApi Zinterface.HandleFunc
	//告知当前连接以及停止
	StopChan chan bool
}

func NewConnection(conn *net.TCPConn,coonId uint32,handleFunc Zinterface.HandleFunc)*Connection{
	return &Connection{
		Conn:      conn,
		ConnID:    coonId,
		IsClosed:  false,
		HandleApi: handleFunc,
		StopChan:  make(chan bool,1),
	}
}

//启动读写业务
func (c *Connection) StartReader(){
	log.Printf("Reader GoRoutine is running...")
	defer c.Stop()
	defer log.Printf("ConnID=%d RemoteAddr=%s stop reading",c.GetConnID(),c.GetRemoteAddr())

	for{
		buf:=make([]byte,512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			log.Printf("Read error=%s",err.Error())
			continue
		}
		//完成读之后让APi处理
		if err:=c.HandleApi(c.Conn,buf,cnt);err!=nil{
			log.Printf("connID=%d, handle is error", c.ConnID)
			c.StopChan<-true
			break
		}
	}
}

func (c *Connection) Start()  {
	log.SetPrefix("[Start]")
	c.StartReader()
	for{
		select {
		case <-c.StopChan:
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

func (c *Connection) Send(data []byte)error  {
	_, err := c.Conn.Write(data)
	return err
}

func (c *Connection) GetConnID()uint32  {
	return c.ConnID
}



