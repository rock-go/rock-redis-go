-- redis 模块,启动脚本
local redis = rock.redis {
    name = "redis",
    addr = "192.168.3.71:6379",
    password = "localtest",
    db = 0,
    pool_size = 100,
    max_conn_age = 100,
}

proc.start(redis)
