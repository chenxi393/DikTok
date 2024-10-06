## DikTok
main分支为微服务版本，持续开发中

分支v1为单体服务版本（不再维护）
## 项目结构
|              |                    |
| ------------ | ------------------ |
| config       | 公共配置信息       |
| gateway      | API网关服务        |
| grpc         | gRPC生成代码       |
| idl          | rpc服务接口定义    |
| package      | 公共依赖包         |
| service      | RPC服务            |
| storage      | MySQL Redis MQ     |
| compose.yaml | docker容器编排     |
| Dockerfile   | 后端服务的镜像文件 |
| Makefile     | 开发部署脚本       |

## 依赖项
* Redis ( single )
* MySQL ( master and slave )
* FFmpeg (弃用)
* Go

## 部署流程
### 自动部署
修改config/config.yaml对应的配置信息
根目录下输入 `make up` 一键部署服务

### 手动部署
* 安装上述依赖项，并修改config/config.yaml对应的配置信息
* 需要配置OSS，七牛云的key
* 使用Makefile分别运行API网关和微服务

### 注意事项
* 请不要上传大于30MB的视频 会返回413（大视频应在客户端压缩）
* 视频支持客户端获取token上传OSS

## 待办
- [x] 拆分成微服务，gRPC服务间内网通信，ETCD服务注册与服务发现
- [x] 分布式ID生成 snowflake雪花算法
- [x] 视频搜索功能 MySQL全文索引实现（可考虑ES复杂条件搜筛）
- [x] 接入OpenTelemetry，完成traces，metrics，logs的上报
- [ ] 评论中台建设，二级评论，点赞，二级评论默认展示首条 其余展开分页5条
- [ ] IM中台建设，Websocket替换轮询，消息模块使用MongoDB存储，消息的全文搜索（考虑ES实现）
- [ ] 系统通知 被关注通知 被点赞 被评论 被艾特通知 
- [ ] 消息模块引入大语言模型√ 每日定时做视频推荐（定时任务怎么写？）
- [ ] 增加视频总结和关键词提取功能（大模型）AI视频推荐助手Agent落地实践
- [ ] token的续期 双token？？
- [ ] 项目快速部署和运维 K8s CICD体系
- [ ] 接入视频推荐算法（Gorse），对用户，视频，画像刻画
- [ ] 视频格式大小校验√ 评论敏感词过滤，视频水印生成（FFmpeg）
- [ ] redis一主两从哨兵模式 MySQL双主互从+Keepalived（redis和MySQL集群引入）
- [ ] 限流操作 redis就行 秒级时间戳当作键 或者token令牌桶（或者具体到ip的限流）
- [ ] 功能测试√ 性能测试 **压力测试**
