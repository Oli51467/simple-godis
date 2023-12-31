package main

import (
	"fmt"
	"os"
	_ "simple-godis/command"
	"simple-godis/config"
	"simple-godis/lib/logger"
	"simple-godis/resp/handler"
	"simple-godis/server"
)

const ConfigFile string = "redis.conf"

// 没有配置文件时的默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6378,
}

// 判断文件是否存在
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "simple-godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
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
	}, handler.MakeRespHandler())
	if err != nil {
		logger.Error(err)
	}
}
