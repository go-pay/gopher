package orm

import (
	"testing"
	"time"

	"github.com/iGoogle-ink/gopher/xlog"
	"github.com/iGoogle-ink/gopher/xtime"
)

var (
	dsn = "root:root@tcp(mysql:3306)/school?parseTime=true&loc=Local&charset=utf8mb4"
)

type Student struct {
	Id   int    `gorm:"column:id;primary_key" xorm:"'id' pk"`
	Name string `gorm:"column:name" xorm:"'name'"`
}

func (m *Student) TableName() string {
	return "student"
}

func TestInitXorm(t *testing.T) {
	// 初始化 Xorm
	gc1 := &MySQLConfig{DSN: dsn, MaxOpenConn: 10, MaxIdleConn: 10, MaxConnTimeout: xtime.Duration(10 * time.Second), ShowSQL: true}
	x := InitXorm(gc1)

	student := new(Student)
	x.Sync2(student)

	_, err := x.Insert(&Student{Name: "Jerry"})
	if err != nil {
		xlog.Error(err)
		return
	}
	_, err = x.Table(student.TableName()).Select("id,name").Where("id = ?", 1).Get(student)
	if err != nil {
		xlog.Error(err)
		return
	}
	xlog.Debug("xorm:", student)
	x.Close()
}
