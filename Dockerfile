FROM golang:latest

# 设置工作目录
WORKDIR /app

# 拷贝本地源代码到容器内
COPY . .

# 安装依赖
RUN go mod download

# 构建可执行文件
RUN go build -o main .

# 容器启动时运行可执行文件
ENTRYPOINT ["./main"]