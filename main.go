package main

import (
	"fmt"
	"os"
	"simple-godis/config"
	"simple-godis/lib/logger"
	"simple-godis/server"
)

const ConfigFile string = "redis.conf"

// 没有配置文件时的默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

// 判断文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path: "logs",
		Name: "simple-godis",
		Ext:  "log",
	})

	if fileExists(ConfigFile) {
		config.SetupConfig(ConfigFile)
	} else {
		config.Properties = defaultProperties
	}

	// 启动服务器并监听
	err := server.ListenAndServeWithSignal(&server.Config{
		Address: fmt.Sprintf("%s:%d",
			config.Properties.Bind,
			config.Properties.Port),
	}, server.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}
