## DikTok main
主分支为微服务版本，正在开发中

分支v1为单体服务
## 项目结构
|              |                   |
| ------------ | ----------------- |
| config       | 公共配置信息      |
| gateway      | API网关           |
| grpc         | gRPC生成代码      |
| idl          | protobuf接口定义  |
| model        | 数据库模型        |
| package      | 公共依赖包        |
| service      | 拆分的微服务      |
| storage      | 数据库 缓存 MQ    |
| compose.yaml | docker容器编排    |
| Dockerfile   | web服务的镜像文件 |
| Makefile     | 一键部署服务      |

## 依赖项
* Redis ( single )
* MySQL ( master and slave )
* RabbitMQ
* FFmpeg (弃用)
* Go

## 部署流程
### 自动部署
根目录下输入 `make up` 一键部署服务

### 手动部署
* 安装上述依赖项，并修改config/config.yaml对应的配置信息
* 需要配置OSS，七牛云的key
* 分别运行网关和6个微服务（可以使用Makefile）

### 注意事项
* 不推荐使用GORM自动建表migration()，一般手动建表（使用config/douyin.sql建表）
* 若数据库已经存在表的情况下，注释掉database里的init.go的migration()（反之取消注释）
* MySQL主从复制在docker容器重启时有可能失败
* 请不要上传大于30MB的视频 会返回413（大视频应在客户端压缩）

## 待办
- [x] 拆分成微服务，考虑gRPC+ETCD，再考虑成熟的微服务框架go-zero
- [x] 分布式ID生成 snowflake雪花算法
- [x] 视频搜索功能 MySQL全文索引实现（可考虑ES）
- [x] 接入OpenTelemetry，完成traces，metric的上报
- [ ] 消息模块引入大语言模型√ 每日定时做视频推荐（定时任务怎么写？）
- [ ] token的续期 双token？？
- [ ] 消息模块使用MongoDB存储，消息的全文搜索（考虑ES实现）
- [ ] 项目快速部署和运维的探究 K8s CICD体系
- [ ] 接入视频推荐算法（Gorse），对用户画像进行刻画
- [ ] 增加视频总结和关键词提取功能（大模型）
- [ ] 完善Websocket替换消息模块轮询
- [ ] 视频格式大小校验√ 评论敏感词过滤，视频水印生成（FFmpeg）
- [ ] redis一主两从哨兵模式 MySQL双主互从+Keepalived（redis和MySQL集群引入）
- [ ] 限流操作 redis就行 秒级时间戳当作键 或者token令牌桶（或者具体到ip的限流）
- [ ] 功能测试√ 性能测试 压力测试

## 下面是开发的杂乱笔记
根据那四篇文章改造自己的抖音系统 推拉模型 feed流

可以用redis set实现抽奖功能 
spop 可以用来抽 一等 二等（会删除）
srandmenber 可以一次抽取（不会删除）

### 微服务拆分
`protoc --go_out=.. --go-grpc_out=.. ./idl/video.proto`
生成两个pb文件

注意Grpc 生成go结构体时，由于proto3 去除了require字段
导致JSON tag 默认 omitempty
解决办法
* 不使用encode/json 使用别的json库（目前没有找到替代的）
* 或者使用sed替换 grpc生成的代码 `sed -i "" -e "s/,omitempty//g" ./api/proto/*.go`
* 在网关层面 不直接使用proto生成的结构体返回（感觉这样更好）
```go
w.Header().Set("Content-Type", "application/json; charset=utf-8")
    m := protojson.Marshaler{EmitDefaults: true}
    m.Marshal(w, resp) // You should check for errors here
```

bytedance/sonic 未支持此特性 考虑提个PR？？？


TODO
目前微服务拆分了，但是缓存redis，数据库，消息队列还存在耦合，需要再拆分
目前的想法 服务之间的依赖关系全都走rpc或者消息队列
数据库的字段也经历减少耦合

注意微服务不一定每个服务都使用一个DB或者缓存
可以多个服务共享 取决于业务的耦合程度

微服务适用于大型团队 更好的独立的进行多模块快速开发

单体对一块内容的更改需要重新部署整个服务，而微服务不用

微服务可以按需水平扩展（部署服务的多个实例）

实际上可以通过MQ进行服务间的通信（）

服务架构强调服务之间的无事务协调

