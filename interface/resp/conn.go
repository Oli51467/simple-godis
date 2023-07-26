package resp

// Connection 接口代表了客户端的一个连接
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
	Close() error
}
