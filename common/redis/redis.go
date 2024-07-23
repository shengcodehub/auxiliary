package redis

import (
	"fmt"

	_redis "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"

	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/redis"
	"github.com/spf13/viper"
)

var redisClient *_redis.Client
var redisClientAuth *_redis.Client
var _redLock *redsync.Redsync

const prefix = "__"

func SetupWithConfig(redisConfig redis.Rdconfig) {

	client, err := redis.NewClient(&redisConfig)
	if err != nil {
		log.Fatalf("open redis failed %s", err.Error())
	}

	redisClient = client
}

func Setup() {
	redisConfig := redis.Rdconfig{
		DB:       viper.GetInt("Redis.Default.DB"),
		Addr:     viper.GetString("Redis.Default.Addr"),
		Password: viper.GetString("Redis.Default.Password"),
	}
	SetupWithConfig(redisConfig)
	setupAuthRedis()
	setupLock()
}

func setupAuthRedis() {
	redisConfig := redis.Rdconfig{
		DB:       viper.GetInt("Redis.Auth.DB"),
		Addr:     viper.GetString("Redis.Auth.Addr"),
		Password: viper.GetString("Redis.Auth.Password"),
	}

	client, err := redis.NewClient(&redisConfig)
	if err != nil {
		log.Fatalf("open redis failed %s", err.Error())
	}

	redisClientAuth = client
}

func GetCache() *_redis.Client {
	return redisClient
}

func GetCacheAuth() *_redis.Client {
	return redisClientAuth
}

// 加载分布式锁
func setupLock() {
	pool := goredis.NewPool(GetCache())
	_redLock = redsync.New(pool)
}

func GetLock() *redsync.Redsync {
	return _redLock
}

func Keys(key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}
