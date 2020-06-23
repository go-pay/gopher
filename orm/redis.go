package orm

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// redis redis config.
type RedisConfig struct {
	Addrs    []string
	Password string
	DB       int
}

type RedisCluster struct {
	Redis *redis.ClusterClient
}

func InitRedisCluster(c *RedisConfig) (rc *RedisCluster) {
	rd := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    c.Addrs,
		Password: c.Password,
	})
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
	rd := redis.NewClient(&redis.Options{
		Addr:     c.Addrs[0],
		Password: c.Password,
		DB:       c.DB,
	})
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