### FIXME 前端用户评论操作后，立即请求加载评论列表，出现拿到脏数据
原因：
* 评论操作为消息队列异步执行，写入数据库之前就返回前端成功
* 主从同步存在延时会导致，读操作在从库拿，写操作在主库

若读请求在消息队列写入数据库之前，必然拿到脏数据（如何解决？？？）
若读请求在写入数据库之后，若去从库拿数据，有概率拿到脏数据（指定主库读可以解决）

目前暂时的解决办法是：前端等待几毫秒毫秒再去请求评论列表

### 统计代码行数
```sh
wc -l `find ./ -name "*.go";find -name "*.yaml"`
```

### 数据库设计
总体思路
* 全都设定为NOT NULL（不占用NULL表空间）
* 若有删除功能 添加is_delete (tinyint 1) 作软删除（或者deleted_at） TODO
* 不适用外键或者级联 使用逻辑外键 外键影响数据库的插入速度
* 确保没有多余无用的索引（可以考虑使用工具检查）
* 要合理的使用组合索引，而不是单列索引。

comment 主键索引
* 根据video_ID 查找评论 建立video_id普通索引

message 主键索引
* 联合索引user_id to_user_id
* （可以考虑MongoDB 来存message）

follow 主键索引  TODO 这里可以考虑联合唯一主键 因为id自增没啥用
* user_id to_user_id 联合唯一索引
* 查询粉丝 需要 to_user_id 建立普通索引
* （查询关注的 需要user_id  复用上面的联合索引）

user 主键索引
* username 唯一索引

favorite 主键索引 TODO 这里可以考虑联合唯一主键 因为id自增没啥用
* user_id video_id 联合唯一索引  user_id 在前（查询喜欢的视频）

video 主键索引
* authorid 普通索引
* publish_time 普通索引 普通索引< >是正常走索引的

### Redis缓存设计
整体思路：
* 对于不同的模块划分不同的redis数据库 降低不同业务数据的影响
* 热点数据过期时间长 冷数据过期时间短
* 可以通过日志分级 去测量一下缓存命中率 99.9%是好的
* 读多写少考虑旁路缓存策略，延时双删，或canal（TODO）
* 读多写多场景（例如点赞），直接在Redis完成自增操作（并持久化），定期同步到数据库TODO

缓存雪崩----大量缓存同时过期 大量请求之间访问数据库
* 设置缓存的时候给缓存过期时间加上一个随机数 降低缓存同时过期的概率
* 在特定场景下对缓存进行预热 预先加载一部分大概率用到的数据到内存

缓存穿透----大量请求查询一个`不存在于缓存中的key`
* 使用`布隆过滤器` 由于请求参数多为userID videoID
* 使用两个布隆过滤器存储userID videoID
* 在用户注册 发布视频 把id加入布隆过滤器
* 限制ip的访问速率 令牌桶算法TODO
* 存储空值（不如布隆过滤器）
    * 注意是缓存查不到数据库也查不到再去 存空值

缓存击穿----大量请求查询`一个过期的key`（或者说少量的key）
* 使用redis的SetNX 实现`分布式锁`
* 设置锁过期时间，值为uuid
* 缓存不存在 尝试加锁
* 执行SQL查询更新缓存后释放锁
* 后续的请求就可以走缓存了
* 分布式锁 还需要续期TODO

具体模块设计
1. comment
   * 方案一：（代码量太大 分开过期不好搞）
   * key: comment_id value: json序列化后的结构体 （评论结构体string）
   * key: video_id   value: int 评论id   （视频的评论id zset 使用主键id排序）
   * 方案二：（可能会有大key问题 本项目采用的）
   * 直接用zset存一个视频的所有评论（删除是数据库层面）
2. video（读多写少） 
   * key: video_id  value: json序列化后的结构体 （视频固定字段string）
   * key: video_id  value: count        （视频的count的计数字段 hash）
   * 分开计数信息和详细信息防止大key对哈希性能的影响？？
   * key: user_id value：video_id1 video_id2 用户publish的视频ID  （zset）
   * key: user_id value: video_id1 video_id2 用户favorite的视频ID （zset）
3. relation
   * 用户的关注列表userID集合 userID---[user_id .....] zset
   * 用户的粉丝列表userID集合                          zset
   * 这里要注意的大V的粉丝可能几百万 但是不一定就存几百万 粉丝列表加载限制一次加载100
