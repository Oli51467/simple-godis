package handler

import (
	"context"
	"io"
	"net"
	"simple-godis/database"
	dbInterface "simple-godis/interface/database"
	"simple-godis/lib/logger"
	"simple-godis/lib/sync/atomic"
	"simple-godis/resp/client"
	"simple-godis/resp/parser"
	"simple-godis/resp/reply"
	"strings"
	"sync"
)

// RespHandler TCP层处理resp协议
type RespHandler struct {
	activeConn sync.Map
	db         dbInterface.Database
	closing    atomic.Boolean
}

func MakeRespHandler() *RespHandler {
	var db dbInterface.Database
	db = database.MakeDatabase()
	return &RespHandler{
		db: db,
	}
}

// Handle 实现handler.Handle方法 处理客户端的连接
func (handler *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	newClient := client.NewClient(conn)             // 使用conn新建一个客户端连接
	handler.activeConn.Store(newClient, struct{}{}) // 将连接存储
	ch := parser.ParseStream(conn)                  // 解析器不断监听管道的数据并将处理后的数据传递到channel中
	for payload := range ch {                       // ch在被关闭之前会一直循环输出
		if payload.Err != nil { // 如果监听的指令存在错误
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				handler.closeClient(newClient) // 客户端关闭连接
				return
			} else { // 协议解析错误
				errReply := reply.MakeErrReply(payload.Err.Error())
				err := newClient.Write(errReply.ToBytes()) // 将错误写回客户端
				if err != nil {                            // 回复用户出错时出错
					handler.closeClient(newClient)
					return
				}
				continue
			}
		} else { // 指令不存在错误
			if payload.Data == nil { // 用户发送的指令为空
				logger.Info("empty command received")
				continue
			}
			// 2.指令不为空 尝试将数据转换成二维字节数组
			command, ok := payload.Data.(*reply.MultiBulkReply)
			if !ok {
				continue
			}
			// 3.转换成功，db执行指令
			execResult := handler.db.Exec(newClient, command.Msg)
			if execResult != nil {
				err := newClient.Write(execResult.ToBytes())
				if err != nil { // 回复用户出错时出错
					handler.closeClient(newClient)
					return
				}
			} else {
				err := newClient.Write(reply.MakeUnknownErrReply().ToBytes())
				if err != nil { // 回复用户出错时出错
					handler.closeClient(newClient)
					return
				}
			}
		}
	}
}

// closeClient 关闭指定的一个客户端的连接
func (handler *RespHandler) closeClient(client *client.Client) {
	logger.Info("Connection closed: " + client.RemoteAddr().String())
	_ = client.Close()
	handler.db.AfterClientClose(client) // 关闭后的处理
	handler.activeConn.Delete(client)   // 连接池移除连接
}

// Close 实现handler.Close方法
func (handler *RespHandler) Close() error {
	logger.Info("Handler shutting down...")
	handler.closing.Set(true)

	handler.activeConn.Range(
		func(key interface{}, value interface{}) bool {
			session := key.(*client.Client)
			_ = session.Close()
			return true
		},
	)
	// 关闭数据库连接
	handler.db.Close()
	return nil
}
