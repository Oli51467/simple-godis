package database

import (
	"simple-godis/aof"
	"simple-godis/config"
	"simple-godis/interface/resp"
	"simple-godis/lib/logger"
	"simple-godis/resp/reply"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

// MakeStandaloneDatabases 初始化数据库和分库以及处理指令文件记录的处理器
func MakeStandaloneDatabases() *StandaloneDatabase {
	databases := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	databases.dbSet = make([]*DB, config.Properties.Databases)
	for i := range databases.dbSet {
		database := MakeDB()
		database.index = i
		databases.dbSet[i] = database
	}
	// 初始化AofHandler
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(databases)
		if err != nil {
			panic(err)
		}
		databases.aofHandler = aofHandler
		// 将落盘方法逐个添加到每个分数据库中
		for _, db := range databases.dbSet {
			finalDb := db
			finalDb.AddAof = func(line CmdLine) {
				databases.aofHandler.AddAof(finalDb.index, line)
			}
		}
	}
	return databases
}

func (db *StandaloneDatabase) Exec(client resp.Connection, args CmdLine) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return executeSelect(client, db, args[1:])
	}
	dbIndex := client.GetDBIndex()
	database := db.dbSet[dbIndex]
	return database.Execute(client, args)
}

func (db *StandaloneDatabase) Close() {
	logger.Info("DB closed")
}

func (db *StandaloneDatabase) AfterClientClose(conn resp.Connection) {

}

// executeSelect 执行选择数据库指令
func executeSelect(conn resp.Connection, databases *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid database index")
	}
	if dbIndex >= len(databases.dbSet) {
		return reply.MakeErrReply("ERR database index is out of range")
	}
	conn.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
