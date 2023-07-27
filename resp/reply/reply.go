package reply

import (
	"bytes"
	"simple-godis/interface/resp"
	"strconv"
)

/*
承载一般的回复体
*/
var (
	nullBulkReplyBytes = []byte("nil")
	CRLF               = "\r\n" // CRLF resp协议的换行符
)

// ErrorReply Redis通信异常回复接口
type ErrorReply interface {
	Error() string
	ToBytes() []byte
	ToClient() []byte
}

// BulkReply 普通字符串通信 将不符合resp通信协议的字符串转化为符合通信协议的字节数组
type BulkReply struct {
	Msg []byte
}

// MakeBulkReply 对外提供的BulkReply构造方法
func MakeBulkReply(msg []byte) *BulkReply {
	return &BulkReply{
		Msg: msg,
	}
}

// ToBytes 实现resp.reply.Reply.ToBytes接口 将BulkReply要发送的字节信息转化为resp协议的信息
func (reply *BulkReply) ToBytes() []byte {
	// 判断传递的信息长度是否为空
	if len(reply.Msg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(reply.Msg)) + CRLF + string(reply.Msg) + CRLF)
}

func (reply *BulkReply) ToClient() []byte {
	// 判断传递的信息长度是否为空
	if len(reply.Msg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte(string(reply.Msg) + CRLF)
}

// MultiBulkReply 多段字符串通信 将不符合resp通信协议的字符串数组转化为符合通信协议的字节数组
type MultiBulkReply struct {
	Msg [][]byte
}

// ToBytes 实现resp.reply.Reply.ToBytes接口 将多行文本转换成resp协议的通信文本
func (reply *MultiBulkReply) ToBytes() []byte {
	msgLen := len(reply.Msg)
	var buf bytes.Buffer // 字符串拼接转bytes
	buf.WriteString("*" + strconv.Itoa(msgLen) + CRLF)
	// 遍历每一个string
	for _, lineMsg := range reply.Msg {
		if lineMsg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(lineMsg)) + CRLF + string(lineMsg) + CRLF)
		}
	}
	return buf.Bytes()
}

func (reply *MultiBulkReply) ToClient() []byte {
	msgLen := len(reply.Msg)
	var buf bytes.Buffer // 字符串拼接转bytes
	buf.WriteString(strconv.Itoa(msgLen) + CRLF)
	// 遍历每一个string
	for _, lineMsg := range reply.Msg {
		if lineMsg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString(string(lineMsg) + CRLF)
		}
	}
	return buf.Bytes()
}

// MakeMultiBulkReply 对外提供的MultiBulkReply构造方法
func MakeMultiBulkReply(msg [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Msg: msg,
	}
}

// StatusReply resp通信协议的状态回复体
type StatusReply struct {
	Status string
}

func (reply *StatusReply) ToClient() []byte {
	return []byte(reply.Status + CRLF)
}

// MakeStatusReply 对外提供的StatusReply构造方法
func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// ToBytes 实现reply.Reply接口
func (reply *StatusReply) ToBytes() []byte {
	return []byte("+" + reply.Status + CRLF)
}

// IntReply resp通信协议的状态回复体，可承载int64
type IntReply struct {
	Code int64
}

// MakeIntReply 对外提供的IntReply的构造方法
func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

// ToBytes 实现reply.Reply接口
func (reply *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(reply.Code, 10) + CRLF)
}

func (reply *IntReply) ToClient() []byte {
	return []byte(strconv.FormatInt(reply.Code, 10) + CRLF)
}

// StandardErrReply 自定义的resp协议的异常状态回复
type StandardErrReply struct {
	Status string
}

// ToBytes 实现reply.Reply接口
func (reply *StandardErrReply) ToBytes() []byte {
	return []byte("-" + reply.Status + CRLF)
}

func (reply *StandardErrReply) ToClient() []byte {
	return []byte(reply.Status + CRLF)
}

// Error 实现reply.Reply接口
func (reply *StandardErrReply) Error() string {
	return reply.Status
}

// MakeErrReply 对外提供的StandardErrReply构造方法
func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

// IsErrorReply 判断回复是异常还是正常 判断第一个字节是不是'-'
func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
