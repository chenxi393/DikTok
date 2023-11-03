package database

import (
	"douyin/config"
	"douyin/model"
	"douyin/package/constant"
	"strings"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func InitMySQL() {
	var ormLogger logger.Interface
	//根据配置文件设置不同的日志等级
	if config.System.Mode == constant.DebugMode {
		ormLogger = logger.Default.LogMode(logger.Info) // Info应该是最低的等级 都会打印
	} else { // default 的慢sql是200ms
		ormLogger = logger.Default // 进去看 这里是Warn级别的
	}
	// dsn := "用户名:密码@tcp(地址:端口)/数据库名"
	masterDNS := strings.Join([]string{
		config.System.MysqlMaster.UserName,
		":",
		config.System.MysqlMaster.Password,
		"@tcp(",
		config.System.MysqlMaster.Host,
		":",
		config.System.MysqlMaster.Port,
		")/",
		config.System.MysqlMaster.Database,
		"?charset=utf8mb4&parseTime=True&loc=Local"}, "",
	)
	slaveDNS := strings.Join([]string{
		config.System.MysqlSlave.UserName,
		":",
		config.System.MysqlSlave.Password,
		"@tcp(",
		config.System.MysqlSlave.Host,
		":",
		config.System.MysqlSlave.Port,
		")/",
		config.System.MysqlSlave.Database,
		"?charset=utf8mb4&parseTime=True&loc=Local"}, "",
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       masterDNS,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置

	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //单数表
		},
	})
	if err != nil {
		zap.L().Fatal("MySQL 连接失败", zap.Error(err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(config.System.MysqlMaster.MaxOpenConn) // 设置数据库最大连接数
	sqlDB.SetMaxIdleConns(config.System.MysqlMaster.MaxIdleConn) // 设置上数据库最大闲置连接数
	// 查询在从库完成，其他操作如写入update在主库操作
	err = db.Use(dbresolver.Register(dbresolver.Config{
		//Sources:  []gorm.Dialector{mysql.Open(masterDNS)}, // update使用 这里应该是默认连接主库
		Replicas: []gorm.Dialector{mysql.Open(slaveDNS)}, // select 使用
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		zap.L().Error("MySQL 读写分离创建失败", zap.Error(err))
	}
	// 连接池什么的不懂 先放着
	constant.DB = db
	//migration()
}

// 企业一般不用自动建表 记得自己在主库里建表
func migration() {
	err := constant.DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Video{},
		&model.Favorite{},
		&model.Comment{},
		&model.Message{},
	)
	if err != nil {
		zap.L().Fatal("数据库migration失败", zap.Error(err))
	}
}
