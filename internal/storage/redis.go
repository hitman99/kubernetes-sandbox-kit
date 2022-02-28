package storage

import (
	"github.com/go-redis/redis/v7"
	"github.com/hitman99/kubernetes-sandbox/internal/config"
)

type RedisClient interface {
	LRange(key string, start, stop int64) *redis.StringSliceCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
}

func MustNewRedisClient() RedisClient {
	cfg := config.GetRedisConfig()
	if !cfg.IsCluster {
		return redis.NewClient(&redis.Options{
			Addr: cfg.Address,
		})
	} else {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{cfg.Address},
		})
	}
}

const REDIS_LIST_KEY = "users"

func WrapRedisErr(err error) error {
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
