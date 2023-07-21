package resp

// Connection 代表了客户端的一个连接
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
