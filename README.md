## DikTok
main分支为微服务版本，持续开发中

分支v1为单体服务（不再维护）
## 项目结构
|              |                   |
| ------------ | ----------------- |
| config       | 公共配置信息      |
| gateway      | API网关服务       |
| grpc         | gRPC生成代码      |
| idl          | rpc服务接口定义   |
| model        | 数据库模型        |
| package      | 公共依赖包        |
| service      | 拆分的RPC服务     |
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
* MySQL主从复制在docker容器重启时偶现失败
* 请不要上传大于30MB的视频 会返回413（大视频应在客户端压缩）

## 待办
- [x] 拆分成微服务，gRPC服务间内网通信，ETCD服务注册与服务发现
- [x] 分布式ID生成 snowflake雪花算法
- [x] 视频搜索功能 MySQL全文索引实现（可考虑ES扩展多维度信息搜索（大模型））
- [x] 接入OpenTelemetry，完成traces，metrics，logs的上报
- [ ] 评论的回复（二级评论） 评论点赞  二级评论默认展示首条 其余展开逻辑
- [ ] 系统通知 被关注通知 被点赞 被评论 被 艾特通知 
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
