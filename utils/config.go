package utils

import (
	"Czinx/Zinterface"
	"encoding/json"
	"io/ioutil"
	"log"
)

//
// Config
// @Description: zinx的全局配置
//
type Config struct {
	Server Zinterface.ServerI
	Host string
	Port int
	Name string
	Version string
	MaxPackageSize int
	MaxConn int
}

var GlobalConfig *Config

var configPath="config/config.json"

func (c *Config)loadFromJson(){
	log.Printf("Read config from %s",configPath)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Config Load error"+err.Error())
		return
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Println("Config Load error"+err.Error())
		return
	}
}

func init(){
	GlobalConfig=&Config{
		Host:           "0.0.0.0",
		Port:           8080,
		Name:           "CZinx",
		Version:        "v0.4",
		MaxPackageSize: 512,
		MaxConn:        1024,
	}
	GlobalConfig.loadFromJson()
}
