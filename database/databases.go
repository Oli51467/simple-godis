package database

import (
	"simple-godis/config"
	"simple-godis/interface/resp"
	"simple-godis/lib/logger"
	"simple-godis/resp/reply"
	"strconv"
	"strings"
)

// NewDatabases 初始化数据库和分库
func NewDatabases() *Databases {
	databases := &Databases{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	databases.dbSet = make([]*Database, config.Properties.Databases)
	for i := range databases.dbSet {
		database := MakeDatabase()
		database.index = i
		databases.dbSet[i] = database
	}
	return databases
}

type Databases struct {
	dbSet []*Database
}

func (db *Databases) Exec(client resp.Connection, args CmdLine) resp.Reply {
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

func (db *Databases) Close() {
	logger.Info("Database closed")
}

func (db *Databases) AfterClientClose(conn resp.Connection) {
	logger.Info("After client close, free some memory")
}

// executeSelect 执行选择数据库指令
func executeSelect(conn resp.Connection, databases *Databases, args [][]byte) resp.Reply {
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
