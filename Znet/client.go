package Znet

import (
	"Czinx/Zinterface"
	"fmt"
	"log"
	"net"
	"sync"
)
const(
	DefaultIp="127.0.0.1"
	DefaultPort=8080
)


type Client struct {
	ip string
	port int
	msgHandler *MsgHandler
	conn Zinterface.ConnectionI
	exitChan chan struct{}
	mutex sync.Mutex
	beforeStart Zinterface.ConnectionFunc
	afterStart Zinterface.ConnectionFunc
	beforeStop Zinterface.ConnectionFunc
	isOkChan chan bool
}

func NewClient(ip string,port int)*Client  {
	return &Client{
		ip:         ip,
		port:       port,
		msgHandler: NewMsgHandler(1, 1),
		exitChan:   make(chan struct{}),
		mutex: sync.Mutex{},
		isOkChan: make(chan bool),
	}
}

func (c *Client)Start(){
	//启动客户端
	go func() {
		addr:=&net.TCPAddr{
			IP: net.ParseIP(c.ip),
			Port: c.port,
			Zone: "",
		}
		conn, err := net.DialTCP("tcp",nil,addr)
		if err != nil {
			log.Printf("[Client]start error"+err.Error())
		}
		c.conn=NewConnection(c,conn,0,c.msgHandler)
		c.isOkChan<-true
		go c.conn.Start()
		select {
		case <-c.exitChan:
			c.conn.Stop()
		}
	}()
}
func (c *Client) Stop(){
	log.Println("[Client]Stop")
	c.exitChan<- struct{}{}
}

func (c *Client) Serve()  {
	//
}
func (c *Client) Conn()Zinterface.ConnectionI{
	return c.conn
}
func (c *Client)GetMsgHandle()Zinterface.MsgHandleI{
	return c.msgHandler
}

func (c *Client) AddRouter(msgId uint32,r Zinterface.RouterI){
	err := c.msgHandler.AddRouter(msgId, r)
	if err != nil {
		fmt.Println("[Client]add router error")
	}
}

func (c *Client) SendMessage(id uint32, s string) {
	fmt.Printf("id:%d,s:%s\n",id,s)
	//等待连接建立完成
	select {
	case <-c.isOkChan:
		break
	}
	err := c.Conn().Send(id, []byte(s))
	if err != nil {
		log.Println("Client")
	}
}

func (c *Client) GetManager()Zinterface.ManagerI{
	return nil
}
func (c *Client) CallAfterConnect(i Zinterface.ConnectionI){
	if c.afterStart!=nil{
		c.afterStart(i)
	}
}
func (c *Client)CallBeforeConnect(i Zinterface.ConnectionI){
	if c.beforeStart!=nil{
		c.beforeStart(i)
	}
}
func (c *Client)CallBeforeStop(i Zinterface.ConnectionI){
	if c.beforeStop!=nil{
		c.beforeStop(i)
	}
}
func (c *Client)SetAfterConnect(f Zinterface.ConnectionFunc){
	c.afterStart=f
}
func (c *Client)SetBeforeConnect(f Zinterface.ConnectionFunc){
	c.beforeStart=f
}
func (c *Client)SetBeforeStop(f Zinterface.ConnectionFunc){
	c.beforeStop=f
}





