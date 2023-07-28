package clus

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"simple-godis/clus/factory"
	"simple-godis/config"
	"simple-godis/database"
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/lib/consistenthashing"
	"simple-godis/resp/reply"
)

// ClusterDatabase 集群模式数据库 数据有三种执行模式 单节点返回、转发、群发
type ClusterDatabase struct {
	self           string   // 记录自己的节点
	nodes          []string // 记录集群中所有的节点
	peerPicker     *consistenthashing.NodeMap
	peerConnection map[string]*pool.ObjectPool // 每个节点需要一个连接池
	db             dbInterface.Database
}

// MakeClusterDatabase 新建了集群之间的连接和连接池的连接，新建了一致性哈希集合和所有节点的列表
func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,             // 配置文件中本机的地址
		db:             database.MakeStandaloneDatabases(), // 本机的数据库
		peerPicker:     consistenthashing.NewNodeMap(nil),  // 存储集群各个节点信息的map
		peerConnection: make(map[string]*pool.ObjectPool),  // 各个节点之间的连接池
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1) // 新建nodes列表，存储所有的节点
	// 遍历所有配置的其他节点，将其他节点和自己都加入到nodes中
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, cluster.self)
	cluster.peerPicker.AddNode(nodes...)
	ctx := context.Background()
	// 初始化该节点与其他各个节点之间的连接池
	for _, peer := range config.Properties.Peers {
		pool.NewObjectPoolWithDefaultConfig(ctx, &factory.ConnectionFactory{Peer: peer})
	}
	cluster.nodes = nodes
	return cluster
}

func (cluster *ClusterDatabase) Exec(client resp.Connection, cmd dbInterface.CmdLine) resp.Reply {
	return reply.MakeOkReply()
}

func (cluster *ClusterDatabase) Close() {

}

func (cluster *ClusterDatabase) AfterClientClose(conn resp.Connection) {

}
