package parser

import (
	"bufio"
	"errors"
	"io"
	"runtime/debug"
	"simple-godis/interface/resp"
	"simple-godis/lib/logger"
	"simple-godis/resp/reply"
	"strconv"
	"strings"
)

/*
parser 将用户根据resp通信协议发来的字节流解析成redis指令
*/

// Payload Client发来的数据解析后的负载，向上返回给上层handler而不是client
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

// ParseStream 对外提供的异步解析流函数 让tcp服务器将io流交给这个函数
// return *Payload 该函数进行解析并通过channel异步告知上层解析结果
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch) // 新建一个协程，每一个用户一个解析器，异步将输出向上层传递
	return ch
}

func parse0(reader io.Reader, ch chan *Payload) {
	// 声明recover，作用是在不断接收的过程中如果panic不让循环退出
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		var ioErr bool
		// 读一行数据
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			// 是io错误
			if ioErr {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		} // 处理错误 如果是io错误则关闭连接，如果是协议错误则继续读取
		// 判断是不是多行解析模式 如果不是 也有可能是未初始化
		if !state.readMultiLine {
			// 1.如果开头是* 是一个数组
			if msg[0] == '*' {
				err := parseMultiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("Protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.MakeEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' { // 2.如果开头是$ 是一个多行字符串
				err := parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: errors.New("Protocol error: " + string(msg)),
					}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: reply.MakeNullBulkReply(),
					}
					state = readState{}
					continue
				}
			} else { // 3.如果开头是+、-、: 解析单行
				result, err := parseSingleLineReply(msg)
				// 直接抛给上层handler
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else { // 进入多行模式 $3/r/n
			err := readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: errors.New("Protocol error: " + string(msg)),
				}
				state = readState{}
				continue
			}
			if finished(&state) {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// finished 判断解析是否结束
func finished(state *readState) bool {
	return state.expectedArgsCount > 0 && state.expectedArgsCount == len(state.args)
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

// parseMultiBulkHeader 前缀是*数据 解析然后修改readState状态
// eg: *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$3\r\nVAL\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine uint64
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 3)
	if err != nil {
		return errors.New("Protocol error " + string(msg))
	}
	if expectedLine == 0 { // 没有参数 这里需要parse0判断expectedArgsCount进行处理
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
	// 如果bulkLen是-1，代表用户发的是空字符串 这里需要parse0判断bulkLen进行处理
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
	default:
		// 当作文本协议解析
		strs := strings.Split(str, " ")
		args := make([][]byte, len(strs))
		for i, s := range strs {
			strings.TrimSpace(s)
			if len(s) == 0 {
				return nil, errors.New("Protocol error " + string(msg))
			}
			args[i] = []byte(s)
		}
		result = reply.MakeMultiBulkReply(args)
	}
	return result, nil
}

// readBody 解析
// eg: $3\r\nSET\r\n$3\r\nKEY\r\n$3\r\nVAL\r\n
// eg: PING\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2] // 先将后面的/r/n切掉
	var err error
	// case: $3 把后面带的长度取出来 赋值到bulkLen
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("Protocol error " + string(msg))
		}
		// $0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
