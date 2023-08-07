package command

import (
	"simple-godis/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
	"strconv"
)

func init() {
	database.RegisterCommand("LPush", executeLPush, -3)
	database.RegisterCommand("LPushX", executeLPushX, -3)
	database.RegisterCommand("RPush", executeRPush, -3)
	database.RegisterCommand("RPushX", executeRPushX, -3)
	database.RegisterCommand("LPop", executeLPop, 2)
	database.RegisterCommand("RPop", executeRPop, 2)
	database.RegisterCommand("LIndex", executeLIndex, 3)
	database.RegisterCommand("LSet", executeLSet, 4)
	database.RegisterCommand("LLen", executeLLen, 2)
}

// executeLIndex 查找下标为index的元素
func executeLIndex(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	indexInt64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("Index value is not an integer or out out range")
	}
	index := int(indexInt64)
	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	listSize := list.Len()
	// 如果查询下标为负，则为回环逆向查找
	if index < -1*listSize || index >= listSize {
		return reply.MakeNullBulkReply()
	} else if index < 0 {
		index = index + listSize
	}
	val, _ := list.Get(index).([]byte)
	return reply.MakeBulkReply(val)
}

// executeLPush 在列表头部插入一个元素
func executeLPush(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	values := args[1:]

	list, _, errorReply := db.GetOrInitList(key)
	if errorReply != nil {
		return errorReply
	}
	for _, value := range values {
		list.Insert(0, value)
	}
	db.AddAof(utils.ToCmdLine3("LPush", args...))
	return reply.MakeIntReply(int64(list.Len()))
}

// executeRPush 在列表尾部插入一个元素
func executeRPush(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	values := args[1:]

	list, _, errorReply := db.GetOrInitList(key)
	if errorReply != nil {
		return errorReply
	}

	for _, value := range values {
		list.Add(value)
	}
	db.AddAof(utils.ToCmdLine3("RPush", args...))
	return reply.MakeIntReply(int64(list.Len()))
}

// executeLPushX 当且仅当列表存在时，在列表头部插入一个元素
func executeLPushX(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	values := args[1:]

	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	// 如果列表不存在，则不插入
	if list == nil {
		return reply.MakeIntReply(0)
	}
	for _, value := range values {
		list.Insert(0, value)
	}
	db.AddAof(utils.ToCmdLine3("LPushX", args...))
	return reply.MakeIntReply(int64(list.Len()))
}

// executeRPushX 当且仅当列表存在时，在列表尾部插入一个元素
func executeRPushX(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	values := args[1:]

	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	// 如果列表不存在，则不插入
	if list == nil {
		return reply.MakeIntReply(0)
	}
	for _, value := range values {
		list.Add(value)
	}
	db.AddAof(utils.ToCmdLine3("RPushX", args...))
	return reply.MakeIntReply(int64(list.Len()))
}

// executeLPop 将列表的第一个元素弹出并返回
func executeLPop(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	removeVal, _ := list.Remove(0).([]byte)
	if list.Len() == 0 {
		db.RemoveEntity(key)
	}
	db.AddAof(utils.ToCmdLine3("LPop", args...))
	return reply.MakeBulkReply(removeVal)
}

// executeRPop 将列表尾部的元素弹出并返回
func executeRPop(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	if list == nil {
		return reply.MakeNullBulkReply()
	}
	removeVal, _ := list.RemoveLast().([]byte)
	if list.Len() == 0 {
		db.RemoveEntity(key)
	}
	db.AddAof(utils.ToCmdLine3("RPop", args...))
	return reply.MakeBulkReply(removeVal)
}

// executeLSet 在列表指定位置放置元素
func executeLSet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	index64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("Index value is not an integer or out out range")
	}
	index := int(index64)
	setVal := args[2]

	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	if list == nil {
		return reply.MakeErrReply("ERR no such key")
	}
	listSize := list.Len()
	// 如果查询下标为负，则为回环逆向查找
	if index < -1*listSize || index >= listSize {
		return reply.MakeErrReply("set index is out of range")
	} else if index < 0 {
		index = index + listSize
	}
	list.Set(index, setVal)
	db.AddAof(utils.ToCmdLine3("LSet", args...))
	return reply.MakeOkReply()
}

// executeLLen 查询列表的长度
func executeLLen(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	list, errorReply := db.GetAsList(key)
	if errorReply != nil {
		return errorReply
	}
	if list == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(list.Len()))
}
