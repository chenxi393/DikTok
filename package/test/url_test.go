package test

import (
	"net/url"
	"testing"
)

func TestUrlParse(t *testing.T) {
	url, err := url.Parse("nacos:///127.0.0.1:8848?group=DEFAULT_GROUP&clusters=default")
	if err != nil {
		return
	}
	queryParams := url.Query()
	// 输出查询参数
	for key, values := range queryParams {
		for _, value := range values {
			println(key, value)
		}
	}
}
