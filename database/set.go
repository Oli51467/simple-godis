package database

import (
	HashSet "simple-godis/datastructure/set"
	dbInterface "simple-godis/interface/database"
	"simple-godis/resp/reply"
)

// GetAsSet 以key为键获取一个set
func (db *DB) GetAsSet(key string) (*HashSet.Set, reply.ErrorReply) {
	entity, existed := db.GetEntity(key)
	if !existed {
		return nil, nil
	}
	hashSet, ok := entity.Data.(*HashSet.Set)
	if !ok {
		return nil, reply.MakeErrReply("Wrong Type")
	}
	return hashSet, nil
}

// GetOrInitSet 根据一个key从数据库尝试获取一个集合 如果没有就创建一个新的
func (db *DB) GetOrInitSet(key string) (set *HashSet.Set, init bool, errorReply reply.ErrorReply) {
	set, errorReply = db.GetAsSet(key)
	if errorReply != nil {
		return nil, false, errorReply
	}
	init = false
	if set == nil {
		set = HashSet.MakeSet()
		db.PutEntity(key, &dbInterface.DataEntity{
			Data: set,
		})
		init = true
	}
	return set, init, nil
}
