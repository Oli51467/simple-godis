package database

import (
	dbInterface "simple-godis/interface/database"
	"simple-godis/interface/resp"
	"simple-godis/resp/reply"
)

type ClusterDatabase struct {
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
