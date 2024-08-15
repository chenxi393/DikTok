# 指定基础镜像，必须为第一个命令
FROM golang:alpine AS builder

# LABEL 为镜像添加标签
LABEL stage=gobuilder

# ADD 将本地文件添加到容器里 为tar类型会自动解压

# 解决go镜像下载慢的问题
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# Install ffmpeg 提取视频第一帧 弃用
# RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /diktok

# 把本地文件拷贝到容器里 这里应该是到工作目录
COPY . .

# 提取参数 构建对应的镜像
# 包变小了 但是会很卡 因为并行的跑几个 之前只跑一个
ARG SERVICE
RUN if [ "$SERVICE" = "gateway" ]; then \
        go install /diktok/gateway; \
        else \
        go install /diktok/service/$SERVICE; \
    fi


FROM alpine

WORKDIR /app

# 设置时区环境变量
ENV TZ=Asia/Shanghai
ENV RUN_ENV=docker

# 安装 tzdata 包以支持时区
RUN apk update && apk add --no-cache tzdata

COPY --from=builder /diktok/config /app/config
COPY --from=builder /go/bin/ /app/

# 容器运行时执行的shell 命令 一般在最后一行 一定要前台运行 不然运行之后容器就关闭了
# 可以被docker run 的tag覆盖
# ENTRYPOINT 容器启动时执行的入口
# EXPOSE 暴露端口 
# VOLUME 加载卷