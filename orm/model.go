package orm

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

var (
	ErrNoRow                = sql.ErrNoRows
	ErrRecordNotFound       = gorm.ErrRecordNotFound
	ErrCantStartTransaction = gorm.ErrCantStartTransaction
	ErrInvalidSQL           = gorm.ErrInvalidSQL
	ErrInvalidTransaction   = gorm.ErrInvalidTransaction
	ErrUnaddressable        = gorm.ErrUnaddressable
)

// MySQLConfig mysql config.
type MySQLConfig struct {
	DSN         string // data source name.
	Active      int    // pool
	Idle        int    // pool
	IdleTimeout int    // connect max life time. second
	ShowSQL     bool
}
