package reply

/*
记录一些固定的错误回复内容
*/
var unknownErrBytes = []byte("-Err unknown\r\n")
var syntaxErrBytes = []byte("-Err syntax error\r\n")
var wrongTypeErrBytes = []byte("-Wrong type operation against a key holding the wrong kind of value\r\n")

var theSyntaxErrReply = &SyntaxErrReply{}
var theUnknownErrReply = &UnknownErrReply{}

// UnknownErrReply 一类未知错误的抽象 实现了reply.ErrorReply接口
type UnknownErrReply struct {
}

// MakeUnknownErrReply UnknownErrReply的构造方法
func MakeUnknownErrReply() *UnknownErrReply {
	return theUnknownErrReply
}

// Error UnknownErrReply实现reply.ErrorReply接口的Error方法
func (err *UnknownErrReply) Error() string {
	return "Err unknown"
}

// ToBytes UnknownErrReply实现reply.ErrorReply接口的ToBytes方法
func (err UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

func (err *UnknownErrReply) ToClient() []byte {
	return unknownErrBytes
}

// ArgsNumErrReply Redis指令参数错误的抽象 实现了reply.ErrorReply接口
type ArgsNumErrReply struct {
	Cmd string
}

// Error ArgsNumErrReply实现reply.ErrorReply接口的Error方法
func (err *ArgsNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + err.Cmd + "' command"
}

// ToBytes ArgsNumErrReply实现reply.ErrorReply接口的ToBytes方法
func (err *ArgsNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + err.Cmd + "' command\r\n")
}

func (err *ArgsNumErrReply) ToClient() []byte {
	return []byte("-ERR wrong number of arguments for '" + err.Cmd + "' command\r\n")
}

// MakeArgNumErrReply ArgNumErrReply的构造方法
func MakeArgNumErrReply(cmd string) *ArgsNumErrReply {
	return &ArgsNumErrReply{
		Cmd: cmd,
	}
}

// SyntaxErrReply 语法错误的抽象 实现了reply.ErrorReply接口
type SyntaxErrReply struct{}

// MakeSyntaxErrReply SyntaxErrReply的构造方法
func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

// ToBytes SyntaxErrReply实现reply.ErrorReply接口的ToBytes方法
func (err *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

// Error SyntaxErrReply实现reply.ErrorReply接口的Error方法
func (err *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// WrongTypeErrReply key/value类型错误的抽象 实现了reply.ErrorReply接口
type WrongTypeErrReply struct{}

// ToBytes WrongTypeErrReply实现reply.ErrorReply接口的ToBytes方法
func (err *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

// Error WrongTypeErrReply实现reply.ErrorReply接口的Error方法
func (err *WrongTypeErrReply) Error() string {
	return "Wrong type operation against a key holding the wrong kind of value"
}

// ProtocolErrReply 接口协议错误的抽象 实现了reply.ErrorReply接口
type ProtocolErrReply struct {
	Msg string
}

// ToBytes ProtocolErrReply实现reply.ErrorReply接口的ToBytes方法
func (err *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + err.Msg + "'\r\n")
}

// Error ProtocolErrReply实现reply.ErrorReply接口的Error方法
func (err *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + err.Msg
}
