package database

import (
	"strings"

	"diktok/config"
	"diktok/package/constant"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

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
		zap.L().Fatal("MySQL主库连接: 失败", zap.Error(err))
	}
	zap.L().Info("MySQL主库连接: 成功")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(config.System.MysqlMaster.MaxOpenConn) // 设置数据库最大连接数
	sqlDB.SetMaxIdleConns(config.System.MysqlMaster.MaxIdleConn) // 设置上数据库最大闲置连接数
	// 查询在从库完成，其他操作如写入update在主库操作
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDNS)}, // update使用 这里应该是默认连接主库
		Replicas: []gorm.Dialector{mysql.Open(slaveDNS)},  // select 使用
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	if err != nil {
		// 主从创建失败 此时不应该写入数据 应该让容器重启的 否则只会写入主库 导致主从不同步
		zap.L().Fatal("MySQL 读写分离创建失败", zap.Error(err))
	}
	zap.L().Info("MySQL从库连接: 成功")
	// 连接池什么的不懂 先放着
	DB = db
	//migration()
}

// 弃用自动建表
// func migration() {
// 	err :=DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
// 		&model.User{},
// 		&model.Follow{},
// 		&model.Video{},
// 		&model.Favorite{},
// 		&model.Comment{},
// 		&model.Message{},
// 	)
// 	if err != nil {
// 		zap.L().Fatal("数据库migration失败", zap.Error(err))
// 	}
// }
