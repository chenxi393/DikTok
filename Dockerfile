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

WORKDIR /douyin

# 把本地文件拷贝到容器里 这里应该是到工作目录
COPY . .

# 构建镜像运行的shell 命令
RUN go install

FROM alpine

WORKDIR /app

COPY --from=builder /douyin/config /app/config
COPY --from=builder /go/bin/douyin /app/douyin

# 容器运行时执行的shell 命令 一般在最后一行 一定要前台运行 不然运行之后容器就关闭了
# 可以被docker run 的tag覆盖
# ENTRYPOINT 容器启动时执行的入口
# EXPOSE 暴露端口 
# VOLUME 加载卷
CMD /app/douyin