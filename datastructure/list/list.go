package list

// Consumer List迭代器
type Consumer func(i int, val interface{}) bool

// Expected 抽象函数 检查给定的值是否等于期望的值
type Expected func(a interface{}) bool

type List interface {
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	Len() int
	ForEach(consumer Consumer)
	Contains(expected Expected) bool
	Range(start int, end int) []interface{}
}
