package list

import "container/list"

const pageSize = 128

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

// Add 新建一个Node添加到列表尾部
func (ql *QuickList) Add(val interface{}) {
	ql.size++
	// 如果链表长度为空 需要新建一页 将数据插入到新的一页中 再将这一页插入到页表中
	if ql.data.Len() == 0 {
		page := make([]interface{}, 0, pageSize)
		page = append(page, val)
		ql.data.PushBack(page)
		return
	}
	// 取出最后一页
	bacKNode := ql.data.Back()
	backPage := bacKNode.Value.([]interface{}) // backNode的Value就是这其所在的zipList
	// 如果最后一页已经满了 需要建新的一页
	if len(backPage) == cap(backPage) {
		// 创建一个zipList
		page := make([]interface{}, 0, pageSize)
		page = append(page, val)
		// 新创建一个quickListNode节点 把这个新创建的节点插入到quickList双向链表中。
		ql.data.PushBack(page)
		return
	}
	backPage = append(backPage, val)
	bacKNode.Value = backPage
}

// find 根据一个指定的index，返回该index所在的zipList(page)以及页内偏移
func (ql *QuickList) find(index int) *iterator {
	if ql == nil {
		panic("quickList is nil")
	}
	if index < 0 || index >= ql.size {
		panic("find in quickList index out of bound")
	}
	var node *list.Element // quickList结点 存储着一个zipList
	var zipListInNode []interface{}
	var totalElementCount int
	// 根据index的大小判断从头部查找还是尾部查找
	if index < ql.size/2 {
		node = ql.data.Front()
		totalElementCount = 0
		for {
			// 取到这个quickListNode所对应的zipList页
			zipListInNode = node.Value.([]interface{})
			// 如果之前累计的Element总数加上该页的元素总数大于index，则所找的节点在该页中
			if totalElementCount+len(zipListInNode) > index {
				break
			}
			totalElementCount += len(zipListInNode)
			node = node.Next()
		}
	} else {
		node = ql.data.Back()
		totalElementCount = ql.size
		for {
			// 取到这个quickListNode所对应的zipList页
			zipListInNode = node.Value.([]interface{})
			totalElementCount -= ql.size
			// 同理，在该页中
			if totalElementCount <= index {
				break
			}
			node = node.Prev()
		}
	}
	// 用所要找的索引减去在它之前页的总元素数 就是所要找的结点在页内的偏移
	pageOffset := index - totalElementCount
	return &iterator{
		node:   node,
		offset: pageOffset,
		ql:     ql,
	}
}

// get 根据迭代器中记录的offset和所在页拿到指定的Element
func (iter *iterator) get() interface{} {
	return iter.page()[iter.offset]
}

// page 返回QuickList迭代器所在的QuickNode结点的zipList[]interface{}
func (iter *iterator) page() []interface{} {
	return iter.node.Value.([]interface{})
}

// Get 根据下标返回QuickList中对应的元素
func (ql *QuickList) Get(index int) (val interface{}) {
	iter := ql.find(index)
	return iter.get()
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
