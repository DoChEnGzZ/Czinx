package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const logPath="./Logs/log.txt"
var lg *zap.Logger

func InitLogger(){
	encoder:=getEncoder()
	consoleEncoder:=zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	writerSyncer:=getWriterSyncer()
	core:=zapcore.NewTee(
		zapcore.NewCore(encoder,writerSyncer,zap.ErrorLevel),
		zapcore.NewCore(consoleEncoder,zapcore.AddSync(os.Stdout),zap.DebugLevel),
		)
	lg:=zap.New(core,zap.AddCaller())
	zap.ReplaceGlobals(lg)
	zap.L().Info("Zap Init success")
}

// core 三个参数之  日志输出路径
func getWriterSyncer() zapcore.WriteSyncer {
	//file, _ := os.Create("./server/zaplog_test/log.log")
	//或者将上面的NewMultiWriteSyncer放到这里来，进行返回
	//return zapcore.AddSync(file)

	//引入第三方库 Lumberjack 加入日志切割功能
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   "./log/log.txt",
		MaxSize:    10, // megabytes
		MaxBackups: 32,
		MaxAge:     14,    // days
		Compress:   false, //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)
}


func getEncoder()zapcore.Encoder{
	c:=zap.NewProductionEncoderConfig()
	c.EncodeTime = zapcore.ISO8601TimeEncoder
	c.TimeKey = "time"
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	c.EncodeDuration = zapcore.SecondsDurationEncoder
	c.EncodeCaller = zapcore.ShortCallerEncoder
	encoder:=zapcore.NewConsoleEncoder(c)
	return encoder
}