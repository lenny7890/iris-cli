package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql" // mysql 数据库驱动
	"gorm.io/gorm"         // 使用 gorm ，操作数据库的 orm 框架
	"gorm.io/gorm/logger"  // 使用 gorm 日志
)

var Db *gorm.DB

/*
数据库初始化。
init() 表示包初始化的时候执行的函数, 如果函数名写成 main() , 会在操作数据的时候报错。
参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情。
*/
func init() {

	// ----------------------- 日志设置 -----------------------
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	// ----------------------- 连接数据库 -----------------------
	var err error
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local&timeout=10ms", // DSN data source name
		DefaultStringSize:         256,                                                                                     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true,                                                                                    // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Printf("mysql connect error %v", err)
	}
	if Db.Error != nil {
		fmt.Printf("database error %v", Db.Error)
	}

	// ----------------------- 连接池设置 -----------------------
	sqlDB, err := Db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

}
