package util

import (
	"time"

	"github.com/sony/sonyflake"
	"go.uber.org/zap"
)

var (
	instance *sonyflake.Sonyflake
)

func init() {
	var err error
	instance, err = sonyflake.New(sonyflake.Settings{
		StartTime: time.Now(),
	})
	if err != nil {
		zap.L().Fatal(err.Error())
	}
}

func GetSonyFlakeID() (uint64, error) {
	return instance.NextID()
}
