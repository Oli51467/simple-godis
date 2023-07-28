package clus

import (
	"context"
	"errors"
	"simple-godis/interface/resp"
	"simple-godis/lib/logger"
	"simple-godis/lib/utils"
	"simple-godis/resp/client"
	"simple-godis/resp/reply"
	"strconv"
)

// getPeerConnection 从连接池中拿到对应peer节点的连接 peer:兄弟节点的地址
func (cluster *ClusterDatabase) getPeerConnection(peer string) (*client.ClusterClient, error) {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return nil, errors.New("connection not found when borrowing")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	peerClient, ok := object.(*client.ClusterClient)
	if !ok {
		return nil, errors.New("get peer client failed")
	}
	return peerClient, nil
}

// returnPeerConnection 使用完连接通信完成后归还连接
func (cluster *ClusterDatabase) returnPeerConnection(peer string, peerClient *client.ClusterClient) error {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection not found when returning")
	}
	return pool.ReturnObject(context.Background(), peerClient)
}

// relay 将指令转发到集群的另一个节点 转发规则由哈希计算
func (cluster *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply {
	// 如果目标节点是自己 直接由本机数据库执行
	if peer == cluster.self {
		return cluster.db.Exec(conn, args)
	}
	peerClient, err := cluster.getPeerConnection(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		err := cluster.returnPeerConnection(peer, peerClient)
		if err != nil {
			logger.Error("return peer client failed")
		}
	}()
	// 由于兄弟节点不知道client的存在，也就不知道用户选择了几号数据库，所有用户选择数据库的记录都在本地，所以在转发指令前先选好数据库
	peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(conn.GetDBIndex())))
	return peerClient.Send(args) // 最后将要执行的指令发到集群节点上
}

// broadcast 向集群内的所有节点广播转发一条指令
func (cluster *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		// 遍历每个节点执行转发
		result := cluster.relay(node, conn, args)
		results[node] = result
	}
	return results
}
