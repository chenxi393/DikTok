package constant

import "time"

// 中间件相关
const (
	TokenTimeOut    = 12 * time.Hour
	TokenMaxRefresh = 3 * time.Hour
)