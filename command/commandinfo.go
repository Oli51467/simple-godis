package command

import (
	"simple-godis/database"
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
)

func init() {
	database.RegisterCommand("commands", executeCommands, 1)
}

// executeDel 执行删除keys方法
func executeCommands(db *database.DB, args [][]byte) resp.Reply {
	return getAllGodisCommandReply()
}

func getAllGodisCommandReply() resp.Reply {
	replies := make([]resp.Reply, 0, len(database.CommandTable))
	for k := range database.CommandTable {
		replies = append(replies, reply.MakeStatusReply(k))
	}
	return reply.MakeMultiRawReply(replies)
}
