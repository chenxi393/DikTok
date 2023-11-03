## 项目结构
* config                配置信息
* database              操作MySQl数据库
* handler               路由的函数
* main.go               程序的入口
* model                 数据库模型
* package               依赖的服务包括redis rabbitmq
* response              返回的数据类型
* router                路由
* service               具体的执行函数
* Dockerfile            web服务的镜像文件
* docker-compose.yaml   docker容器编排
* Makefile              一键部署服务

## 依赖项
* Redis ( single )
* MySQL (master and slave)
* RabbitMQ
* FFmpeg
* Go

## 部署流程
### 自动部署
需要将config.yaml的myIP 替换为自己的ip

根目录下输入 `make up` 一键部署服务

### 手动部署
* 安装上述依赖项，并修改config/config.yaml对应的配置信息
* 需要手动安装FFMpeg 否则视频封面将是默认封面
* 需要配置OSS，七牛云的key 上传到本地后异步上传到七牛云（非必须）
* 完成上述步骤后 `make` 或 `make run` 或 `go run .`启动项目

### 注意事项
* 企业中不使用GORM提供的自动建表migration()，一般手动建表（使用config/douyin.sql建表）
* 若数据库已经存在表的情况下，注释掉database里的init.go的migration()（反之取消注释）
* 否则MySQL主从复制有可能失败
* 请不要上传大于5MB的视频 处理速度会非常慢（大视频理应在客户端压缩）

## 待办
- [ ] 拆分成微服务（考虑go-zero），使用RPC调用
- [ ] 整理日志系统并且接入ELK体系（或者使用OpenTelemetry）
- [ ] 接入视频推荐算法，对用户画像进行刻画
- [ ] 增加视频总结和关键词提取功能
- [ ] 消息模块引入大语言模型 每日定时做视频推荐
- [ ] token的续期 双token？？ 需要客户端支持（延后）
- [x] 主键自增消耗很快 考虑分布式ID生成 snowflake雪花算法？
- [x] 视频搜索功能 全文索引实现
- [ ] 用户名 评论 视频描述 消息敏感词过滤
- [ ] 视频格式和大小 校验（感觉这个应该前端做）
- [ ] 限流操作 redis就行 秒级时间戳当作键 或者token令牌桶（或者具体到ip的限流）
- [ ] 功能测试 性能测试 压力测试


## 下面是开发的杂乱笔记

### 统计代码行数
```
wc -l `find ./ -name "*.go";find -name "*.yaml"`
```

### 数据库设计
总体思路
* 不建议使用NULL 为什么 占用NULL表？？
全都设定为NOT NULL 给默认值
* 若有删除功能 添加is_delete (tinyint 1) 作软删除 TODO
* 不适用外键或者级联 使用逻辑外键 似乎影响数据库的插入速度
* 使用SQLyog检查 确保没有多余无用的索引

comment 主键索引
* 根据video_ID 查找评论 建立video_id普通索引
* 根据user_id 和 video_id 删除评论  user_id 普通索引
* 软删除？？

message 主键索引
* 建立联合索引user_id to_user_id
* create_time 建立普通索引

follow 主键索引
* 取消关注时 user_id to_user_id 建立为联合唯一索引
* 查询粉丝 需要 to_user_id 建立普通索引
* 查询关注的 需要user_id  复用上面的联合索引
* 软删除？？

user 主键索引
* username 唯一索引

favorite 主键索引
* user_id video_id 联合唯一索引  user_id 在前（查询喜欢的视频）
* is_deleted 软删除

video 主键索引
* authorid 普通索引
* publish_time 普通索引 < 小于好像用不到索引（-1 然后<=）

### redis 缓存设计  解决读的问题
整体思路：
* 读请求量大的热点数据 缓存到redis里降低数据库的压力
* 对于不同的模块划分不同的redis数据库 降低不同业务数据的影响
* 热点数据过期时间长 冷数据过期时间短
* 可以通过日志分级 去测量一下缓存命中率 99.9%是好的

缓存雪崩----大量缓存同时过期 大量请求之间访问数据库
* 设置缓存的时候给缓存过期时间加上一个随机数 降低缓存同时过期的概率
* 在特定场景下对缓存进行预热 预先加载一部分大概率用到的数据到内存

缓存穿透----大量请求查询一个`不存在于缓存中的key`
* 使用`布隆过滤器` 由于请求参数多为userID videoID
* 使用两个布隆过滤器存储userID videoID
* 在用户注册 发布视频 把id加入布隆过滤器
* 限制ip的访问速率 bucket桶算法？？
* 回种空值？？

缓存击穿----大量请求查询`一个过期的key`
* 使用redis的SetNX 实现`分布式锁`
* 请求缓存数据加锁 缓存不存在 继续更新锁的有效时间
* 保证同一时间只有一个协程在进行SQL查询
* 执行SQL查询更新缓存后释放锁
* 后续的请求就可以走缓存了

为什么分布式锁有效 
* 因为redis是单线程 参考一下别人的实现TODO

具体模块设计
1. comment
   * 方案一：（代码量太大 分开过期不好搞）
   * key: comment_id value: json序列化后的结构体 （评论结构体string）
   * key: video_id   value: int 评论id   （视频的评论id zset 使用主键id排序）
   * 方案二：（可能会有大key问题 本项目采用的）
   * 直接用zset存一个视频的所有评论（删除是数据库层面）
2. video 
   * key: video_id  value: json序列化后的结构体 （视频固定字段string）
   * key: video_id  value: count        （视频的count的计数字段 hash）
   * 分开计数信息和详细信息防止大key对哈希性能的影响？？
   * key: user_id value：video_id1 video_id2 用户publish的视频ID  （zset）
   * key: user_id value: video_id1 video_id2 用户favorite的视频ID （zset）
