package resp

// Reply 将回复给客户端的内容转换成字节 是客户端-服务端双向通信接口
type Reply interface {
	ToBytes() []byte
	ToClient() []byte
}
