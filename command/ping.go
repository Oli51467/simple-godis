package command

import (
	"simple-godis/database"
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
)

// Ping 指令的ExecuteCommand方法
func Ping(db *database.DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

// init 初始化时执行
func init() {
	database.RegisterCommand("ping", Ping, -1)
}
