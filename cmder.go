package redis

import (
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
)

type cmder func(redis.Cmdable , *lua.LState) int

func hmset( r redis.Cmdable, L *lua.LState) int {
	n := L.GetTop()
	if n < 3 {
		L.RaiseError("hmset(key , ...)  #args must >= 3")
		return 0
	}

	if (n - 1) % 2 != 0 {
		L.RaiseError("hmset(key , key1 , val1 , key2 , val2) must be name kv")
		return 0
	}

	key := L.CheckString(1)
	fields := make(map[string]interface{})
	for i := 1 ; i<= n; i += 2 {
		field := L.CheckString(i + 1)
		val := L.Get(i+2)
		switch val.Type() {
		case lua.LTNumber:
			fields[field] = int(val.(lua.LNumber))
		case lua.LTString:
			fields[field] = val.String()
		default:
			L.RaiseError("hmset must value must be int or string")
			return 0
		}
	}

	var err error
	//var res string
	_, err = r.HMSet(key, fields).Result()
	if err != nil {
		L.RaiseError("redis hmset error: %v", err)
		return 0
	}

	return 0
}

func hmget( r redis.Cmdable, L *lua.LState) int {
	var err error
	var res []interface{}
	res, err = r.HMGet(L.CheckString(1), L.CheckString(2)).Result()
	if err != nil {
		L.RaiseError("redis hmget error: %v", err)
		return 0
	}

	if len(res) == 0 {
		return 0
	}

	switch v := res[0].(type) {
	case int:
		L.Push(lua.LNumber(v))
	case string:
		L.Push(lua.LString(v))
	default:
		return 0
	}

	return 1
}

func incr(r redis.Cmdable, L *lua.LState) int {
	var err error

	_, err = r.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
	if err != nil {
		L.RaiseError("redis hincr error: %v", err)
	}

	return 0
}

func expire(r redis.Cmdable , L *lua.LState) int {
	var err error

	_, err = r.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
	if err != nil {
		L.RaiseError("redis hincr error: %v", err)
	}

	return 0
}

func del(r redis.Cmdable , L *lua.LState) int {
	var err error
	_, err = r.Del(L.CheckString(1)).Result()
	if err != nil {
		L.RaiseError("redis delete error: %v", err)
	}
	return 0
}

func hdel(r redis.Cmdable , L *lua.LState) int {
	var err error
	_, err = r.HDel(L.CheckString(1), L.CheckString(2)).Result()
	if err != nil {
		L.RaiseError("redis hdelete error: %v", err)
	}
	return 0
}