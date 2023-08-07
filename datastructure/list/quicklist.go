package list

import "container/list"

const pageSize = 128

// QuickList zipList和linkedList的结合
type QuickList struct {
	data *list.List
	size int
}

// locator QuickList的定位器
type locator struct {
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
func (ql *QuickList) find(index int) *locator {
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
	return &locator{
		node:   node,
		offset: pageOffset,
		ql:     ql,
	}
}

// get 根据定位器中记录的offset和所在页拿到指定的Element
func (locator *locator) get() interface{} {
	return locator.page()[locator.offset]
}

// page 返回QuickList定位器所在的QuickNode结点的zipList[]interface{}
func (locator *locator) page() []interface{} {
	return locator.node.Value.([]interface{})
}

// Get 根据下标返回QuickList中对应的元素
func (ql *QuickList) Get(index int) (val interface{}) {
	locator := ql.find(index)
	return locator.get()
}

// set 根据定位器的记录 将QuickList指定位置的键设置为指定的值
func (locator *locator) set(val interface{}) {
	page := locator.page()
	page[locator.offset] = val
}

// Set 将指定位置的值更新 指定位置应该在[0, list.size]之间
func (ql *QuickList) Set(index int, val interface{}) {
	locator := ql.find(index)
	locator.set(val)
}

// Insert 将元素插入链表的指定位置
func (ql *QuickList) Insert(index int, val interface{}) {
	// 如果插入位置在链表尾部 则等同于在尾部新加一个元素
	if index == ql.size {
		ql.Add(val)
		return
	}
	// 找到要插入的那一页
	locator := ql.find(index)
	page := locator.node.Value.([]interface{})
	// 如果这一页还有容量可以添加
	if len(page) < pageSize {
		// 将这一页中要插入位置后面的元素往后移动一位 从offset+1位置起，依次追加offset以及后面的元素
		page = append(page[:locator.offset+1], page[locator.offset:]...)
		page[locator.offset] = val
		// 更新该元素所在的page
		locator.node.Value = page
		ql.size++
		return
	}
	// 这一页已经满了 插入整页可能会引起内存复制，所以将整页拆分成两个半页
	var nextPage []interface{}
	nextPage = append(nextPage, page[pageSize/2:]...)
	page = page[:pageSize/2]
	// 要插入的位置在前半段
	if locator.offset < len(page) {
		page = append(page[:locator.offset+1], page[locator.offset:]...)
		page[locator.offset] = val
	} else {
		// 要插入的位置在后半段
		afterOffset := locator.offset - pageSize/2
		nextPage = append(nextPage[:afterOffset+1], nextPage[afterOffset:]...)
		nextPage[afterOffset] = val
	}
	// 更新该元素所在的page
	locator.node.Value = page
	ql.data.InsertAfter(nextPage, locator.node)
	ql.size++
}

// remove 移除定位器中的元素
func (locator *locator) remove() interface{} {
	page := locator.page()
	val := page[locator.offset]
	// 将要删除的元素移除 从第offset+1个位置开始往前移动一个位置
	page = append(page[:locator.offset], page[locator.offset+1:]...)
	if len(page) > 0 {
		locator.node.Value = page
		if locator.offset == len(page) {
			if locator.node != locator.ql.data.Back() {
				locator.node = locator.node.Next()
				locator.offset = 0
			}
		}
	} else {
		if locator.node == locator.ql.data.Back() {
			locator.ql.data.Remove(locator.node)
			locator.node = nil
			locator.offset = 0
		} else {
			nextNode := locator.node.Next()
			locator.ql.data.Remove(locator.node)
			locator.node = nextNode
			locator.offset = 0
		}
	}
	locator.ql.size--
	return val
}

// Remove 移除指定位置的元素
func (ql *QuickList) Remove(index int) (val interface{}) {
	locator := ql.find(index)
	return locator.remove()
}

// RemoveLast 将列表中的最后一个元素移除 如果最后一页只有一个元素 则将最后一页也删除
func (ql *QuickList) RemoveLast() (val interface{}) {
	if ql.Len() == 0 {
		return nil
	}
	ql.size--
	lastNode := ql.data.Back()
	lastPage := lastNode.Value.([]interface{})
	if len(lastPage) == 1 {
		ql.data.Remove(lastNode) // 将最后一整页删除
		return lastPage[0]       // 返回最后一页的第一个元素
	}
	removeVal := lastPage[len(lastPage)-1]
	lastPage = lastPage[:len(lastPage)-1]
	lastNode.Value = lastPage
	return removeVal
}

// atEnd 判断定位器是否已经在列表的尾部
func (locator *locator) atEnd() bool {
	if locator.ql.data.Len() == 0 {
		return true
	}
	if locator.node != locator.ql.data.Back() {
		return false
	}
	page := locator.page()
	return locator.offset == len(page)
}

// next 将定位器移动到下一个位置
func (locator *locator) next() bool {
	page := locator.page()
	if locator.offset < len(page)-1 {
		locator.offset++
		return true
	}
	// 已经移动到了最后一个元素
	if locator.node == locator.ql.data.Back() {
		locator.offset = len(page)
		return false
	}
	// 移动到下一页
	locator.offset = 0
	locator.node = locator.node.Next()
	return true
}

// RemoveAllByVal 移除列表中所有值为val的元素
func (ql *QuickList) RemoveAllByVal(expected Expected) int {
	locator := ql.find(0) // 从第一个元素开始找
	removedCount := 0
	for !locator.atEnd() {
		// 如果定位器的所在的值等于给定元素的值，则删除定位器所在位置的元素
		if expected(locator.get()) {
			locator.remove()
			removedCount++
		} else {
			locator.next() // 否则定位器移动到下一个位置
		}
	}
	return removedCount
}

// RemoveByVal 从左到右扫描列表，移除该列表中的给定的值，并最多移除count个
func (ql *QuickList) RemoveByVal(expected Expected, count int) int {
	if ql.size == 0 {
		return 0
	}
	locator := ql.find(0)
	removedCount := 0
	for !locator.atEnd() {
		if expected(locator.get()) {
			locator.remove()
			removedCount++
			if removedCount == count {
				break
			}
		} else {
			locator.next()
		}
	}
	return removedCount
}

// Len 返回列表的长度
func (ql *QuickList) Len() int {
	return ql.size
}

// ForEach 遍历列表的每一个元素，如果consumer返回false则终止遍历，否则一直遍历到列表尾部
func (ql *QuickList) ForEach(consumer Consumer) {
	if ql == nil {
		panic("list is nil")
	}
	if ql.Len() == 0 {
		return
	}
	locator := ql.find(0)
	i := 0
	for {
		goOn := consumer(i, locator.get())
		if !goOn {
			break
		}
		i++
		if !locator.next() {
			break
		}
	}
}

// Contains 检查列表中是否包含指定元素
func (ql *QuickList) Contains(expected Expected) bool {
	contains := false
	ql.ForEach(func(i int, valInList interface{}) bool {
		if expected(valInList) {
			contains = true
			return false
		}
		return true
	})
	return contains
}

// Range 返回列表中下标从[start, end)的所有元素
func (ql *QuickList) Range(start int, end int) []interface{} {
	if start < 0 || start > ql.Len() {
		panic("`start` index out of range")
	}
	if end < start || end > ql.Len() {
		panic("`end` index out of range")
	}
	sliceSize := end - start
	rangeElements := make([]interface{}, 0, sliceSize)
	locator := ql.find(start)
	i := 0
	if i < sliceSize {
		rangeElements = append(rangeElements, locator.get())
		locator.next()
		i++
	}
	return rangeElements
}
