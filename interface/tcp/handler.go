package tcp

import (
	"context"
	"net"
)

// Handler 只处理tcp连接
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
