package parser

import (
	"bufio"
	"errors"
	"io"
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
	"strconv"
	"strings"
)

/*
ParserStream 将用户根据resp通信协议发来的字节流解析成redis指令
*/

// Payload Client发来的数据解析后的负载
type Payload struct {
	Data resp.Reply
	Err  error
}

// readState ParserStream的解析状态
type readState struct {
	readMultiLine     bool     // 是否是多行数据
	expectedArgsCount int      // 期望的参数个数
	msgType           byte     // 指令类型
	args              [][]byte // 每一个参数对应一个字节数组
	bulkLen           int64    // 接下来的一个块要读取的字节数
}

// finished 判断解析是否结束
func finished(state *readState) bool {
	return state.expectedArgsCount > 0 && state.expectedArgsCount == len(state.args)
}

// ParseStream 对外提供的解析流函数 让tcp服务器将io流交给这个函数
// return *Payload 该函数进行解析并通过channel异步告知上层解析结果
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch) // 新建一个协程 异步将输出向上层传递
	return ch
}

func parse0(reader io.Reader, ch chan *Payload) {

}

// readLine在io.Reader流中读取一行
// return: []byte 这一行数据的byte数组 bool 有没有io错误 error 错误本身
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	// 如果读到了$+数字，严格读取字符串个数，否则按/r/n切分
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		// 判断/n的前一个字符是不是/r 或者/n前面没有/r
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("Protocol error " + string(msg))
		}
	} else {
		msg = make([]byte, state.bulkLen+2)
		_, err := io.ReadFull(bufReader, msg) // 将bufReader的内容强行全部读到msg中
		if err != nil {
			return nil, true, err
		}
		// 判断读到的特定长度的字节是不是/r/n结尾 或者根本没有/r/n
		var msgLen = len(msg)
		if msgLen == 0 || msg[msgLen-2] != '\r' || msg[msgLen-1] != '\n' {
			return nil, false, errors.New("Protocol error " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// parseMultiBulkHeader 将readLine切出来的*数据进行解析 然后修改readState状态
// eg: *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$3\r\nVAL\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 3)
	if err != nil {
		return errors.New("Protocol error " + string(msg))
	}
	if expectedLine == 0 { // 没有参数
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		// 有参数，告诉readState还要接着读
		state.expectedArgsCount = int(expectedLine)
		state.msgType = msg[0]                       // 设置指令类型
		state.readMultiLine = true                   // 是读数组 -> 多行
		state.args = make([][]byte, 0, expectedLine) // 初始化参数空间
		return nil
	} else {
		return errors.New("Protocol error " + string(msg))
	}
}

// parseBulkHeader $数据进行解析 然后修改readState状态
// eg: $4\r\nPING\r\n
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("Protocol error " + string(msg))
	}
	// 如果bulkLen是-1，代表用户发的是空字符串
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	}
	return errors.New("Protocol error " + string(msg))
}

// parseSingleLineReply 解析：+OK\r\n -err\r\n :5\r\n
func parseSingleLineReply(msg []byte) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n")
	var result resp.Reply
	switch str[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("Protocol error " + string(msg))
		}
		result = reply.MakeIntReply(val)
	}
	return result, nil
}
