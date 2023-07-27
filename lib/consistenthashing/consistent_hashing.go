package consistenthashing

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

// NodeMap stores nodes and you can pick node from NodeMap
type NodeMap struct {
	hashFunc  HashFunc       // hash函数
	hashNodes []int          // 每个节点的哈希
	hashMap   map[int]string // 将每个节点的hash映射到节点
}

// NewNodeMap 构造nodes节点map
func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc: fn,
		hashMap:  make(map[int]string),
	}
	if fn == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// IsEmpty returns if there is no node in NodeMap
func (m NodeMap) IsEmpty() bool {
	return len(m.hashNodes) == 0
}

// AddNode add the given nodes into consistent hash circle
func (m *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.hashNodes = append(m.hashNodes, hash)
		m.hashMap[hash] = key
	}
	sort.Ints(m.hashNodes)
}

// PickNode gets the closest item in the hash to the provided key.
func (m *NodeMap) PickNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.hashNodes), func(i int) bool {
		return m.hashNodes[i] >= hash
	})
	// Means we have cycled back to the first replica.
	if idx == len(m.hashNodes) {
		idx = 0
	}
	return m.hashMap[m.hashNodes[idx]]
}
