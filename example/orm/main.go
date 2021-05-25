package main

import (
	"time"

	"github.com/iGoogle-ink/gotil/orm"
	"github.com/iGoogle-ink/gotil/xlog"
	"github.com/iGoogle-ink/gotil/xtime"
)

type MxCity struct {
	Id    int        `gorm:"column:id;primary_key"`
	Ctime xtime.Time `gorm:"column:ctime"`
	Mtime xtime.Time `gorm:"column:mtime"`
}

func main() {

	c := &orm.MySQLConfig{
		DSN:            "uname:password@tcp(host:3306)/db_name?timeout=10s&readTimeout=10s&writeTimeout=10s&parseTime=true&loc=Local&charset=utf8mb4",
		MaxOpenConn:    10,
		MaxIdleConn:    10,
		MaxConnTimeout: xtime.Duration(10 * time.Second),
		ShowSQL:        true,
	}
	db := orm.InitGorm(c)

	var mcs []*MxCity

	err := db.Table("mx_city FORCE INDEX (`idx_mtime`)").Where("mtime < ?", "2020-10-26 10:00:00").Find(&mcs).Error
	if err != nil {
		xlog.Error(err)
		return
	}
	for _, v := range mcs {
		xlog.Debug(v)
	}
}
