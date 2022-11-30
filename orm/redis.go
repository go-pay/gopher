package orm

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-pay/gopher/xtime"
	"github.com/go-redis/redis/v8"
)

// RedisConfig redis config.
type RedisConfig struct {
	Addrs        []string       `json:"addrs" yaml:"addrs" toml:"addrs"`
	Password     string         `json:"password" yaml:"password" toml:"password"`
	DB           int            `json:"db" yaml:"db" toml:"db"`
	ReadTimeout  xtime.Duration `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
	WriteTimeout xtime.Duration `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
	TLSCfg       *tls.Config    `json:"-" yaml:"-" toml:"-"`
}

type RedisCluster struct {
	Redis *redis.ClusterClient
}

func InitRedisCluster(c *RedisConfig) (rc *RedisCluster) {
	opts := &redis.ClusterOptions{
		Addrs:    c.Addrs,
		Password: c.Password,
	}
	if c.TLSCfg != nil {
		opts.TLSConfig = c.TLSCfg
	}
	if c.ReadTimeout != 0 {
		opts.ReadTimeout = time.Duration(c.ReadTimeout)
	}
	if c.WriteTimeout != 0 {
		opts.WriteTimeout = time.Duration(c.WriteTimeout)
	}
	rd := redis.NewClusterClient(opts)
	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect redis error:%+v", err))
	}
	rc = &RedisCluster{Redis: rd}
	return rc
}

type Redis struct {
	Redis *redis.Client
}

func InitRedis(c *RedisConfig) (r *Redis) {
	opts := &redis.Options{
		Addr:     c.Addrs[0],
		Password: c.Password,
		DB:       c.DB,
	}
	if c.TLSCfg != nil {
		opts.TLSConfig = c.TLSCfg
	}
	if c.ReadTimeout != 0 {
		opts.ReadTimeout = time.Duration(c.ReadTimeout)
	}
	if c.WriteTimeout != 0 {
		opts.WriteTimeout = time.Duration(c.WriteTimeout)
	}
	rd := redis.NewClient(opts)
	_, err := rd.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect redis error:%+v", err))
	}
	r = &Redis{Redis: rd}
	return r
}

func (r *Redis) Transaction(ctx context.Context, fc func(tx redis.Pipeliner) error) error {
	txPip := r.Redis.TxPipeline()
	defer txPip.Close()
	err := fc(txPip)
	if err != nil {
		txPip.Discard()
		return err
	}
	txPip.Exec(ctx)
	return nil
}

func (r *RedisCluster) Transaction(ctx context.Context, fc func(tx redis.Pipeliner) error) error {
	txPip := r.Redis.TxPipeline()
	defer txPip.Close()
	err := fc(txPip)
	if err != nil {
		txPip.Discard()
		return err
	}
	txPip.Exec(ctx)
	return nil
}
