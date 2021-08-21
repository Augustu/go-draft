--[[
    redis-cli -p 32379 eval "$(cat json.lua)" 1 stats:a:0:b 1619834827 1619834827 10
    redis-cli -p 32379 evalsha d7212d10d43da5e97bd2f680f649b64ec79c6a36 1 stats:a:0:b 161983482 161983482 10
]]

local key = KEYS[1]
local min = ARGV[1]
local max = ARGV[2]
local num = ARGV[3]

local res = redis.call('zrangebyscore', key, min, max)
local data = res[1]
if ( data == nil ) then
        local newmember = cjson.encode({s = num})
        redis.call('zadd', key, min, newmember)
        return 1
end

local json = cjson.decode(data)
json.s = json.s + num

local updated = cjson.encode(json)
redis.call('zadd', key, min, updated)

redis.call('zrem', key, data)
return 0