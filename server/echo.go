package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"simple-godis/lib/logger"
	"simple-godis/lib/sync/atomic"
	"simple-godis/lib/sync/wait"
	"sync"
	"time"
)

/**
 * A echo server to response client
 */

// EchoHandler 向客户端回复消息 维护一个连接
type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

// EchoClient 客户端数据结构
type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

// Close 关闭客户端连接
func (client *EchoClient) Close() error {
	client.Waiting.WaitWithTimeout(10 * time.Second) // 等待Client的WaitGroup的计数器清0 会在收到message时Add(1)
	_ = client.Conn.Close()
	return nil
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// 判断状态 如果引擎正在关闭 则直接将连接关闭
	if handler.closing.Get() {
		_ = conn.Close()
	}

	// 封装成一个EchoClient 将连接Conn作为初始化参数
	client := &EchoClient{
		Conn: conn,
	}
	// activeConn存储所有客户端的连接
	handler.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	// 不断读取Client传递过来的信心
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection close by client probably because of EOF")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)     // 将msg转成byte
		_, _ = conn.Write(b) // 将数据写回
		client.Waiting.Done()
	}
}

// Close stops echo handler
func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	// 将业务引擎的状态设置为关闭中
	handler.closing.Set(true)
	// 将所有客户端的连接关闭
	handler.activeConn.Range(func(key, value interface{}) bool { // 返回的bool表示要不要遍历下一个key
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
