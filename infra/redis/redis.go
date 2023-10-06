package redisUtils

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/swarajkumarsingh/job-data-digger/conf"
)

var (
	ctx = context.Background()
	Rdb *redis.Client
)

func Init() {
	enableSSL, _ := conf.RedisConf["SSL"].(bool)
	endpoint, _ := conf.RedisConf["Addr"].(string)
	userName, _ := conf.RedisConf["Username"].(string)
	password, _ := conf.RedisConf["Password"].(string)
	redisOption := &redis.Options{
		Addr:         endpoint,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     20,
		Username:     userName,
		Password:     password,
		PoolTimeout:  30 * time.Second,
	}
	if enableSSL {
		redisOption.TLSConfig = &tls.Config{}
	}
	Rdb = redis.NewClient(redisOption)

	if val := Rdb.Ping(ctx); val != nil {
		if val.Val() == "PONG" {
			fmt.Println("Connected to redis successfully")
			return
		}
		fmt.Println("error connecting to redis client: ", val)
	}
}
