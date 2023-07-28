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

// Members 返回集合中的所有元素的切片
func (set *Set) Members() []string {
	slice := make([]string, set.Len())
	i := 0
	set.s.ForEach(func(key string, val interface{}) bool {
		if i < len(slice) {
			slice[i] = key
		} else {
			// 如果遍历过程中有元素被追加进去，则选择追加的方式
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

// Intersect 将集合取交集
func Intersect(sets ...*Set) *Set {
	result := MakeSet()
	if len(sets) == 0 {
		return result
	}
	counterMap := make(map[string]int)
	for _, set := range sets {
		set.ForEach(func(key string) bool {
			counterMap[key]++
			return true
		})
	}
	for k, v := range counterMap {
		// 如果一个元素在所有集合中都出现过，即出现的次数等于集合的个数
		if v == len(sets) {
			result.Add(k)
		}
	}
	return result
}

// Union 将集合取并集
func Union(sets ...*Set) *Set {
	result := MakeSet()
	for _, set := range sets {
		set.ForEach(func(key string) bool {
			result.Add(key)
			return true
		})
	}
	return result
}

// Diff 将集合取差集
func Diff(sets ...*Set) *Set {
	if len(sets) == 0 {
		return MakeSet()
	}
	// 将第一个集合的元素拷贝到另一个result集合
	result := sets[0].ShallowCopy()
	for i := 1; i < len(sets); i++ {
		sets[i].ForEach(func(key string) bool {
			// 将另一个集合中存在的键移除
			result.Remove(key)
			return true
		})
		if result.Len() == 0 {
			break
		}
	}
	return result
}

// ShallowCopy 将一个集合中的所有元素拷贝到另一个集合
func (set Set) ShallowCopy() *Set {
	result := MakeSet()
	set.ForEach(func(key string) bool {
		result.Add(key)
		return true
	})
	return result
}

// RandomMembers 随机返回给定数量的键，可能包含重复的键
func (set *Set) RandomMembers(limit int) []string {
	if set == nil || set.s == nil {
		return nil
	}
	return set.s.RandomKeys(limit)
}

// RandomDistinctMembers 随机返回给定数量的键，不会包含重复的键
func (set *Set) RandomDistinctMembers(limit int) []string {
	return set.s.RandomDistinctKeys(limit)
}
