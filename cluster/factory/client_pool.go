package factory

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"simple-godis/resp/client"
)

// ConnectionFactory 连接工厂 提供给ClusterDatabase.peerConnection连接池使用
type ConnectionFactory struct {
	Peer string
}

// MakeObject 新建一个集群之间的连接
func (c *ConnectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	clusterClient, err := client.MakeClusterClient(c.Peer)
	if err != nil {
		return nil, err
	}
	clusterClient.Start()
	return pool.NewPooledObject(clusterClient), nil
}

// DestroyObject 关闭一个集群之间的连接
func (c *ConnectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	clusterClient, ok := object.Object.(*client.ClusterClient)
	if !ok {
		return errors.New("destroy peer connection failed")
	}
	clusterClient.Close()
	return nil
}

func (c *ConnectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (c *ConnectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c *ConnectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
