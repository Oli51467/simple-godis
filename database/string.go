package database

import "simple-godis/resp/reply"

func (db *DB) GetAsString(key string) ([]byte, reply.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return nil, reply.MakeWrongTypeReply()
	}
	return bytes, nil
}
