package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	redisdb *redis.Client
)

func InitRedisClient() (err error) {

	configfile := GetYml()

	redisdb = redis.NewClient(&redis.Options{
		Addr:     configfile.Redis.Host + ":" + configfile.Redis.Port,
		Password: configfile.Redis.Passowrd, // no password set
		DB:       configfile.Redis.Db,       // use default DB
		PoolSize: configfile.Redis.Poolsize, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = redisdb.Ping(ctx).Result()
	return err
}
func RedisSet(key string, val string, ext time.Duration) {
	ctx := context.Background()
	if err := InitRedisClient(); err != nil {
		return
	}

	err := redisdb.Set(ctx, key, val, ext).Err()
	if err != nil {
		fmt.Println(err.Error())
	}

}

func RedisGet(key string) string {

	ctx := context.Background()
	val, err := redisdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
