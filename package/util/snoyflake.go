package util

import (
	"douyin/package/constant"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
)

var (
	instance *sonyflake.Sonyflake
)

func init() {
	var err error
	instance, err = sonyflake.New(sonyflake.Settings{
		// 这里若设置成time.Now 那么运行之后就不应该停止
		// 否则可能出现ID重复
		StartTime: time.UnixMilli(constant.SnoyFlakeStartTime),
	})
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

func GetSonyFlakeID() (uint64, error) {
	return instance.NextID()
}

func GetUUid() (string, error) {
	value, err := uuid.NewV4()
	if err != nil {
		zap.L().Sugar().Error(err)
		return "", err
	}
	uuidValue := value.String()
	return uuidValue, nil
}
