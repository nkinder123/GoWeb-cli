package redisCon

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

// 声明一个全局的rdb变量
var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化连接
func Init() (err error) {
	addstr := fmt.Sprintf("%s:%d",
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"))
	client = redis.NewClient(&redis.Options{
		Addr:     addstr,
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"), // use default DB
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return
}

func Close() {
	client.Close()
}

func GetRedis() *redis.Client {
	return client
}
