package utils

import (
	"errors"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"time"
)

var(
	snowFlake *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineId()(uint16,error){
	return sonyMachineID,nil
}

func InitSnowFlake(machineID uint16)  {
	sonyMachineID=machineID
	t,_:=time.Parse("2006-01-02", "2023-01-01")
	settings:=sonyflake.Settings{
		StartTime:      t,
		MachineID:      getMachineId,
	}
	snowFlake=sonyflake.NewSonyflake(settings)
	zap.L().Info("sony flake init success")
}

func GetId()(uint64,error){
	if snowFlake==nil{
		return 0,errors.New("no snowflake init")
	}
	return snowFlake.NextID()
}
