FROM golang:latest

# Ignore APT warnings about not having a TTY
ENV DEBIAN_FRONTEND noninteractive

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Install ffmpeg 提取视频第一帧
RUN apt-get install -y ffmpeg

WORKDIR /go/projects/douyin
COPY . .

RUN go install
CMD /go/bin/douyin