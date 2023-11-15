
---使用hash结构 HashMap<String,HashMap<String,Object>>
--首先判断是否有对应的Keys，如果没有，就可以加锁

if redis.call('EXISTS',KEYS[1]) == 0 then

    -- 加锁，设置过期时间
    redis.call('HINCRBY',KEYS[1],ARGV[1],1)
    redis.call('EXPIRE',KEYS[1],ARGV[2])
    return 1
    -- 如果已经存在该锁，体现可重入性
elseif redis.call('HEXISTS',KEYS[1],ARGV[1]) then
    redis.call('HINCRBY',KEYS[1],ARGV[1],1)
    return 1
else
    return 0
end