package command

import (
	"simple-godis/database"
	HashSet "simple-godis/datastructure/set"
	"simple-godis/interface/resp"
	"simple-godis/lib/utils"
	"simple-godis/resp/reply"
	"strconv"
)

func init() {
	database.RegisterCommand("sAdd", executeSAdd, -3)
	database.RegisterCommand("sIsMember", executeSIsMember, 3)
	database.RegisterCommand("sRem", executeSRemove, -3)
	database.RegisterCommand("sMembers", executeSMembers, 2)
	database.RegisterCommand("sCard", executeSCard, 2)
	database.RegisterCommand("sInter", executeSIntersection, -2)
	database.RegisterCommand("sUnion", executeUnion, -2)
	database.RegisterCommand("sDiff", executeDiff, -2)
	database.RegisterCommand("sPop", executeSPop, -2)
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

// executeSCard 返回集合中的元素个数
func executeSCard(db *database.DB, args [][]byte) resp.Reply {
	key := string(args[0])

	// 尝试获取一个集合
	set, errorReply := db.GetAsSet(key)
	if errorReply != nil {
		return errorReply
	}
	if set == nil {
		return reply.MakeIntReply(0)
	}
	return reply.MakeIntReply(int64(set.Len()))
}

// executeSIntersection 将多个集合取交集
func executeSIntersection(db *database.DB, args [][]byte) resp.Reply {
	sets := make([]*HashSet.Set, 0, len(args))
	for _, arg := range args {
		key := string(arg)
		set, errorReply := db.GetAsSet(key)
		if errorReply != nil {
			return errorReply
		}
		if set.Len() == 0 {
			return reply.MakeEmptyMultiBulkReply()
		}
		sets = append(sets, set)
	}
	result := HashSet.Intersect(sets...)
	return set2reply(result)
}

// executeSIntersection 将多个集合取并集
func executeUnion(db *database.DB, args [][]byte) resp.Reply {
	sets := make([]*HashSet.Set, 0, len(args))
	for _, arg := range args {
		key := string(arg)
		set, errorReply := db.GetAsSet(key)
		if errorReply != nil {
			return errorReply
		}
		sets = append(sets, set)
	}
	result := HashSet.Union(sets...)
	return set2reply(result)
}

// executeDiff 将多个集合取差集
func executeDiff(db *database.DB, args [][]byte) resp.Reply {
	sets := make([]*HashSet.Set, 0, len(args))
	for _, arg := range args {
		key := string(arg)
		set, errReply := db.GetAsSet(key)
		if errReply != nil {
			return errReply
		}
		sets = append(sets, set)
	}
	result := HashSet.Diff(sets...)
	return set2reply(result)
}

// executeSPop 随机从集合中移除指定数量的元素 不指定就是1个
func executeSPop(db *database.DB, args [][]byte) resp.Reply {
	if len(args) != 1 && len(args) != 2 {
		return reply.MakeErrReply("ERR wrong number of arguments for sPop command")
	}
	key := string(args[0])

	// 拿到对应的数据库
	set, errReply := db.GetAsSet(key)
	if errReply != nil {
		return errReply
	}
	if set == nil {
		return reply.MakeNullBulkReply()
	}
	count := 1
	if len(args) == 2 {
		count64, err := strconv.ParseInt(string(args[1]), 10, 64)
		if err != nil || count64 <= 0 {
			return reply.MakeErrReply("ERR value is out of range, must be positive")
		}
		count = int(count64)
	}
	if count > set.Len() {
		count = set.Len()
	}
	members := set.RandomDistinctMembers(count) // 随机取出count个集合内的元素
	result := make([][]byte, len(members))
	for i, val := range members {
		set.Remove(val)
		result[i] = []byte(val)
	}
	if count > 0 {
		db.AddAof(utils.ToCmdLine3("sPop", args...))
	}
	return reply.MakeMultiBulkReply(result)
}

// set2reply 将一个set转成二维字节数组的reply
func set2reply(set *HashSet.Set) resp.Reply {
	arr := make([][]byte, set.Len())
	i := 0
	set.ForEach(func(key string) bool {
		arr[i] = []byte(key)
		i++
		return true
	})
	return reply.MakeMultiBulkReply(arr)
}
