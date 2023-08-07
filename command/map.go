package command

import (
	"simple-godis/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
)

func init() {
	database.RegisterCommand("HSet", executeHSet, 4)
	database.RegisterCommand("HSetNx", executeHSetNx, 4)
	database.RegisterCommand("HGet", executeHGet, 3)
	database.RegisterCommand("HDel", executeHDel, -3)
	database.RegisterCommand("HExists", executeHExists, 3)
	database.RegisterCommand("HMSet", executeHMSet, -3)
	database.RegisterCommand("HMGet", executeHMGet, -3)
	database.RegisterCommand("HKeys", executeHKeys, 2)
	database.RegisterCommand("HValues", executeHValues, 2)
	database.RegisterCommand("HGetAll", executeHGetAll, 2)
	database.RegisterCommand("HLen", executeHLen, 2)
	database.RegisterCommand("HStrlen", executeHStrlen, 3)
}

// executeHSet 将以key1为键的实体中添加映射(field,val)
func executeHSet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	value := args[2]

	iMap, _, errorReply := db.GetOrInitMap(key)
	if errorReply != nil {
		return errorReply
	}
	result := iMap.Put(field, value)
	db.AddAof(utils.ToCmdLine3("HSet", args...))
	return reply.MakeIntReply(int64(result))
}

// executeHSetNx 当且仅当field的键不存在时，将以key1为键的实体中添加映射(field,val)
func executeHSetNx(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])
	value := string(args[2])

	iMap, _, errorReply := db.GetOrInitMap(key)
	if errorReply != nil {
		return errorReply
	}
	// 该键不存在时才放入
	result := iMap.PutIfAbsent(field, value)
	if result > 0 {
		db.AddAof(utils.ToCmdLine3("HSetNx", args...))
	}
	return reply.MakeIntReply(int64(result))
}

// executeHGet 在以key1为键的实体中获取以field为键的value
func executeHGet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])

	// 获取实体
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeErrReply("entity not existed")
	}

	rawVal, exists := iMap.Get(field)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	value, _ := rawVal.([]byte)
	return reply.MakeBulkReply(value)
}

// executeHExists 在以key1为键的实体中检查存不存在以field的键
func executeHExists(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])

	// 获取实体
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeIntReply(0)
	}

	_, exists := iMap.Get(field)
	if exists {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

// executeHDel 在以key1为键的实体中删除一个或多个键值对
func executeHDel(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	fields := make([]string, len(args)-1)
	fieldsArgs := args[1:]
	for i, field := range fieldsArgs {
		fields[i] = string(field)
	}

	// 获取实体
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeIntReply(0)
	}

	deletedCount := 0
	for _, field := range fields {
		result := iMap.Remove(field)
		deletedCount += result
	}
	if iMap.Len() == 0 {
		db.RemoveEntity(key)
	}
	if deletedCount > 0 {
		db.AddAof(utils.ToCmdLine3("HDel", args...))
	}
	return reply.MakeIntReply(int64(deletedCount))
}

// executeHLen 获取以key为键的实体中键值对的数量
func executeHLen(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(iMap.Len()))
}

// executeHStrlen 获取以key为键的实体中指定field为key的value的字符串的长度
func executeHStrlen(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	field := string(args[1])

	// 获取实体
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeIntReply(0)
	}

	// 获取value
	rawVal, exists := iMap.Get(field)
	if exists {
		value, _ := rawVal.([]byte)
		return reply.MakeIntReply(int64(len(value)))
	}
	return reply.MakeIntReply(0)
}

// executeHMSet 获取以key为键的实体中设置多个键值对
func executeHMSet(db *database.DB, args [][]byte) resp.Reply {
	// 参数必须是key filed1 val1 field2 val2 ... 奇数个
	if len(args)%2 == 0 {
		return reply.MakeSyntaxErrReply()
	}
	key := string(args[0])
	kvPairCount := (len(args) - 1) / 2
	fields := make([]string, kvPairCount)
	values := make([][]byte, kvPairCount)

	for i := 0; i < kvPairCount; i++ {
		fields[i] = string(args[i*2+1])
		values[i] = args[i*2+2]
	}

	// 获取实体
	iMap, _, errorReply := db.GetOrInitMap(key)
	if errorReply != nil {
		return errorReply
	}
	// 将多个键值对放入
	for i, field := range fields {
		iMap.Put(field, values[i])
	}
	db.AddAof(utils.ToCmdLine3("HMSet", args...))
	return reply.MakeOkReply()
}

// executeHMGet 在以key为键的实体中获取多个以field为键的值
func executeHMGet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	fieldsCount := len(args) - 1
	fields := make([]string, fieldsCount)
	// 解析参数获取所有field放到fields中
	for i := 0; i < fieldsCount; i++ {
		fields[i] = string(args[i+1])
	}
	// 获取实体
	result := make([][]byte, fieldsCount)
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeMultiBulkReply(result)
	}
	// 遍历每个field
	for i, field := range fields {
		rawValue, exists := iMap.Get(field)
		if !exists {
			result[i] = nil
		} else {
			value, _ := rawValue.([]byte)
			result[i] = value
		}
	}
	return reply.MakeMultiBulkReply(result)
}

// executeHKeys 获取一个哈希表内所有的键
func executeHKeys(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeEmptyMultiBulkReply()
	}
	fields := make([][]byte, iMap.Len())
	cnt := 0
	iMap.ForEach(func(key string, val interface{}) bool {
		fields[cnt] = []byte(key)
		cnt++
		return true
	})
	return reply.MakeMultiBulkReply(fields[:cnt])
}

// executeHValues 获取一个哈希表内所有的值
func executeHValues(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeEmptyMultiBulkReply()
	}
	values := make([][]byte, iMap.Len())
	cnt := 0
	iMap.ForEach(func(key string, val interface{}) bool {
		values[cnt], _ = val.([]byte)
		cnt++
		return true
	})
	return reply.MakeMultiBulkReply(values[:cnt])
}

// executeHGetAll 获取一个哈希表内的所有键值对
func executeHGetAll(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	// 获取实体
	iMap, errorReply := db.GetAsMap(key)
	if errorReply != nil {
		return errorReply
	}
	if iMap == nil {
		return reply.MakeEmptyMultiBulkReply()
	}

	mapSize := iMap.Len()
	result := make([][]byte, mapSize*2)
	cnt := 0
	iMap.ForEach(func(key string, val interface{}) bool {
		result[cnt] = []byte(key)
		cnt++
		result[cnt], _ = val.([]byte)
		cnt++
		return true
	})
	return reply.MakeMultiBulkReply(result)
}
