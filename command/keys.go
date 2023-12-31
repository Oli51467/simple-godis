package command

import (
	"simple-godis/database"
	List "simple-godis/datastructure/list"
	"simple-godis/datastructure/set"
	"simple-godis/datastructure/smap"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/lib/wildcard"
	"simple-godis/resp/reply"
)

/*
有关key的相关操作 都作为executeCommand的执行方法
*/

func init() {
	database.RegisterCommand("del", executeDel, -2)
	database.RegisterCommand("exists", executeExists, -2)
	database.RegisterCommand("flush", executeFlush, -1)
	database.RegisterCommand("type", executeType, 2)
	database.RegisterCommand("rename", executeRename, 3)
	database.RegisterCommand("renameNx", executeRenameNx, 3)
	database.RegisterCommand("keys", executeKeys, 2)
}

// executeDel 执行删除keys方法
func executeDel(db *database.DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.RemoveEntities(keys...)
	if deleted > 0 {
		db.AddAof(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// executeDel 给定一个或多个key，判读在指定数据库中key是否存在
func executeExists(db *database.DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// executeFlush 在指定数据库中删除所有key
func executeFlush(db *database.DB, args [][]byte) resp.Reply {
	db.FlushKeys()
	db.AddAof(utils.ToCmdLine2("flush", args...))
	return reply.MakeOkReply()
}

// executeType 给定一个key 返回key对应value的类型
func executeType(db *database.DB, args [][]byte) resp.Reply {
	// 第一个参数就是key
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeStatusReply("None")
	}
	// TODO:实现其他数据类型
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	case List.List:
		return reply.MakeStatusReply("list")
	case set.Set:
		return reply.MakeStatusReply("set")
	case smap.Map:
		return reply.MakeStatusReply("hash")
	}
	return reply.MakeUnknownErrReply()
}

// executeRename 键的重命名 rename key1 key2 执行会覆盖key2
func executeRename(db *database.DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])
	destKey := string(args[1])
	entity, exists := db.GetEntity(srcKey)
	if !exists {
		return reply.MakeErrReply(srcKey + "not exists")
	}
	db.PutEntity(destKey, entity)
	db.RemoveEntity(srcKey)
	db.AddAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeOkReply()
}

// executeRenameNx 键的重命名 rename key1 key2 执行会检查key2是否存在
func executeRenameNx(db *database.DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])
	destKey := string(args[1])
	_, ok := db.GetEntity(destKey)
	if !ok {
		return reply.MakeIntReply(0)
	}
	entity, exists := db.GetEntity(srcKey)
	if !exists {
		return reply.MakeErrReply(srcKey + "not exists")
	}
	db.PutEntity(destKey, entity)
	db.RemoveEntity(srcKey)
	db.AddAof(utils.ToCmdLine2("renameNx", args...))
	return reply.MakeIntReply(1)
}

// executeKeys 返回所有的key
func executeKeys(db *database.DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.Data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}
