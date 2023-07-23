package database

import (
	"simple-godis/interface/resp"
	"simple-godis/lib/logger"
	"simple-godis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connection, cmd [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(cmd)
}

func (e *EchoDatabase) Close() {
	logger.Info("Database closed")
}

func (e *EchoDatabase) AfterClientClose(conn resp.Connection) {
	logger.Info("After client close, free some memory")
}
