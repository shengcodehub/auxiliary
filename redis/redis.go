package redis

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

var (
	// redisCache redis cache
	redisCache *redis.Redis
)

type NodeConf struct {
	Addr      string `json:"addr"`
	Password  string `json:"password"`
	DB        int    `json:"db"`
	KeyPrefix string `json:"keyPrefix"` //前缀key
}

func SetCache(c cache.ClusterConf) {
	redisCache = NewRds(c)
}

func GetCache() *redis.Redis {
	return redisCache
}

func NewRds(cacheRedis cache.ClusterConf) *redis.Redis {

	conf := redis.RedisConf{
		Host:        cacheRedis[0].Host,
		Type:        cacheRedis[0].Type,
		Pass:        cacheRedis[0].Pass,
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	return rds
}
