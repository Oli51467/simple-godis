package set

import "simple-godis/datastructure/smap"

type Consumer func(key string) bool

// Set 线程安全的集合
type Set struct {
	s smap.Map
}

// MakeSet 构造方法
func MakeSet(members ...string) *Set {
	set := &Set{
		s: smap.MakeSyncMap(),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

// Add 向集合中添加一个元素
func (set *Set) Add(val string) int {
	return set.s.Put(val, nil)
}

// Remove 从集合中移除一个元素
func (set *Set) Remove(val string) int {
	ret := set.s.Remove(val)
	return ret
}

// Has 判断集合中有无该元素
func (set *Set) Has(val string) bool {
	if set == nil || set.s == nil {
		return false
	}
	_, exists := set.s.Get(val)
	return exists
}

// Len 集合的长度
func (set *Set) Len() int {
	if set == nil || set.s == nil {
		return 0
	}
	return set.s.Len()
}

func (set *Set) Members() []string {
	slice := make([]string, set.Len())
	i := 0
	set.s.ForEach(func(key string, val interface{}) bool {
		if i < len(slice) {
			slice[i] = key
		} else {
			// 如果遍历过程中有元素被追加进去
			slice = append(slice, key)
		}
		i++
		return true
	})
	return slice
}

// ForEach 遍历集合的每个元素
func (set *Set) ForEach(consumer Consumer) {
	if set == nil || set.s == nil {
		return
	}
	set.s.ForEach(func(key string, val interface{}) bool {
		return consumer(key)
	})
}
