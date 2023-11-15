for _, v in ipairs(ARGV) do
    redis.call("SETBIT",KEYS[1],v,1)
end