3. relation 和user使用一个redis库
   * 用户的关注列表userID集合 userID---[user_id .....] zset
   * 用户的粉丝列表userID集合                          zset
   * 这里要注意的大V的粉丝可能几百万 但是不一定就存几百万 粉丝列表加载限制一次加载100
4. user
   * 用户的个人信息的不变字段（基本不变的字段）  string
   * 用户的个人信息的计数字段（变化多的计数字段）  hash
   * 因为hash更新字段 会重新分配内存 拆一下 

数据结构选用原则
* 容易变更的字段选用hash结构
* 固定的字段可以使用JSON序列化存string
* 排序集合使用Zset 无序使用set
* 注意video和user的信息会个别过期 需要一组user或者video信息的时候
* 应当先去redis查出id 再查redis存在的信息，不存在的记录下来批量查数据库

缓存预热
* 在登录成功的时候 就把用户相关的数据加入缓存
* feed 视频流 异步将视频数据加入的缓存中？？

数据一致性 [推荐的文章](https://www.cnblogs.com/crazymakercircle/p/14853622.html)
* 使用GROM的事务进行数据更新，引入lua脚本（redis的原子性），把更新redis当作数据库事务的一部分
* 查询时，先查redis，没有查数据库并加入缓存（可以异步）
* 当修改数据库时，为了确保数据一致性，需要`删除`掉redis的脏数据！！！（目前的方案）
* 更新缓存需要去数据库里拿所有数据不能增量更新 即使这样也有数据不一致的风险 两个请求 先更新数据库的后到redis执行 待验证是否可以增量或者全量更新 TODO
* 或增量更新前判断键存不存在 引入lua脚本 如果放在GORM事务里应该没有数据数据不一致的风险 因为会等redis执行成功再返回 两个事务的并发情况分析？？ 也就是mysql事务和redis lua原子操作打包成一个整体执行TODO 
    * 引入消息队列Canal订阅binlog 收到增量更改时删除redis
    * 先更新MySQL，再删除（修改）redis的key（本项目采用的）
    * 缓存延时双删

### 消息队列设计  解决写的问题
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
2. 注意这里视频和封面（ffmeng截取第一帧 只在docker内部署了）的存储
   * 全都写在本地 且是当前目录下 这不合适 而且封面截取 需要ffmeng
   * 使用对象存储（七牛云每月有免费额度） 对象存储需要自己的域名 CDN是默认开启的
3. 只要服务端收到了请求了就只能返回200
4. 生成唯一ID的三种算法
   * 主键自增
       * 有序 存在数量泄漏风险 但是可以优化主键自增法
   * UUID生成（基于MAC地址和时间和随机数）
       * 无序字符串
   * 雪花算法（基于相对时间69年。机器码，序号）
       * 按时间趋势递增
       * 分布式系统中不会ID碰撞
       * 时间回拨会乱序和重复

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
| id  | name  | title   |
| --- | ----- | ------- |
| 1   | Alice | Video 1 |
| 2   | Bob   | Video 2 |
| 2   | Bob   | Video 3 |

### Fiber
fiber 要注意一个点 Fiber.ctx的值是可变的(会被重复使用-这也是我们是Zero Allocation)
例如 result := c.Params("foo")  result可能会被修改 
尽管go的string被认为值类型不可变，但是实际上可以修改底层的字节数据改变
string只能在handler里面有效  若要传参或返回值 
得深拷贝copy(buffer, result)  或者调用utils.CopyString(c.Params("foo")) 

也可以配置为不可变（Immutable） 但是有性能开销
[Zero Allocation](https://docs.gofiber.io/#zero-allocation)

### 遇到的问题 
1. MySQL 主从同步 1032 error 主库用来update，从库同来select
   * [[MySQL] SQL_ERROR 1032解决办法](https://www.cnblogs.com/langdashu/p/5920436.html)
   * 解决办法就是 查看 日志 插入重复失败 就删除 删除失败就插入 但是为什么会重复插入啊 明明已经插入了 `找到不同步的点` 让他们同步 ？？ 
   * 查看binlog 但是这个问题老是出现 出现的原因是什么 应该时GORM的自动迁移导致的 
   * [错误复现](https://cloud.tencent.com/developer/article/1564571)
   * show binlog events in 'binlog.000004';
2. GORM Scan的两个问题 Scan的要求类型是什么，它是如何匹配相应字段的
   * [GORM Scan源码](https://blog.csdn.net/xz_studying/article/details/107095153)

### 不同阶段的测试集合
其他测试 压力测试 接入redis 测耗时 接入消息队列 测耗时
#### 测试1
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
规划之初设定的被清晰定义的、可实现的且可测量的项目成功标准

目标一 完成数据库设计 done
1. 做好数据库的字段设计
2. 设计好索引结构，外键逻辑
3. 对数据长度，量级等等进行分析，预估数据瓶颈

目标二 完成后端基本功能开发 done
1. 完成基本模块功能开发
2. 完成互动模块功能开发
3. 完成社交模块功能开发

目标三 完成 Redis 缓存设计 done
1. 完善缓存基本设计
2. 解决三大经典缓存问题
3. 结合场景对缓存进行预热处理

目标四 完成消息队列部分 done
1. 完成消息队列选型和需要接入的接口分析
2. 完成消息队列代码接入
3. 完成消息队列相关代码的性能测试（整个图）

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
3. 仿照真的抖音，多做关注和朋友动态两个页面