4. user
   * 用户的个人信息的不变字段（基本不变的字段）  string
   * 用户的个人信息的计数字段（变化多的计数字段）  hash
   * 因为hash更新字段 会重新分配内存 拆一下 

数据结构选用原则
* 容易变更的字段选用hash结构
* 固定的字段可以使用JSON序列化存string（可以使用MsgPack）
* 应当先去redis查出id 再查redis存在的信息，不存在的记录下来批量查数据库

数据一致性 [推荐的文章](https://www.cnblogs.com/crazymakercircle/p/14853622.html)
* 读操作：先读缓存，缓存没有，查数据库，（异步）写入缓存（这叫旁路缓存）
* 写操作：先写数据库，再删缓存，但是需要保证两个操作的完整性（GORM保证）
    * 使用GROM的事务进行数据更新，引入lua脚本（redis的原子性），把更新redis当作数据库事务的一部分
    * 引入消息队列Canal订阅binlog 收到增量更改时删除redis
* 写操作：缓存延时双删（先删缓存 再更新数据库 睡一会 再删除缓存）

⭐无论是MQ重试还是延迟双删或是Canal等方法，在极端条件下都有会不一致的情况，只能尽量降低不一致性的可能（可以通过队列保证请求时序性）

如果业务对缓存命中率很高，可以采用[更新数据库]+[更新缓存]的方案
* 更新缓存有数据不一致的风险 两个请求 先更新数据库的后到redis执行
* 更新前加入分布式锁，引入lua脚本，放在GORM事务里，等redis执行成功再返回

实际上双写不一致使用一个队列，或者加分布式锁 都可以
只需要保证它的执行的先后即可

### 点赞系统的设计（读多写多） 参考B站实现（也可以看看得物的文章）
[得物的点赞方案](https://xie.infoq.cn/article/f6840380238de0761abe39e08)
目前我的方案有点像得物的1.0方案 比较鸡肋
[B站的技术方案](https://mp.weixin.qq.com/s/4T_S7nR8-HXJ59IbK4FBWQ)
#### 点赞服务系统能力分析 
业务需求分析：
* 对视频（取消）点赞
* 查询视频点赞数
* 查询是否对单个视频点赞
* 用户点赞列表
* 用户收到的总点赞数
* 视频的点赞列表（暂时客户端没有）

容灾能力
* DB不可用（依托缓存）
* 缓存不可用（尽量依靠DB）
* 消息队列不可用（RPC降级）
* 数据同步延迟（不一致）
* 点赞消息堆积

#### 点赞压力分析
全局流量压力
* 写流量比较大，写时可以在内存中聚合数据
* 例如聚合10S的点赞数，一次性写入，减少IO
* 异步写入数据库，可以用消息队列（削峰）
* 点赞的正确性保证，不能重复点赞

单点压力
* 热门视频需要考虑单点压力可以缓存到内存中

#### 存储分析
持久层（第一层）
* 点赞记录表 userID videoID （点赞时间 来源）
* 点赞数表 放在视频表里 videoID 
* 数据量大考虑分库分表（怎么分？）可以使用TiDB（分布式数据库） 

Redis层（第二层）
* 点赞数 string存 key: favorite:count:{video_id} value: likes
* 用户点赞列表 ZSet key: user:favorite:userID value: member:{video_ID} score{timestamp}

本地缓存（第三次）暂未实现 为了应对热点问题

#### 异步任务（主要针对写入）
* 保证数据的写入不会超过数据库的负荷同时也不会出现数据堆积（导致查询延迟）

* 可以先写本地缓存，或者先写redis然后定期同步到MySQL
* 还可以秒级聚合所有请求 同步到数据库
* 旁路缓存的策略不太适合读多写多的场景

### 消息队列设计 TODO再优化使用
为什么要使用消息队列
* 热门视频会有大量的点赞和评论操作 写入MySQL压力大
* 关注和发送消息流量也比较大

好处
* 异步 
    * 之前是串行处理 写入数据库 redis 提取视频封面 上传到OSS 时间是叠加的
    * 使用消息队列可以异步 类似于并行处理 大大减少了响应时间
    * 异步下游服务不可用 例如上传OSS不可用 不影响上游服务
* 解耦
    * 上面的不可用也是解耦的好处
    * 如果我们要新增视频主题提取 视频总结 情色分析 发布订阅模型 订阅对应的主题即可
    * 不用影响原来的业务 解耦合
* 削峰
    * 每天的业务不是均衡的 只能承载1000  但是流量来了10000
    * 消息队列暂存请求

坏处
* 消息队列存在写入数据库不成功或者不及时的场景
* 注册和发视频 比较重要 不适用消息队列（不能失败）
* 点赞 关注 评论 数据库写入失败对用户影响小 异步写入 

### Point
1. 读写分离（基于主从复制） 查询select 在从库 插入更新update在主库
   * 主从同步存在延时时间 主库更新 从备库读到脏数据怎么解决
       * 实时性强的，例如即时写即时查，直接指定主库进行操作，避免主从同步延迟而导致查询出bug，即查询比插入快
   * 为什么需要主从复制
       * 一条sql语句可能需要锁表，导致暂时不能使用读的操作，影响业务
       * 做数据的热备？？？
       * 业务量比较大，单机无法满足需求，使用多机器的存储，提高单机的IO性能（单机多库感觉没有意义）
   * [参考配置的文章](https://zhuanlan.zhihu.com/p/650314645)
2. 只要服务端收到了请求了就只能返回200
3. 生成唯一ID的三种算法
   * 主键自增
       * 有序 存在数量泄漏风险 但是可以优化主键自增法
   * UUID生成（基于MAC地址和时间和随机数）
       * 无序字符串
   * 雪花算法（基于相对时间69年。机器码，序号）
       * 按时间趋势递增
       * 分布式系统中不会ID碰撞
       * 时间回拨会乱序和重复

### Fiber
fiber 要注意一个点 Fiber.ctx的值是可变的(会被重复使用-这也是我们是Zero Allocation)
例如 result := c.Params("foo")  result可能会被修改 
尽管go的string被认为值类型不可变，但是实际上可以修改底层的字节数据改变
string只能在handler里面有效  若要传参或返回值 
得深拷贝copy(buffer, result)  或者调用utils.CopyString(c.Params("foo")) 

也可以配置为不可变（Immutable） 但是有性能开销
[Zero Allocation](https://docs.gofiber.io/#zero-allocation)

### 搜索功能的实现探究
* MySQL的[全文索引](https://blog.csdn.net/mrzhouxiaofei/article/details/79940958)
* ElasticSearch
* 都是基于倒排索引

### websocket
HTTP 长轮询（数据更新之后才回复）VS短轮询（有轮询间隔） 

长轮询
* 客户端发起一个请求，服务器收到客户端发来的请求后，服务器端不会直接进行响应，而是先将这个请求挂起，然后判断请求的数据是否有更新。如果有更新，则进行响应，如果一直没有数据，则等待一定的时间后才返回。

websocket 应用层协议
ws默认端口80  wss端口为443（TLS）

优点
* 开销小：连接建立后，交换数据头部小
* 实时性：全双工协议，服务器可以主动发给客户端
* 长连接：有状态的协议，连接建立后可以省略部分状态信息
* 二进制支持：定义了二进制帧，HTTP应该是文本吧
* 可扩展

websocket 生命周期
* 客户端HTTP请求GET/ws
* 服务端发出握手升级HTTP为Websocket
* 互相通信
* 一方关闭连接则连接关闭

#### 流程
客户端请求：
```http
GET ws://echo.websocket.org/ HTTP/1.1
Host: echo.websocket.org
Origin: file://
Connection: Upgrade
Upgrade: websocket
Sec-WebSocket-Version: 13
Sec-WebSocket-Key: Zx8rNEkBE4xnwifpuh8DHQ==
Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
```

服务端响应：
```http
HTTP/1.1 101 Web Socket Protocol Handshake
Connection: Upgrade
Upgrade: websocket
Sec-WebSocket-Accept: 52Rg3vW4JQ1yWpkvFlsTsiezlqw=
```
101 表示未完成的连接

可以利用HTTP服务器根据具体流程实现websockt服务器

websocket特点（主要对比HTTP）
* 与HTTP一样默认端口80/443 握手阶段采用HTTP
* 数据格式轻量，开销小
* 可以二进制也可以发送文本
* 没有同源现在
* 协议标识ws/wss

TODO 目前怎么推送消息还没有完成，至少客户端不用HTTP轮询了
而是改用websocket轮询

### 遇到的问题 
1. MySQL 主从同步 出现的问题
   * 出现问题的根本原因：容器重启 导致同步的进度被刷新 容器重启relay log不一致（官方已知bug）
   * [[MySQL] SQL_ERROR 1032解决办法](https://www.cnblogs.com/langdashu/p/5920436.html)
   * [错误复现](https://cloud.tencent.com/developer/article/1564571)
   * show binlog events in 'binlog.000004';
   * 新的错误 ERROR 1872 (HY000): Slave failed to initialize relay log info structure from the repository
   * docker重新启动 导致主机名变化
   * 通用的解决办法 找到同步进度（binlog的位置很重要） 手动重置
2. GORM Scan的两个问题 Scan的要求类型是什么，它是如何匹配相应字段的
   * [GORM Scan源码](https://blog.csdn.net/xz_studying/article/details/107095153)

### 不同阶段的测试集合
功能测试 性能测试 压力测试

简单压力测试下遇到的问题
1. "Error 1040: Too many connections" 数据库会返回这个 设置比较高的连接数 似乎也会这样

#### 功能测试和简单的延迟测试
条件：docker-compose 部署 开启mysql主从复制 没有redis 没有消息队列

APIfox返回的时间
```t
/douyin/user/register/  66ms
/douyin/user/login/  57ms
/douyin/user/ 5ms
/douyin/publish/action/ 1MB 的视频 180ms（偶尔蹦到很高）  3MB 的视频 在270ms  （似乎第一次发会蹦到很高） 接口时间随视频大小上升 500KB的视频 只需要140ms
/douyin/feed/ 登录7ms 未登录6ms
/douyin/publish/list/ 7ms 视频数量176个

/douyin/favorite/action/ 15ms
/douyin/favorite/list/ 6ms
/douyin/comment/action/ 14ms
/douyin/comment/list/ 6ms

/douyin/relation/action/ 13ms
/douyin/relation/follow/list/ 6ms
/douyin/relation/follower/list/ 6ms
/douyin/relation/friend/list/ 6ms

/douyin/message/action/  11ms
/douyin/message/chat/ 4ms
```

fiber返回的时间  这是在docker内
```t
/douyin/user/register/  62ms
/douyin/user/login/  54ms
/douyin/user/ 0.8ms
/douyin/publish/action/ 1MB 的视频 145ms（偶尔蹦到很高）  3MB 的视频 在230ms  （似乎第一次发会蹦到很高） 接口时间随视频大小上升 500KB的视频 只需要130ms
/douyin/feed/ 登录3ms 未登录1.8ms
/douyin/publish/list/ 3ms 视频数量176个

/douyin/favorite/action/ 13ms
/douyin/favorite/list/ 3ms
/douyin/comment/action/ 11ms
/douyin/comment/list/ 2ms

/douyin/relation/action/ 12ms
/douyin/relation/follow/list/ 2.2ms
/douyin/relation/follower/list/ 2.0ms
/douyin/relation/friend/list/ 2.7ms

/douyin/message/action/  8ms
/douyin/message/chat/ 1.2ms
```

### 项目规划
目标一 完成数据库设计 done
1. 做好数据库的字段设计
2. 设计好索引结构，外键逻辑
3. 对数据长度，量级等等进行分析，预估数据瓶颈

目标二 完成后端基本功能开发 done
1. 完成基本模块功能开发
2. 完成互动模块功能开发
3. 完成社交模块功能开发

目标三 完成 Redis 缓存设计 done
1. 完善缓存基本设计（读多写少 读多写多场景）
2. 解决三大经典缓存问题
3. 结合场景对缓存进行预热处理

目标四 完成消息队列部分 done
1. 完成消息队列选型和需要接入的接口分析
2. 完成消息队列代码接入
3. 完成消息队列相关代码的对比性能测试

目标五 完成全面的测试
1. 完成代码部分的集成测试
2. 完成Apifox的功能测试
3. 造好测试场景用的数据，并进行一定程度上的压测

目标六 完善运维体系
1. 实现 CI/CD
2. 实现链路追踪

目标七 提高项目安全质量
1. 完善用户密码相关（加密，复杂度校验等）
2. 完成接口限流，部分接口进行防刷设置
3. 继续进行安全优化（敏感词检查，文件检查等）

目标八 进行功能扩展和其他部分完善
1. 完成基本的推荐算法以及测试
2. 添加部分常见功能接口（如修改个人信息等）
3. 仿照真的抖音，给自己提需求