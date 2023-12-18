package util

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// 解决grpc返回成功状态码为0 会被忽略 字段为默认零值都会被忽略
func GrpcMarshal(v interface{}) ([]byte, error) {
	data, ok := v.(proto.Message)
	if ok {
		// FIXME protojson 会将64位的整数序列化成 字符串
		return protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}.Marshal(data)
	}
	return json.Marshal(v)
}