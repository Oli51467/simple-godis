package clus

import (
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
)

// clusterRename 判断改名后key的hash是否改变，如果改变则不支持rename
func clusterRename(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within one slot in cluster mode")
	}
	return cluster.relay(destPeer, conn, cmdArgs)
}
