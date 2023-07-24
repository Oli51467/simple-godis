package database

import (
	"simple-godis/datastructure/smap"
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
	"strings"
)

type CmdLine = [][]byte

// Database 一个子数据库 实现了smap.Map接口
type Database struct {
	index int
	Data  smap.Map
}

// ExecuteCommand 所有redis指令都要使用该函数执行
type ExecuteCommand func(db *Database, args [][]byte) resp.Reply

// MakeDatabase 构建一个数据库
func MakeDatabase() *Database {
	db := &Database{
		Data: smap.MakeSyncMap(),
	}
	return db
}

// Execute 执行指令 conn 一个对应的连接 cmdLine 具体的指令[][]byte
func (db *Database) Execute(conn resp.Connection, cmdLine CmdLine) resp.Reply {
	// 统一将指令转为小写
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := commandTable[cmdName]
	if !ok { // 不存在该指令集
		return reply.MakeErrReply("ERR unknown command " + cmdName)
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	executor := cmd.executor
	return executor(db, cmdLine[1:]) // 将参数切出来
}

// validateArity 验证参数的个数
// 如果参数个数定长 则arity=n，如果参数不定长，则arity=-n，n为参数个数的最小值
func validateArity(arity int, commandArgs [][]byte) bool {
	argNum := len(commandArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

// GetEntity 从该索引的数据库中拿一个key对应的DataEntity return: DataEntity, 是否拿到
func (db *Database) GetEntity(key string) (*dbInterface.DataEntity, bool) {
	raw, ok := db.Data.Get(key) // Get返回的是空接口 需要转换为DataEntity
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*dbInterface.DataEntity)
	return entity, true
}

// PutEntity 从该索引的数据库中放入一个key对应的DataEntity
func (db *Database) PutEntity(key string, entity *dbInterface.DataEntity) int {
	return db.Data.Put(key, entity)
}

// PutEntityIfExists 如果key在该索引对应的数据库中存在，
// 从该索引的数据库中放入一个key对应的DataEntity
func (db *Database) PutEntityIfExists(key string, entity *dbInterface.DataEntity) int {
	return db.Data.PutIfExists(key, entity)
}

// PutEntityIfAbsent 如果key在该索引对应的数据库中不存在，
// 从该索引的数据库中放入一个key对应的DataEntity
func (db *Database) PutEntityIfAbsent(key string, entity *dbInterface.DataEntity) int {
	return db.Data.PutIfAbsent(key, entity)
}

// RemoveEntity 从该索引的数据库中删除一个key对应的DataEntity
func (db *Database) RemoveEntity(key string) {
	db.Data.Remove(key)
}

// RemoveEntities 从该索引的数据库中删除一个或多个key对应的DataEntity
func (db *Database) RemoveEntities(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		deleted += db.Data.Remove(key)
	}
	return deleted
}

// FlushKeys 从该索引的数据库中删除所有key
func (db *Database) FlushKeys() {
	db.Data.Clear()
}
