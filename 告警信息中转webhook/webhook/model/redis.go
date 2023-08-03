package model

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisConfig struct {
	IP  string `json:"IP,omitempty"`
	PWD string `json:"PWD,omitempty"`
	DB  int    `json:"DB,omitempty"`
}

func (r *RedisConfig) InitRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     r.IP,
		Password: r.PWD,
		DB:       r.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Errorf("redis连接失败")
		panic(err.Error())
	}
	return client

}
