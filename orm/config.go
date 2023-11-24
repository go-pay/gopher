package orm

import (
	"database/sql"

	"github.com/go-pay/xtime"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNoRow          = sql.ErrNoRows
	ErrRecordNotFound = gorm.ErrRecordNotFound

	LogLevelMap = map[string]logger.LogLevel{
		"info":   logger.Info,
		"warn":   logger.Warn,
		"error":  logger.Error,
		"silent": logger.Silent,
	}
)

// MySQLConfig mysql config.
type MySQLConfig struct {
	DSN            string         `json:"dsn" toml:"dsn" yaml:"dsn"`                                        // data source name.
	MaxOpenConn    int            `json:"max_open_conn" toml:"max_open_conn" yaml:"max_open_conn"`          // pool, e.g:10
	MaxIdleConn    int            `json:"max_idle_conn" toml:"max_idle_conn" yaml:"max_idle_conn"`          // pool, e.g:100
	MaxConnTimeout xtime.Duration `json:"max_conn_timeout" toml:"max_conn_timeout" yaml:"max_conn_timeout"` // connect max lifetime. Unmarshal config file e.g: 10s、2m、1m10s
	MaxIdleTimeout xtime.Duration `json:"max_idle_timeout" toml:"max_idle_timeout" yaml:"max_idle_timeout"` // connect max idle time. Unmarshal config file e.g: 10s、2m、1m10s
	LogLevel       string         `json:"log_level" toml:"log_level" yaml:"log_level"`                      // enum: info、warn、error、silent, default warn
	Colorful       bool           `json:"colorful" toml:"colorful" yaml:"colorful"`                         // is colorful log, default false
	SlowThreshold  xtime.Duration `json:"slow_threshold" toml:"slow_threshold" yaml:"slow_threshold"`       // slow sql log. Unmarshal config file e.g: 100ms、200ms、300ms、1s, default 200ms
}
