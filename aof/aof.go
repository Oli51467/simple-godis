package aof

import (
	"io"
	"os"
	"simple-godis/config"
	dbInterface "simple-godis/interface/database"
	"simple-godis/lib/logger"
	"simple-godis/lib/utils"
	"simple-godis/resp/client"
	"simple-godis/resp/parser"
	"simple-godis/resp/reply"
	"strconv"
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
	database      dbInterface.Database // 持有数据库
	aofChan       chan *payload        // 数据缓冲区 缓存的是指令的集合
	aofFile       *os.File
	aofFilename   string
	currenDbIndex int // 该文件对应哪个分数据库
}

// NewAofHandler AofHandler的构造方法
func NewAofHandler(database dbInterface.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	handler.database = database
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
	handler.currenDbIndex = 0
	for payload := range handler.aofChan {
		// 需要切换分数据库
		if payload.dbIndex != handler.currenDbIndex {
			handler.currenDbIndex = payload.dbIndex
			// 转成字节数组写入文件中
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(payload.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
		}
		// 如果不需要切换数据库 或者已经切换好数据库 直接将指令的字节数组写入文件
		data := reply.MakeMultiBulkReply(payload.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
}

// loadAof 在服务启动时将磁盘中的resp格式的指令当作用户发来的指令恢复
func (handler *AofHandler) loadAof() {
	file, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			logger.Error("Close AofFile Failed", err)
		}
	}()
	// 使用解析器解析Aof文件的历史指令 并将解析结果吐到ch管道里 再遍历管道还原指令
	ch := parser.ParseStream(file)
	dummyClient := &client.Client{}
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF {
				break
			}
			logger.Error(payload.Err)
			continue
		}
		if payload.Data == nil {
			logger.Error("Empty payload")
			continue
		}
		res, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("Aof order need multi bulkReplyFormat", payload.Data)
			continue
		}
		// 将取到的指令送到数据库执行
		executeReply := handler.database.Exec(dummyClient, res.Msg)
		if reply.IsErrorReply(executeReply) {
			logger.Error(executeReply)
		}
	}
}
