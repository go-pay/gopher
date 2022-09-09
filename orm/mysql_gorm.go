package orm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(c *MySQLConfig) (db *gorm.DB) {
	lc := logger.Config{
		SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:                  logger.Warn,            // Log level
		Colorful:                  c.Colorful,             // 日志颜色
		IgnoreRecordNotFoundError: false,                  // 忽略记录未找到错误
	}
	if c.LogLevel != "" {
		if ll, ok := LogLevelMap[c.LogLevel]; ok {
			lc.LogLevel = ll
			if lc.LogLevel == logger.Error || lc.LogLevel == logger.Silent {
				lc.IgnoreRecordNotFoundError = true // 忽略记录未找到错误
			}
		}
	}
	if c.SlowThreshold != 0 {
		lc.SlowThreshold = time.Duration(c.SlowThreshold)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] >> ", log.Lmsgprefix|log.Ldate|log.Lmicroseconds), // io writer
		lc,
	)
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{Logger: newLogger, SkipDefaultTransaction: true})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database error:%+v", err))
	}
	sql, err := db.DB()
	if err != nil {
		panic(err)
	}
	sql.SetMaxIdleConns(c.MaxIdleConn)
	sql.SetMaxOpenConns(c.MaxOpenConn)
	sql.SetConnMaxLifetime(time.Duration(c.MaxConnTimeout))
	sql.SetConnMaxIdleTime(time.Duration(c.MaxIdleTimeout))
	return db
}
