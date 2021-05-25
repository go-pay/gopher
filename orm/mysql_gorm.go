package orm

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitGorm(c *MySQLConfig) (db *gorm.DB) {
	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database error:%+v", err))
	}
	db.DB().SetMaxOpenConns(c.MaxOpenConn)
	db.DB().SetMaxIdleConns(c.MaxIdleConn)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxConnTimeout))
	db.LogMode(c.ShowSQL)
	return db
}
