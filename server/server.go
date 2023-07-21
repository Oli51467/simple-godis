package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"simple-godis/interface/tcp"
	"simple-godis/lib/logger"
	"sync"
	"syscall"
)

/**
 * A tcp server
 */

// Config 定义了tcp服务器的属性
type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	signalChan := make(chan os.Signal) // 负载是系统的信号
	// 当系统调用发出如下信号时传递给signalChan
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 监听系统调用 如果被外部关闭 则给closeChan发信号
	go func() {
		signals := <-signalChan
		switch signals {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

// ListenAndServe binds port and handle requests, blocking until close
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	go func() { // 如果直接关闭程序 defer仍然可能走不到 通过closeConn通知该方法进行关闭
		<-closeChan // closeChan的负载时空结构体，监听closeChan的信号时如果没有数据则会在这里阻塞
		logger.Info("shutting down...")
		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()  // close connections
	}()

	defer func() { // 遇到错误退出时关闭连接
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background() // 获取上下文 这里的上下文是空的 可以设置过期时间等
	var waitDone sync.WaitGroup // 新建一个WaitGroup 等待客户端退出
	for {
		conn, err := listener.Accept()
		if err != nil { // 接收连接错误时等待正在服务的客户端退出
			break
		}
		logger.Info("Accept Connection")
		waitDone.Add(1)
		go func() { // 一个协程监听一个连接 处理这个新的连接需要新建一个协程
			defer func() { // 防止Handle时出现panic时执行不到waitDone.Done
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
