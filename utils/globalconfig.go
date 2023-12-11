package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type GlobalConf struct {
	Name           string
	Host           string
	TcpPort        int
	MaxConn        int
	MaxPackageSize uint32
}

var GlobalConfig *GlobalConf

func (g GlobalConf) LoadFromConf() {
	data, err := os.ReadFile("config/conf.json")
	if err != nil {
		fmt.Println("Load config err:", err)

		return
	}

	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		fmt.Println("Wrong conf file data. err:", err)

		return
	}
}

func init() {
	GlobalConfig = &GlobalConf{
		Name:           "GoLinServer",
		TcpPort:        2333,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalConfig.LoadFromConf()
}