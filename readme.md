# 说明

rock-redis-go模块基于rock-go框架开发,用于连接和操作redis

# 使用

该模块主要用于日志技术分析中.

## 导入

```go
import redis "github.com/rock-go/rock-redis-go"
```

## 组件注册

```go
rock.Inject(xcall.Rock, redis.LuaInjectApi)
```

## lua 脚本调用
调用时，一般声明为一个全局变量redis，供数据分析模块使用。
```lua
-- redis 模块
redis = rock.redis {
    name = "redis",
    addr = "192.168.3.71:6379",
    password = "localtest",
    db = 0,
    pool_size = 100,
    max_conn_age = 100,
}

proc.start(redis)

-- redis 操作

-- hmset: 设置test_key的test_field字段的值为1
redis.hmset("test_key", "test_field", 1)

-- hmget: 获取test_key的test_field字段的值
local v = redis.hmget("test_key", "test_field")
print(v) -- 1

-- incr: 对test_key的test_filed字段加1操作
redis.incr("test_key", "test_field")
local v = redis.hmget("test_key", "test_field")
print(v) -- 2

-- hdelete: 删除test_key的test_field字段
redis.hmdel("test_key", "test_field")

-- delete: 删除test_key
redis.delete("test_key")

-- pipeline 操作,批量执行操作
local pipeline = redis.pipeline()
redis.hmset("test_key", "test_field", 1)
redis.hmset("test_key_1", "test_field", 1)
redis.hmset("test_key_2", "test_field", 1)
redis.hmset("test_key_3", "test_field", 1)
redis.hmset("test_key_3", "test_field_1", 1)
pipeline.exec()
pipeline.close()

local pipeline_1 = redis.pipeline()
pipeline_1.incr("test_key", "test_field")
pipeline_1.incr("test_key_1", "test_field")
pipeline_1.incr("test_key_2", "test_field")
pipeline_1.incr("test_key_3", "test_field")
pipeline_1.exec()
pipeline_1.close()

local pipeline_2 = redis.pipeline()
pipeline_2.hmdel("test_key_3", "test_field_1")
pipeline_2.hmdel("test_key_2", "test_field")
pipeline_2.delete("test_key_1")
pipeline_2.exec()
pipeline_2.close()
```