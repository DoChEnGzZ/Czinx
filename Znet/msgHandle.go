package Znet

import (
	"context"
	"errors"
	"fmt"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/DoChEnGzZ/Czinx/utils"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

//根据收到的不同的消息id调用不同的Router进行处理
type MsgHandler struct {
	MsgRouterMap map[uint32]Zinterface.RouterI //每个消息id对应一个Router
	MaxPoolSize int
	TaskQueue []chan Zinterface.RequestI //[workerId]chan,通过向chan中放request找到对应的worker
	buckets []*rate.Limiter //令牌桶
}

func NewMsgHandler(maxPoolSize,maxTaskSize int)*MsgHandler{
	h:=&MsgHandler{
		MsgRouterMap: make(map[uint32]Zinterface.RouterI),
		MaxPoolSize: maxPoolSize,
		TaskQueue:    make([]chan Zinterface.RequestI,maxPoolSize),
		buckets: make([]*rate.Limiter,maxPoolSize),
	}
	return h
}

func NewMsgHandlerByConfig()*MsgHandler{
	h:=&MsgHandler{
		MsgRouterMap: make(map[uint32]Zinterface.RouterI),
		MaxPoolSize: utils.GlobalConfig.MaxWorkPoolSize,
		TaskQueue:    make([]chan Zinterface.RequestI,utils.GlobalConfig.MaxWorkPoolSize),
		buckets:  make([]*rate.Limiter,utils.GlobalConfig.MaxWorkPoolSize),
	}
	return h
}

func (h *MsgHandler) StartWorker(ctx context.Context,workerId int,taskChan chan Zinterface.RequestI,bucket *rate.Limiter)  {
	zap.L().Info(fmt.Sprintf("WorkerPool:%d is startting",workerId))
	for{
		select {
		case r:=<-taskChan:
			if err:=bucket.Wait(ctx);err!=nil{
				zap.L().Error(fmt.Sprintf("WorkerPool:%d handle error:%s",workerId,err))
				continue
			}
			err := h.Handle(r)
			if err != nil {
				zap.L().Error(fmt.Sprintf("WorkerPool:%d handle error:%s",workerId,err))
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
			zap.L().Error(fmt.Sprintf("workerpool %d already exited",i))
			return
		}
		h.TaskQueue[i]=make(chan Zinterface.RequestI,h.MaxPoolSize)
		h.buckets[i]=rate.NewLimiter(rate.Every(100*time.Millisecond),h.MaxPoolSize)
		go h.StartWorker(ctx,i,h.TaskQueue[i],h.buckets[i])
	}
}
func (h *MsgHandler) SendMessage(r Zinterface.RequestI){
	workerId:=int(r.GetConnection().GetConnID())%h.MaxPoolSize
	zap.L().Debug(fmt.Sprintf("connection:%d MsgID:%d, send msg by WorkerId=%d in worker pool",r.GetConnection().GetConnID(),
		r.GetMessageID(),workerId))
	h.TaskQueue[workerId]<-r
}

func (h *MsgHandler) AddRouter(id uint32,r Zinterface.RouterI)error{
	if _,ok:=h.MsgRouterMap[id];ok{
		return errors.New(fmt.Sprintf("Msgid= %d,Router has been registed",id))
	}
	//zap.L().Error(fmt.Sprintf("Msgid= %d,Router regist",id))
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