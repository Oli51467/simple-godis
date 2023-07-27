package command

import (
	"simple-godis/database"
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
)

/*
有关key的相关操作 都作为executeCommand的执行方法
*/

func init() {
	database.RegisterCommand("get", executeGet, -2)
	database.RegisterCommand("set", executeSet, 3)
	database.RegisterCommand("setnx", executeSetnx, 3)
	database.RegisterCommand("getset", executeGetAndSet, 3)
	database.RegisterCommand("strlen", executeStrLen, 2)
}

// executeGet 执行获取一个键对应的value
func executeGet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	respMes := entity.Data.([]byte)
	return reply.MakeBulkReply(respMes)
}

// executeGet 执行获取一个键对应的value
func executeSet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &dbInterface.DataEntity{
		Data: val,
	}
	db.PutEntity(key, entity)
	db.AddAof(utils.ToCmdLine2("set", args...))
	return reply.MakeOkReply()
}

// executeSetnx 如果对应的key不存在，执行获取一个键对应的value
func executeSetnx(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := &dbInterface.DataEntity{
		Data: val,
	}
	result := db.PutEntityIfAbsent(key, entity)
	db.AddAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(result))
}

// executeGetAndSet 获取key对应的value，并设置为新的值，返回旧的值
func executeGetAndSet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity, exist := db.GetEntity(key)
	db.PutEntity(key, &dbInterface.DataEntity{
		Data: val,
	})
	if !exist {
		return reply.MakeNullBulkReply()
	}
	db.AddAof(utils.ToCmdLine2("getset", args...))
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// executeStrLen 获取一个key对应val的长度
func executeStrLen(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}
