package resp

// Reply 将回复给客户端的内容转换成字节
type Reply interface {
	ToBytes() []byte
}
