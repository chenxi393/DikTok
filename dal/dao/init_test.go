package dao_test

import (
	"douyin/config"
	"douyin/dal/dao"
	"testing"
)

func TestInit(t *testing.T) {
	config.Init()
	dao.InitMysql()
}
