package utils

import (
	"encoding/json"
	"github.com/DoChEnGzZ/Czinx/Zinterface"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
)

//
// Config
// @Description: zinx的全局配置
//
type Config struct {
	/*Server*/
	Server Zinterface.ServerI
	Host string
	Port int
	Name string
	/*CZinx*/
	Version string //zinx版本
	MaxPackageSize int //最大包长
	MaxConn int //最大连接数
	MaxWorkPoolSize int //最大工作池数
	MaxPoolTaskSize int //每个池的最大任务数
	MaxBuffChanSize int //服务器最大的发送缓冲区
}

var GlobalConfig *Config

var configPath="./config/config.json"

func (c *Config)loadFromJson(){
	zap.L().Info("Read config from %s"+configPath)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		zap.L().Error("Config Load error"+err.Error())
		return
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		zap.L().Error("Config Load error"+err.Error())
		return
	}
}

func InitConfig(){
	GlobalConfig=&Config{
		Host:           "0.0.0.0",
		Port:           8080,
		Name:           "CZinx",
		Version:        "v0.4",
		MaxPackageSize: 512,
		MaxConn:        1024,
		MaxWorkPoolSize: 1,
		MaxPoolTaskSize: 1,
		MaxBuffChanSize: 1024,
	}
	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config change in running")
		err := viper.Unmarshal(&GlobalConfig)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	})
	err := viper.ReadInConfig()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	err=viper.Unmarshal(&GlobalConfig)
	if err!=nil{
		zap.L().Error(err.Error())
		return
	}
	zap.L().Info("init config success")
	//GlobalConfig.loadFromJson()
}
