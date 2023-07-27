package smap

type SimpleMap struct {
	m map[string]interface{}
}

func MakeSimpleMap() *SimpleMap {
	return &SimpleMap{
		m: make(map[string]interface{}),
	}
}

func (s *SimpleMap) Get(key string) (val interface{}, exists bool) {
	val, ok := s.m[key]
	return val, ok
}

func (s *SimpleMap) Len() int {
	if s.m == nil {
		panic("m is nil")
	}
	return len(s.m)
}

func (s *SimpleMap) Put(key string, val interface{}) (result int) {
	_, existed := s.m[key]
	s.m[key] = val
	if existed {
		return 0
	}
	return 1
}

func (s *SimpleMap) PutIfAbsent(key string, val interface{}) (result int) {
	_, existed := s.m[key]
	if existed {
		return 0
	}
	s.m[key] = val
	return 1
}

func (s *SimpleMap) PutIfExists(key string, val interface{}) (result int) {
	_, existed := s.m[key]
	if existed {
		s.m[key] = val
		return 1
	}
	return 0
}

func (s *SimpleMap) Remove(key string) (result int) {
	_, existed := s.m[key]
	delete(s.m, key)
	if existed {
		return 1
	}
	return 0
}

func (s *SimpleMap) ForEach(consumer Consumer) {
	for k, v := range s.m {
		if !consumer(k, v) {
			break
		}
	}
}

func (s *SimpleMap) Keys() []string {
	result := make([]string, len(s.m))
	i := 0
	for k := range s.m {
		result[i] = k
		i++
	}
	return result
}

func (s *SimpleMap) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		for k := range s.m {
			result[i] = k
			break
		}
	}
	return result
}

func (s *SimpleMap) RandomDistinctKeys(limit int) []string {
	size := limit
	if size > len(s.m) {
		size = len(s.m)
	}
	result := make([]string, size)
	i := 0
	for k := range s.m {
		if i == size {
			break
		}
		result[i] = k
		i++
	}
	return result
}

func (s *SimpleMap) Clear() {
	*s = *MakeSimpleMap()
}
