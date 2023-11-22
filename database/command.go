package database

import (
	"strings"
)

/*
每一个类型的指令对应一个执行方法command
*/

var CommandTable = make(map[string]*command)

// command 一种类型的指令对应一个command
type command struct {
	executor ExecuteCommand // 具体对应的是哪个执行函数
	arity    int            // 参数数量
}

// RegisterCommand input: name指令名称 executor具体的执行函数 arity参数个数
// 新建一个command放到commandTable中
func RegisterCommand(name string, executor ExecuteCommand, arity int) {
	name = strings.ToLower(name)
	CommandTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}
