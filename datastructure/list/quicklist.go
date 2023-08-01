package list

import "container/list"

// QuickList zipList和linkedList的结合
type QuickList struct {
	data *list.List
	size int
}

// iterator QuickList的迭代器
type iterator struct {
	node   *list.Element
	offset int
	ql     *QuickList
}

// MakeQuickList 初始化一个QuickList
func MakeQuickList() *QuickList {
	return &QuickList{
		data: list.New(),
	}
}

func (ql *QuickList) Add(val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Get(index int) (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Set(index int, val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Insert(index int, val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Remove(index int) (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) RemoveLast() (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) RemoveAllByVal(expected Expected) int {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) RemoveByVal(expected Expected, count int) int {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Len() int {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) ForEach(consumer Consumer) {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Contains(expected Expected) bool {
	//TODO implement me
	panic("implement me")
}

func (ql *QuickList) Range(start int, end int) []interface{} {
	//TODO implement me
	panic("implement me")
}
