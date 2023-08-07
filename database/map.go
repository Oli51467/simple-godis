package database

import (
	"simple-godis/datastructure/smap"
	dbInterface "simple-godis/interface/database"
	"simple-godis/resp/reply"
)

// GetAsMap 以key为键获取一个map
func (db *DB) GetAsMap(key string) (smap.Map, reply.ErrorReply) {
	entity, existed := db.GetEntity(key)
	if !existed {
		return nil, nil
	}
	iMap, ok := entity.Data.(smap.Map)
	if !ok {
		return nil, reply.MakeWrongTypeReply()
	}
	return iMap, nil
}

// GetOrInitMap 尝试以key为键从数据库中获取一个Map实体，如果获取不到则创建一个新的并放入到数据库实体中
func (db *DB) GetOrInitMap(key string) (iMap smap.Map, init bool, errorReply reply.ErrorReply) {
	iMap, errorReply = db.GetAsMap(key)
	if errorReply != nil {
		return nil, false, errorReply
	}
	init = false
	if iMap == nil {
		iMap = smap.MakeSimpleMap()
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: iMap,
		})
		init = true
	}
	return iMap, init, nil
}
