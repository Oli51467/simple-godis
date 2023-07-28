package smap

import "sync"

// SyncMap 实现Map接口
type SyncMap struct {
	m sync.Map
}

// MakeSyncMap SyncMap的构造方法
func MakeSyncMap() *SyncMap {
	return &SyncMap{}
}

func (s *SyncMap) Get(key string) (val interface{}, exists bool) {
	val, ok := s.m.Load(key)
	return val, ok
}

func (s *SyncMap) Len() int {
	length := 0
	s.m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (s *SyncMap) Put(key string, val interface{}) (result int) {
	_, existed := s.m.Load(key)
	s.m.Store(key, val)
	if existed {
		return 0
	}
	return 1
}

func (s *SyncMap) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := s.m.Load(key)
	if existed {
		return 0
	}
	s.m.Store(key, val)
	return 1
}

func (s *SyncMap) PutIfExists(key string, val interface{}) (result int) {
	_, existed := s.m.Load(key)
	if !existed {
		return 0
	}
	s.m.Store(key, val)
	return 1
}

func (s *SyncMap) Remove(key string) (result int) {
	_, existed := s.m.Load(key)
	if !existed {
		return 0
	}
	s.m.Delete(key)
	return 1
}

func (s *SyncMap) ForEach(consumer Consumer) {
	s.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncMap) Keys() []string {
	result := make([]string, s.Len())
	i := 0
	s.m.Range(func(key, value interface{}) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

func (s *SyncMap) RandomKeys(limit int) []string {
	result := make([]string, s.Len())
	for i := 0; i < limit; i++ {
		s.m.Range(func(key, value interface{}) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (s *SyncMap) RandomDistinctKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	s.m.Range(func(key, value interface{}) bool {
		result[i] = key.(string)
		i++
		return i < limit
	})
	return result
}

func (s *SyncMap) Clear() {
	*s = *MakeSyncMap()
}
