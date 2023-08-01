package command

import (
	"simple-godis/database"
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
	"strconv"
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
	database.RegisterCommand("append", executeAppend, 3)
	database.RegisterCommand("getDel", executeGetAndDel, 2)
	database.RegisterCommand("incr", executeIncr, -2)
	database.RegisterCommand("decr", executeDecr, -2)
}

// executeGet 执行获取一个键对应的value
func executeGet(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	set, _ := db.GetAsSet(key)
	// 判断要取的key是不是一个集合 如果是就改为查询集合内的所有元素
	if set != nil {
		return executeSMembers(db, args)
	}
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
	entity, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	db.PutEntity(key, &dbInterface.DataEntity{
		Data: val,
	})
	db.AddAof(utils.ToCmdLine2("getset", args...))
	return reply.MakeBulkReply(entity)
}

// executeStrLen 获取一个key对应val的长度
func executeStrLen(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	bytes, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(len(bytes)))
}

// executeAppend 向某个value为字符串的key追加元素
func executeAppend(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	bytes, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	bytes = append(bytes, args[1]...)
	db.PutEntity(key, &dbInterface.DataEntity{
		Data: bytes,
	})
	db.AddAof(utils.ToCmdLine3("append", args...))
	return reply.MakeIntReply(int64(len(bytes)))
}

// executeGetAndDel 通过key拿到一个value然后再删除
func executeGetAndDel(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	if val == nil {
		return reply.MakeNullBulkReply()
	}
	db.RemoveEntity(key)
	db.AddAof(utils.ToCmdLine3("del", args...))
	return reply.MakeBulkReply(val)
}

// executeIncr 如果value是int类型 则将value的值增加指定的值 如果不指定 则增加1
func executeIncr(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	bytes, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	if bytes != nil {
		val, err := strconv.ParseInt(string(bytes), 10, 64)
		if err != nil {
			return reply.MakeErrReply("ERR value is not an integer or out of range")
		}
		var increment int64
		// 如果制定了增量的值 则尝试将增量的值转为Int 如果转化失败 则抛出异常
		if len(args) > 1 && args[1] != nil {
			increment, err = strconv.ParseInt(string(args[1]), 10, 64)
			if err != nil {
				return reply.MakeErrReply("ERR the increment argument is NaN")
			}
		} else {
			increment = 1
		}
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: []byte(strconv.FormatInt(val+increment, 10)),
		})
		db.AddAof(utils.ToCmdLine3("incr", args...))
		return reply.MakeIntReply(val + increment)
	} else {
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: []byte("1"),
		})
		db.AddAof(utils.ToCmdLine3("incr", args...))
		return reply.MakeIntReply(1)
	}
}

// executeDecr 如果value是int类型 则将value的值减少指定的值 如果不指定 则减少1
func executeDecr(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	bytes, err := db.GetAsString(key)
	if err != nil {
		return err
	}
	if bytes != nil {
		val, err := strconv.ParseInt(string(bytes), 10, 64)
		if err != nil {
			return reply.MakeErrReply("ERR value is not an integer or out of range")
		}
		var decrement int64
		// 如果制定了减量的值 则尝试将减量的值转为Int 如果转化失败 则抛出异常
		if len(args) > 1 && args[1] != nil {
			decrement, err = strconv.ParseInt(string(args[1]), 10, 64)
			if err != nil {
				return reply.MakeErrReply("ERR the increment argument is NaN")
			}
		} else {
			decrement = 1
		}
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: []byte(strconv.FormatInt(val-decrement, 10)),
		})
		// 落盘
		db.AddAof(utils.ToCmdLine3("decr", args...))
		return reply.MakeIntReply(val - decrement)
	} else {
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: []byte("-1"),
		})
		db.AddAof(utils.ToCmdLine3("decr", args...))
		return reply.MakeIntReply(-1)
	}
}
