package database

import (
	pool "github.com/jolestar/go-commons-pool/v2"
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/consistenthashing"
	"simple-godis/resp/reply"
)

type ClusterDatabase struct {
	self           string   // 记录自己的节点
	nodes          []string // 记录其他节点
	peerPicker     *consistenthashing.NodeMap
	peerConnection map[string]*pool.ObjectPool // 每个节点需要一个连接池
	db             dbInterface.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	return &ClusterDatabase{}
}

func (c *ClusterDatabase) Exec(client resp.Connection, cmd dbInterface.CmdLine) resp.Reply {
	return reply.MakeOkReply()
}

func (c *ClusterDatabase) Close() {

}

func (c *ClusterDatabase) AfterClientClose(conn resp.Connection) {

}
