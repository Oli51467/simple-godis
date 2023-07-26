package reply

/*
记录一些固定的回复格式或内容
*/
var pongBytes = []byte("PONG\r\n") // pong的字节数组
var okBytes = []byte("OK\r\n")
var nullBulkBytes = []byte("nil\r\n")     // nil 空字符串回复
var emptyMultiBulkBytes = []byte("0\r\n") // 空数组回复
var noBytes = []byte("")                  // 空回复

/*
本地持有一些固定回复 节约内存
*/
var theOkReply = new(OkReply)
var thePongReply = new(PongReply)
var theNullBulkReply = new(NullBulkReply)
var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)
var theNoReply = new(NoReply)

// PongReply 回复客户端的Ping
type PongReply struct {
}

// MakePongReply PongReply制作方法
func MakePongReply() *PongReply {
	return thePongReply
}

// ToBytes PongReply实现Reply接口的ToBytes方法
func (reply *PongReply) ToBytes() []byte {
	return pongBytes
}

func (reply *PongReply) ToClient() []byte {
	return pongBytes
}

// OkReply 回复客户端OK
type OkReply struct {
}

// MakeOkReply OkReply制作方法
func MakeOkReply() *OkReply {
	return theOkReply
}

// ToBytes OkReply实现Reply接口的ToBytes方法
func (reply *OkReply) ToBytes() []byte {
	return okBytes
}

func (reply *OkReply) ToClient() []byte {
	return okBytes
}

// NullBulkReply 空回复
type NullBulkReply struct {
}

// MakeNullBulkReply NullBulkReply制作方法
func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkReply
}

// ToBytes NullBulkReply实现Reply接口的ToBytes方法
func (reply *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func (reply *NullBulkReply) ToClient() []byte {
	return nullBulkBytes
}

// EmptyMultiBulkReply 空数组回复
type EmptyMultiBulkReply struct {
}

// MakeEmptyMultiBulkReply EmptyMultiBulkReply制作方法
func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

// ToBytes EmptyMultiBulkReply实现Reply接口的ToBytes方法
func (reply *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func (reply *EmptyMultiBulkReply) ToClient() []byte {
	return emptyMultiBulkBytes
}

// NoReply 空回复
type NoReply struct {
}

// MakeNoReply NoReply制作方法
func MakeNoReply() *NoReply {
	return theNoReply
}

// ToBytes EmptyMultiBulkReply实现Reply接口的ToBytes方法
func (reply *NoReply) ToBytes() []byte {
	return noBytes
}
