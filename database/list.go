package database

import (
	List "simple-godis/datastructure/list"
	"simple-godis/interface/database"
	"simple-godis/resp/reply"
)

// GetAsList 以key为键获取一个list
func (db *DB) GetAsList(key string) (List.List, reply.ErrorReply) {
	entity, existed := db.GetEntity(key)
	if !existed {
		return nil, nil
	}
	list, ok := entity.Data.(List.List)
	if !ok {
		return nil, reply.MakeWrongTypeReply()
	}
	return list, nil
}

// GetOrInitList 尝试以key为键从数据库中获取一个List实体，如果获取不到则创建一个新的并放入到数据库实体中
func (db *DB) GetOrInitList(key string) (list List.List, init bool, errorReply reply.ErrorReply) {
	list, errorReply = db.GetAsList(key)
	if errorReply != nil {
		return nil, false, errorReply
	}
	init = false
	if list == nil {
		// 初始化一个新的List，并将list放到数据库实例中
		list = List.MakeQuickList()
		db.PutEntity(key, &database.DataEntity{
			Data: list,
		})
		init = true
	}
	return list, init, nil

}
