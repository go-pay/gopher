package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig redis config.
type RedisConfig struct {
	Addrs        []string
	Password     string
	DB           int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type RedisCluster struct {
	Redis *redis.ClusterClient
}

func InitRedisCluster(c *RedisConfig) (rc *RedisCluster) {
	opts := &redis.ClusterOptions{
		Addrs:    c.Addrs,
		Password: c.Password,
	}
	if c.ReadTimeout != 0 {
		opts.ReadTimeout = c.ReadTimeout
	}
	if c.WriteTimeout != 0 {
		opts.WriteTimeout = c.WriteTimeout
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
	if c.ReadTimeout != 0 {
		opts.ReadTimeout = c.ReadTimeout
	}
	if c.WriteTimeout != 0 {
		opts.WriteTimeout = c.WriteTimeout
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
