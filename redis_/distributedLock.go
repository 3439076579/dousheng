package redisService

import (
	"awesomeProject/utils"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

/*
	本文件用于构建分布式锁
	分布式锁：
		- 高可用
		- 独占性
		- 可重入
 		- 防死锁
*/

type RedisDistributeLock struct {
	key          string
	Ctx          context.Context
	uuID         int64
	ScriptLock   *redis.Script
	ScriptUnLock *redis.Script
	GoroutinueID int64
	expireTime   int64
	Keys         []string
	Argv1        string
	Argv2        int64
}

func (r RedisDistributeLock) Lock() {
	for IsLock, _ := r.ScriptLock.Run(r.Ctx, GlobalRedis, r.Keys, r.Argv1, r.Argv2).Bool(); !IsLock; {
		time.Sleep(20 * time.Millisecond)
	}

}

func (r RedisDistributeLock) Unlock() error {
	flag := r.ScriptUnLock.Run(r.Ctx, GlobalRedis, r.Keys, r.Argv1)

	if flag.Err() == redis.Nil {
		return errors.New("try to unlock a lock isn't locked")
	}

	return nil

}
func GetDistributedLock(expireTime int64) RedisDistributeLock {
	distributedLock := RedisDistributeLock{
		Ctx:  context.Background(),
		key:  "RedisLock",
		uuID: utils.GenerateUuid(),
		ScriptLock: redis.NewScript(`
if redis.call('EXISTS',KEYS[1]) == 0 then
    -- 加锁，设置过期时间
    redis.call('HINCRBY',KEYS[1],ARGV[1],1)
    redis.call('EXPIRE',KEYS[1],ARGV[2])
    return 1
    -- 如果已经存在该锁，体现可重入性
else if redis.call('HEXISTS',KEYS[1],ARGV[1]) then
    redis.call('HINCRBY',KEYS[1],ARGV[1],1)
    return 1
else
    return 0
end
end`),
		ScriptUnLock: redis.NewScript(`
if redis.call('HEXISTS',KEYS[1],ARGV[1])==0 then
    return nil
elseif redis.call('HINCRBY',KEYS[1],ARGV[1],-1) ==0 then
    return redis.call('DEL',KEYS[1])
    
else 
    return 0
    
end `),
		GoroutinueID: utils.GetGoID(),
		expireTime:   expireTime,
	}
	distributedLock.Keys = append(distributedLock.Keys, distributedLock.key)
	distributedLock.Argv1 = strconv.FormatInt(distributedLock.GoroutinueID, 10) + ":" +
		strconv.FormatInt(distributedLock.uuID, 10)

	distributedLock.Argv2 = distributedLock.expireTime

	return distributedLock

}
