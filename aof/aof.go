package aof

import (
	"os"
	"simple-godis/config"
	"simple-godis/database"
)

const aofBufferSize = 1 << 8

type CmdLine [][]byte

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// AofHandler 落盘处理器
// 构造方法
// 将用户指令包装成payload放到缓冲区aofChan中去，再将aofChan中的数据落到硬盘中
// 加载 将磁盘中的aof指令加载出来
type AofHandler struct {
	databases       database.Databases // 持有数据库
	aofChan         chan *payload      // 数据缓冲区 缓存的是指令的集合
	aofFile         *os.File
	aofFilename     string
	currentDatabase int // 该文件对应哪个分数据库
}

// NewAofHandler AofHandler的构造方法
func NewAofHandler(databases database.Databases) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.databases = databases
	// LoadAof程序启动时将磁盘中的aof文件加载出来
	handler.loadAof()
	aofFile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	// 创建缓冲区
	handler.aofChan = make(chan *payload, aofBufferSize)
	// 新建协程用于接收
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

// AddAof 将一条指令语句追加到缓冲区aofChan中
func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}
}

// HandleAof 将缓冲区aofChan中的内容源源不断地往外取，并保存到磁盘中
func (handler *AofHandler) handleAof() {

}

func (handler *AofHandler) loadAof() {

}
