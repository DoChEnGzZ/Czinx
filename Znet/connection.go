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
	//HandleApi Zinterface.HandleFunc
	Router Zinterface.RouterInterface
	//告知当前连接以及停止
	StopChan chan bool
}

func NewConnection(conn *net.TCPConn,coonId uint32,router Zinterface.RouterInterface)*Connection{
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
		buf:=make([]byte,512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			log.Printf("Read error=%s",err.Error())
			continue
		}
		//完成读之后让APi处理
		//if err:=c.HandleApi(c.Conn,buf,cnt);err!=nil{
		//	log.Printf("connID=%d, handle is error", c.ConnID)
		//	c.StopChan<-true
		//	break
		//}
		req:=&Request{
			conn: c,
			data: buf[:cnt],
		}
		//理由路由绑定的handler执行
		go func(requestInterface Zinterface.RequestInterface) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)
	}
}

func (c *Connection) Start()  {
	log.SetPrefix("[Start]")
	if c.IsClosed{
		log.Printf("%d connection is closed",c.ConnID)
		return
	}
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



