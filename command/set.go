package command

import (
	"simple-godis/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
)

func init() {
	database.RegisterCommand("sAdd", executeSAdd, -3)
	database.RegisterCommand("sIsMember", executeSIsMember, 3)
	database.RegisterCommand("sRem", executeSRemove, -3)
	database.RegisterCommand("sMembers", executeSMembers, 2)
}

// executeGet 执行获取一个键对应的value
func executeSAdd(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	members := args[1:]

	// 从数据库获取这个key对应的set，如果获取不到，就初始化一个
	set, _, errorReply := db.GetOrInitSet(key)
	if errorReply != nil {
		return errorReply
	}
	counter := 0
	for _, member := range members {
		counter += set.Add(string(member))
	}
	db.AddAof(utils.ToCmdLine3("sAdd", args...))
	return reply.MakeIntReply(int64(counter))
}

// executeGet 执行一个value是否在一个key的集合中
func executeSIsMember(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	member := string(args[1])

	set, errorReply := db.GetAsSet(key)
	if errorReply != nil {
		return errorReply
	}
	if set == nil {
		return reply.MakeIntReply(0)
	}
	has := set.Has(member)
	if has {
		return reply.MakeIntReply(1)
	}
	return reply.MakeIntReply(0)
}

// executeSRemove 执行将一个或多个元素从集合中删除，如果删除后集合中没有元素，则删除该key
func executeSRemove(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])
	members := args[1:]

	set, errorReply := db.GetAsSet(key)
	if errorReply != nil {
		return errorReply
	}
	if set == nil {
		return reply.MakeIntReply(0)
	}
	counter := 0
	for _, member := range members {
		counter += set.Remove(string(member))
	}
	if set.Len() == 0 {
		db.RemoveEntity(key)
	}
	if counter > 0 {
		db.AddAof(utils.ToCmdLine3("sRem", args...))
	}
	return reply.MakeIntReply(int64(counter))
}

// executeSMembers 列出以key为键的集合中的所有元素
func executeSMembers(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	// 尝试获取一个集合
	set, errorReply := db.GetAsSet(key)
	if errorReply != nil {
		return errorReply
	}
	if set == nil {
		return reply.MakeIntReply(0)
	}
	members := make([][]byte, set.Len())
	i := 0
	set.ForEach(func(member string) bool {
		members[i] = []byte(member)
		i++
		return true
	})
	return reply.MakeMultiBulkReply(members)
}
