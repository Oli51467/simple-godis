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
- Redis集群（待定）

#### 指令

- keys

```keys```
```exists```
```del```
```type```
```rename```
```renamenx```
```flush```

- strings

```get```
```set```
```getset```
```strlen```
```setnx```

- common

```ping```
```exit```

