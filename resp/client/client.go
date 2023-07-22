package client

import (
	"net"
	"simple-godis/lib/sync/wait"
	"sync"
	"time"
)

// Client 抽象了客户端的连接
type Client struct {
	conn       net.Conn
	waiting    wait.Wait
	mutex      sync.Mutex
	selectedDB int
}

// NewClient 指定conn新建一个客户端的连接
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// RemoteAddr 获取连接会话的远程地址
func (session *Client) RemoteAddr() net.Addr {
	return session.conn.RemoteAddr()
}

// Close 实现io.Close
func (session *Client) Close() error {
	// 关闭连接时需要先等待一次读写完成
	session.waiting.WaitWithTimeout(10 * time.Second)
	_ = session.conn.Close()
	return nil
}

// Write 给客户端发送数据
func (session *Client) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	// 同一时刻只能有一个协程向客户端写
	session.mutex.Lock()
	session.waiting.Add(1)
	defer func() {
		session.waiting.Done()
		session.mutex.Unlock()
	}()
	_, err := session.conn.Write(bytes)
	return err
}

// GetDBIndex 返回客户端指定的数据库
func (session *Client) GetDBIndex() int {
	return session.selectedDB
}

// SelectDB 修改客户端连接的数据库
func (session *Client) SelectDB(dbIndex int) {
	session.selectedDB = dbIndex
}
