package orm

import (
	"github.com/go-pay/gopher/xtime"
)

// MySQLConfig mysql config.
type MySQLConfig struct {
	DSN            string         `json:"dsn" yaml:"dsn" toml:"dsn"`                                        // data source name.
	MaxOpenConn    int            `json:"max_open_conn" yaml:"max_open_conn" toml:"max_open_conn"`          // pool, e.g:10
	MaxIdleConn    int            `json:"max_idle_conn" yaml:"max_idle_conn" toml:"max_idle_conn"`          // pool, e.g:100
	MaxConnTimeout xtime.Duration `json:"max_conn_timeout" yaml:"max_conn_timeout" toml:"max_conn_timeout"` // connect max life time. Unmarshal config file e.g:10s、2m、1m10s
	ShowSQL        bool           `json:"show_sql" yaml:"show_sql" toml:"show_sql"`
}
