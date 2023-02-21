package Znet

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	log.SetPrefix("[ClientTest]")
	log.Println("test start after 3s")
	time.Sleep(3*time.Second)
	conn, err := net.Dial("tcp4","127.0.0.1:8080")
	if err != nil {
		log.Printf("connect error:%s",err.Error())
		return
	}
	log.Printf("client connection establishded with %s",conn.RemoteAddr().String())
	for i:=0;i<10;i++{
		_, err := conn.Write([]byte("hello world!v0.2"))
		if err != nil {
			log.Printf("write error:%s",err.Error())
			continue
		}
		buf :=make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ")
			return
		}

		log.Printf(" server call back : %s, cnt = %d\n", buf[:cnt ],  cnt)

		time.Sleep(1*time.Second)
	}
}

func TestServer(t *testing.T) {
	s:=NewServer("test")
	go ClientTest()
	s.Serve()

}