package redisService

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var (
	GlobalRedis *redis.Client
)

func InitRedis() {

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis Connect failed", err)
		return
	} else {
		log.Println("Redis Connect Succeed", result)
	}
	GlobalRedis = client
}

type RedisService struct {
	PreKey      string
	Ctx         context.Context
	GlobalRedis interface{}
}

// SetNX 包含了设置过期时间，原子操作
func (r *RedisService) SetNX(key string, value interface{}, exp time.Duration) *redis.BoolCmd {
	return GlobalRedis.SetNX(r.Ctx, key, value, exp)
}

func (r *RedisService) LPush(key string, value ...interface{}) *redis.IntCmd {
	return GlobalRedis.LPush(r.Ctx, key, value)
}

// SAdd 向Set集合里添加成员
func (r *RedisService) SAdd(key string, members ...interface{}) {
	GlobalRedis.SAdd(r.Ctx, key, members)
}

// SIsMember 判断该集合中是否存在member
func (r *RedisService) SIsMember(key string, member interface{}) *redis.BoolCmd {
	return GlobalRedis.SIsMember(r.Ctx, key, member)
}

// SRemove 删除集合中对应members的value
func (r *RedisService) SRemove(key string, members ...interface{}) {
	GlobalRedis.SRem(r.Ctx, key, members)
}

// SGet 取出集合中的所有元素
func (r *RedisService) SGet(key string) *redis.StringSliceCmd {
	return GlobalRedis.SMembers(r.Ctx, key)
}

// Get 用于获取对应键的值
func (r *RedisService) Get(key string) *redis.StringCmd {
	return GlobalRedis.Get(r.Ctx, key)
}

func (r *RedisService) SCard(key string) *redis.IntCmd {
	return GlobalRedis.SCard(r.Ctx, key)
}

func (r *RedisService) Del(key string) *redis.IntCmd {
	return GlobalRedis.Del(r.Ctx, key)
}
