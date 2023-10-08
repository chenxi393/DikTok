### TODO
参照gin-mall和其他单体 架构 快速构建一个可用的。 先模仿再改进
之后进行压力测试，然后考虑使用微服务，对比测试（考虑接入AI）
本机运行打开了window的端口 docker内似乎就不用 记得关闭 还做了端口转发
只要接受到请求了就只能返回200 其他状态吗 应该是别的地方做的

### Point
1. 读写分离（基于主从复制） 查询select 在从库 插入更新update在主库
   * 实时性强的，例如即时写即时查，直接指定主库进行操作，避免主从同步延迟而导致查询出bug，即查询比插入快
   * 为什么需要主从复制
       * 一条sql语句可能需要锁表，导致暂时不能使用读的操作，影响业务
       * 做数据的热备？？？
       * 业务量比较大，单机无法满足需求，使用多机器的存储，提高单机的IO性能（单机多库感觉没有意义）
   * [参考配置的文章](https://zhuanlan.zhihu.com/p/650314645)
2. redis各个地方要统筹考虑

### RIGHT JOIN LEFT JOIN WHERE
where和inner join是内连接 只保留公共部分
外连接会保留不满足条件的
总之 外连接至少会保留一张表的所有信息

user表
| id  | name  |
| --- | ----- |
| 1   | Alice |
| 2   | Bob   |
| 3   | Carol |

video表
| id  | title   | author_id |
| --- | ------- | --------- |
| 1   | Video 1 | 1         |
| 2   | Video 2 | 2         |
| 3   | Video 3 | 2         |
| 4   | Video 4 | 4         |

user right join video on author_id=id
| id   | name  | title   |
| ---- | ----- | ------- |
| 1    | Alice | Video 1 |
| 2    | Bob   | Video 2 |
| 2    | Bob   | Video 3 |
| NULL | NULL  | Video 4 |

user right join video等价于video left join user

where author_id=id 或者 直接join
| id   | name  | title   |
| ---- | ----- | ------- |
| 1    | Alice | Video 1 |
| 2    | Bob   | Video 2 |
| 2    | Bob   | Video 3 |

### Fiber
fiber 要注意一个点 Fiber.ctx的值是可变的(会被重复使用-这也是我们是Zero Allocation)

只能在handler里面处理  若要传参或返回值 
得深拷贝copy(buffer, result)  或者调用utils.CopyString(c.Params("foo")) 

也可以配置为不可变（Immutable ） 但是有性能开销
[Zero Allocation](https://docs.gofiber.io/#zero-allocation)


### 可以考虑的 观看青训营答辩得出的 9.28 
1. go-zero进行微服务开发 使用数据库集群存储数据 实现软删除和分页分表 `redis缓存`降低数据库的压力  redis缓存热点数据（注意缓存命中率）---但是感觉只有真实业务才能确定什么是热点 `缓存穿透 缓存击穿 缓存雪崩`
    * 缓存穿透 ---限制ip的访问速率 bucket桶算法？？ 
    * `分布式锁` 解决缓存击穿的问题  --和穿透一回事吧 都是不查缓存去查DB
    * 缓存雪崩（很多key同时过期 ） ---- 解决办法 随机化过期时间
2. 日志进行分级zap比较好用 
3. 提问
    * 为什么用gin不用hertz (技术选型)
    * 缓存什么场景下使用（经常变更的字段）热点数据？？
    * 表的索引设计，消息队列的时效性考虑 
    * 循环查库是bug 不可取
    * 日志需要有一个数据库去归纳 打控制台不太好（日志很重要）
    * 细致的限流操作 具体到ip  Token令牌桶的限流
    * ETCD集群 投票 redia哨兵集群 mysql集群 这是高可用
    * 数据一致性 --订阅 mysql binlog？？？
    * 可以使用大模型做视频推荐系统？？？⭐⭐ 这很好 可以了解
    * 这个项目还是太拉了 不如用抖音那个


* [第一名 必看](https://z37kw7eggp.feishu.cn/docx/Y3KCdaFMSoKKNjxPOHAcWMiInZb)
* [这个是第三名 有详细的开发流程和规范 可以参考](https://gagjcxhxrb.feishu.cn/docx/SCEddZcB3oQwKOxrWQNcqicQnxd)
* 还有两个二等奖 可以看看

### 消息队列的作用
1. `削峰限流` 下游服务器只能承载2000 通过消息队列 高并发慢慢推
2. `解耦` 上下游没有直接的接口调用 下游下线也不会导致整个服务不可用 也可任意新增服务
3. `异步` 上游不需要等待下游返回的结果 可以增加上游的吞吐量

#### 选择合适的消息队列
| 消息队列 | RabbitMQ                        | Kafka                                          | RocketMQ                         |
| -------- | ------------------------------- | ---------------------------------------------- | -------------------------------- |
| 特点1    | 轻量 迅捷     10wQPS            | 依赖ZooKeeper 兼容性好 大数据/日志场景 100wQPS | 延时低 毫秒响应 几十万QPS        |
| 特点2    | 客户端支持丰富                  | 分布式(broker) 性能好 可扩展 可持久化（disk）  | 阿里开发 java开发 二次开发扩展好 |
| 缺点1    | 消息堆积导致性能下降            | 批量发送会导致延时比较高（当消息少的时候）     | 生态不好  兼容性不好             |
| 缺点2    | Erlang开发 扩展和二次开发成本高 | topic上百个吞吐量大幅下降 不适用在线业务       |                                  |

kafaka
* `Topic` (半结构化的数据) topic可以分区 （`partition`）
* `offset` 是消息的位置 每个分区里是唯一且递增的（保障了顺序）
* `record` 消息记录 Key value 键值对
* key值为空 会轮询partition写入 否则相同key的消息可以写到相同的partition
* `Replica-factor` 的副本数量（包括主） 会选取一个`leader` 数据写入和读入都是从leader
* `ISR` 同步副本集，若副本相差比较多 会被踢出ISR的集合（等待副本追赶上）
* Kafka集群由`broker 消息代理`组成 一个服务器启动一个broker实例

#### 问题
* 异步怎么保证 下游一定成功呢 失败了怎么办
* 一致性怎么办 返回给用户成功 怎么保证时效性----所以消息队列是有适用场景的
* 当上游服务器必须等待下游的处理结果返回就不适用
* 消息不丢失？？
* 重复消息