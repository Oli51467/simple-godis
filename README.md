### 使用Go语言实现Redis

- 通信协议：```Redis Serialization Protocol```
- Go版本 ```1.17```

#### 技术栈

- Epoll
- TCP
- GC
- HashMap
- 协程
- 锁
- defer

#### 实现

- 使用TCP协议的应用层服务器
- Redis协议解析器
- 内存数据库
- Redis持久化
- Redis集群

#### 指令

- Keys

```keys```
```exists```
```del```
```type```
```rename```
```renamenx```
```flush```

- Strings

```get```
```set```
```getset```
```strlen```
```setnx```
```getdel```
```append```
```incr```
```decr```

- Common

```ping```
```exit```

- Set

```sAdd```
```sIsMember```
```sRem```
```sMembers```
```sCard```
```sInter```
```sUnion```
```sDiff```
```sPop```

- List

```lpush```
```lpushx```
```rpush```
```rpushx```
```lpop```
```rpop```
```lindex```
```lset```
```lrange```
```lrem```
```llen```

- Hash

```hset```
```hsetnx```
```hget```
```hdel```
```hexists```
```hmset```
```hmget```
```hkeys```
```hvalues```
```hgetall```
```hlen```
```hstrlen```
