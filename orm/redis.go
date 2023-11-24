package orm

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-pay/xtime"
	"github.com/redis/go-redis/v9"
)

// RedisConfig redis config.
type RedisConfig struct {
	Addrs        []string                                        `json:"addrs" yaml:"addrs" toml:"addrs"`
	Username     string                                          `json:"username" yaml:"username" toml:"username"`
	Password     string                                          `json:"password" yaml:"password" toml:"password"`
	DB           int                                             `json:"db" yaml:"db" toml:"db"`
	ReadTimeout  xtime.Duration                                  `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout xtime.Duration                                  `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	PoolSize     int                                             `json:"pool_size" yaml:"pool_size" toml:"pool_size"`
	MaxIdleConn  int                                             `json:"max_idle_conn" yaml:"max_idle_conn" toml:"max_idle_conn"`
	TLS          bool                                            `json:"tls" yaml:"tls" toml:"tls"`
	TLSCfg       *tls.Config                                     `json:"-" yaml:"-" toml:"-"`
	Limiter      redis.Limiter                                   `json:"-" yaml:"-" toml:"-"`
	OnConnFunc   func(ctx context.Context, cn *redis.Conn) error `json:"-" yaml:"-" toml:"-"`
}

func InitRedis(c *RedisConfig) (rd *redis.Client) {
	opts := &redis.Options{
		Addr:         c.Addrs[0],
		OnConnect:    c.OnConnFunc,
		Username:     c.Username,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		ReadTimeout:  time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
		TLSConfig:    c.TLSCfg,
	}
	rd = redis.NewClient(opts)
	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect redis error:%+v", err))
	}
	return rd
}

func InitRedisCluster(c *RedisConfig) (rc *redis.ClusterClient) {
	opts := &redis.ClusterOptions{
		Addrs:        c.Addrs,
		OnConnect:    c.OnConnFunc,
		Username:     c.Username,
		Password:     c.Password,
		PoolSize:     c.PoolSize,
		ReadTimeout:  time.Duration(c.ReadTimeout),
		WriteTimeout: time.Duration(c.WriteTimeout),
		TLSConfig:    c.TLSCfg,
	}
	rc = redis.NewClusterClient(opts)
	_, err := rc.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect redis error:%+v", err))
	}
	return rc
}
