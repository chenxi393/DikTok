package database

import (
	"douyin/config"
	"douyin/model"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var global_db *gorm.DB

func InitMysql() {
	var ormLogger logger.Interface
	if config.SystemConfig.Mode == "debug" { //根据配置文件设置不同的日志等级
		ormLogger = logger.Default.LogMode(logger.Info) // Info应该是最低的等级 都会打印
	} else {
		ormLogger = logger.Default // 进去看 这里是Warn级别的
	}
	// dsn := "用户名:密码@tcp(地址:端口)/数据库名"
	masterDNS := strings.Join([]string{
		config.SystemConfig.MysqlMaster.UserName,
		":",
		config.SystemConfig.MysqlMaster.Password,
		"@tcp(",
		config.SystemConfig.MysqlMaster.Host,
		":",
		config.SystemConfig.MysqlMaster.Port,
		")/",
		config.SystemConfig.MysqlMaster.Database,
		"?charset=utf8mb4&parseTime=True&loc=Local"}, "",
	)
	slaveDNS := strings.Join([]string{
		config.SystemConfig.MysqlSlave.UserName,
		":",
		config.SystemConfig.MysqlSlave.Password,
		"@tcp(",
		config.SystemConfig.MysqlSlave.Host,
		":",
		config.SystemConfig.MysqlSlave.Port,
		")/",
		config.SystemConfig.MysqlSlave.Database,
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
		panic(err.Error())
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(config.SystemConfig.MysqlMaster.MaxOpenConn) // 设置数据库最大连接数
	sqlDB.SetMaxIdleConns(config.SystemConfig.MysqlMaster.MaxIdleConn) // 设置上数据库最大闲置连接数
	// point:读写分离
	// 查询在从库完成，其他操作如写入update在主库操作
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDNS)}, // update使用
		Replicas: []gorm.Dialector{mysql.Open(slaveDNS)},  // select 使用
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	// 连接池什么的不懂 先放着
	global_db = db
	// 自动建表 企业一般不用 这里为了方便 就不手动建表了
	migration()
}

func migration() {
	err := global_db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Video{},
		&model.UserFavoriteVideo{},
	)
	if err != nil {
		panic(err.Error())
	}
}
