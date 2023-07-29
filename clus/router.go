package clus

import "simple-godis/interface/resp"

type CmdLine = [][]byte
type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["ping"] = LocalRouter
	routerMap["select"] = LocalRouter

	routerMap["del"] = ClusterDel
	routerMap["flush"] = ClusterFlushDB

	routerMap["rename"] = clusterRename
	routerMap["renamenx"] = clusterRename

	routerMap["exists"] = defaultClusterRouter
	routerMap["type"] = defaultClusterRouter

	routerMap["set"] = defaultClusterRouter
	routerMap["setnx"] = defaultClusterRouter
	routerMap["get"] = defaultClusterRouter
	routerMap["getset"] = defaultClusterRouter

	routerMap["sAdd"] = defaultClusterRouter
	routerMap["sIsMember"] = defaultClusterRouter
	routerMap["sRem"] = defaultClusterRouter
	routerMap["sMembers"] = defaultClusterRouter
	routerMap["sCard"] = defaultClusterRouter
	routerMap["sInter"] = defaultClusterRouter
	routerMap["sUnion"] = defaultClusterRouter
	routerMap["sDiff"] = defaultClusterRouter
	routerMap["sPop"] = defaultClusterRouter
	return routerMap
}

// defaultClusterRouter 集群间转发的默认方法
func defaultClusterRouter(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.peerPicker.PickNode(key) // PickNode是要找到对应的peer地址
	return cluster.relay(peer, conn, cmdArgs)
}

// LocalRouter 将指令转发到本地
func LocalRouter(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(conn, cmdArgs)
}
