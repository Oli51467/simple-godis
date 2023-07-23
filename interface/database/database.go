package database

import "simple-godis/interface/resp"

type CmdLine = [][]byte

// Database 数据库接口
// Method：
// Exec执行数据库指令 client：指定客户端连接 cmd: 字节数组指令
// Close关闭数据库连接
// AfterClientClose 连接关闭后的处理
type Database interface {
	Exec(client resp.Connection, cmd CmdLine) resp.Reply
	Close()
	AfterClientClose(conn resp.Connection)
}

// DataEntity 抽象了Redis中所有的数据结构
type DataEntity struct {
	Data interface{}
}
