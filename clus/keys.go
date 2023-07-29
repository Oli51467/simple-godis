package clus

import (
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
)

// ClusterFlushDB removes all data in current database
func ClusterFlushDB(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(conn, cmdArgs)
	var errReply reply.ErrorReply
	for _, v := range replies {
		if reply.IsErrorReply(v) {
			errReply = v.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("error occurs: " + errReply.Error())
}

// ClusterDel 以原子方式从集群中删除给定的键，键可以分布在任何节点上
// 如果给定的键分布在不同的节点上，Del将使用try-commit-catch删除它们
func ClusterDel(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(conn, cmdArgs)
	var errReply reply.ErrorReply
	var deletedCount int64 = 0
	for _, rep := range replies {
		if reply.IsErrorReply(rep) {
			errReply = rep.(reply.ErrorReply)
			break
		}
		intReply, ok := rep.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeErrReply("error")
		}
		deletedCount += intReply.Code
	}
	if errReply == nil {
		return reply.MakeOkReply()
	}
	return reply.MakeErrReply("error occurs: " + errReply.Error())
}
