package main

import (
	"Czinx/Zinterface"
	"Czinx/Znet"
	"Czinx/utils"
	"encoding/binary"
	_ "fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	s:=Znet.NewServer("test")
	s.AddRouter(0,Znet.NewBaseRouter("client 0 test message"))
	s.AddRouter(1,Znet.NewBaseRouter("高打低打傻逼啊彬子"))
	s.SetBeforeConnect(func(i Zinterface.ConnectionI) {
		log.Printf("server %s is starting on %s:%d,maxbufsize is %d maxconnection nums is %d," +
			"connection id is %d",
			utils.GlobalConfig.Name,utils.GlobalConfig.Host,
			utils.GlobalConfig.Port,utils.GlobalConfig.MaxPackageSize,
			utils.GlobalConfig.MaxConn,i.GetConnID())
		i.SetProperty("name","ZinxV1.0")
	})
	go ClientTest(0)
	go ClientTest(1)
	s.Serve()

}

func ClientTest(msgID uint32) {
	log.Println("[Client]test start after 3s")
	time.Sleep(3*time.Second)
	conn, err := net.Dial("tcp4","127.0.0.1:8080")
	if err != nil {
		log.Printf("[Client]connect error:%s",err.Error())
		return
	}
	log.Printf("[Client]client connection establishded with %s",conn.RemoteAddr().String())
	for i:=1;i<10;i++{
		//_, err := conn.Write([]byte("hello world!v0.2"))
		//if err != nil {
		//	log.Printf("write error:%s",err.Error())
		//	continue
		//}
		//buf :=make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf error ")
		//	return
		//}
		//向服务器发数据
		bytes, err := Znet.DefaultDataPack.Pack(Znet.NewMessage([]byte("hello world ZinxV0.5 with message"), msgID))
		if err != nil {
			log.Println("[Client]Pack error:",err)
			return
		}
		_, err = conn.Write(bytes)
		if err != nil {
			log.Printf("[Client]write error:%s",err.Error())
			return
		}
		//从服务器接收数据
		head:=make([]byte,Znet.DefaultDataPack.GetHeadLen())
		_, err = io.ReadFull(conn,head)
		if err != nil {
			log.Printf("[Client]Read error=%s",err.Error())
			return
		}
		dataLen:=binary.LittleEndian.Uint32(head[:4])
		data:=make([]byte,dataLen)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			log.Printf("[Client]Read error=%s",err.Error())
			continue
		}
		//将包头和数据合并
		buf:=append(head,data...)
		msg, err :=Znet.DefaultDataPack.UnPack(buf)
		if err!=nil{
			log.Printf("[Client]Read error=%s",err.Error())
			continue
		}
		log.Printf("[Client]server call back : %s\n", string(msg.GetData()))
		time.Sleep(1*time.Second)
	}
}
