package database_test

import (
	"douyin/config"
	"douyin/database"
	"testing"
)

func TestInit(t *testing.T) {
	config.Init()
	database.InitMysql()
}
