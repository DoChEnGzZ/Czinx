package Znet

import (
	"Czinx/Zinterface"
	"Czinx/utils"
	"context"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"time"
)

//根据收到的不同的消息id调用不同的Router进行处理
type MsgHandler struct {
	MsgRouterMap map[uint32]Zinterface.RouterI //每个消息id对应一个Router
	MaxPoolSize int
	TaskQueue []chan Zinterface.RequestI //[workerId]chan,通过向chan中放request找到对应的worker
	bucket *rate.Limiter //令牌桶
}

func NewMsgHandler()*MsgHandler{
	return &MsgHandler{
		MsgRouterMap: make(map[uint32]Zinterface.RouterI),
		MaxPoolSize: utils.GlobalConfig.MaxWorkPoolSize,
		TaskQueue:    make([]chan Zinterface.RequestI,utils.GlobalConfig.MaxWorkPoolSize),
		bucket: rate.NewLimiter(rate.Every(100*time.Millisecond),utils.GlobalConfig.MaxWorkPoolSize),
	}
}

func (h *MsgHandler) StartWorker(ctx context.Context,workerId int,taskChan chan Zinterface.RequestI)  {
	log.Printf("WorkerPool:%d is startting",workerId)
	for{
		select {
		case r:=<-taskChan:
			if err:=h.bucket.Wait(ctx);err!=nil{
				log.Printf("WorkerPool:%d handle error:%s",workerId,err)
				continue
			}
			err := h.Handle(r)
			if err != nil {
				log.Printf("WorkerPool:%d handle error:%s",workerId,err)
				continue
			}
		case <-ctx.Done():
			break
		}
	}
}

func (h *MsgHandler) StartWorkerPool(ctx context.Context){
	for i:=0;i<h.MaxPoolSize;i++{
		if h.TaskQueue[i]!=nil{
			log.Printf("workerpool %d already exited",i)
			return
		}
		h.TaskQueue[i]=make(chan Zinterface.RequestI,utils.GlobalConfig.MaxPoolTaskSize)
		go h.StartWorker(ctx,i,h.TaskQueue[i])
	}
}
func (h *MsgHandler) SendMessage(r Zinterface.RequestI){
	workerId:=int(r.GetConnection().GetConnID())%h.MaxPoolSize
	log.Printf("ConnID=%d,MsgID=%d send msg with WorkerId=%d in worker pool",r.GetConnection().GetConnID(),
		r.GetMessageID(),workerId)
	h.TaskQueue[workerId]<-r
}

func (h *MsgHandler) AddRouter(id uint32,r Zinterface.RouterI)error{
	if _,ok:=h.MsgRouterMap[id];ok{
		return errors.New(fmt.Sprintf("Msgid= %d,Router has been registed",id))
	}
	fmt.Printf("Msgid= %d,Router regist",id)
	h.MsgRouterMap[id]=r
	return nil
}

func (h *MsgHandler) Handle(r Zinterface.RequestI)error  {
	handle,ok:=h.MsgRouterMap[r.GetMessageID()]
	if !ok{
		return errors.New(fmt.Sprintf("Msgid= %d,Router has been registed",r.GetMessageID()))
	}
	handle.PreHandle(r)
	handle.Handle(r)
	handle.PostHandle(r)
	return nil
}

func (h *MsgHandler) Close()  {
	//todo:close all handler
}