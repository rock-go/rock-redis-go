package redis

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"time"
)

func LuaInjectApi(env xcall.Env) {
	env.Set("redis", lua.NewFunction(createRedisUserData))
}

func createRedisUserData(L *lua.LState) int {
	opt := L.CheckTable(1)
	cfg := Config{
		name:       opt.CheckString("name", "redis"),
		addr:       opt.CheckSockets("addr", L),
		password:   opt.CheckString("password", ""),
		db:         opt.CheckInt("db", 0),
		poolSize:   opt.CheckInt("pool_size", 10),
		maxConnAge: opt.CheckInt("max_conn_age", 10),
	}

	redis := &Redis{C: cfg}

	var obj *Redis
	var ok bool

	proc := L.NewProc(redis.C.name)
	if proc.Value == nil {
		proc.Value = redis
		goto done
	}

	obj, ok = proc.Value.(*Redis)
	if !ok {
		L.RaiseError("invalid redis proc")
		return 0
	}
	obj.C = cfg

done:
	L.Push(proc)

	return 1
}

func (r *Redis) start(L *lua.LState) int {
	if err := r.Start(); err != nil {
		L.RaiseError("redis start error")
		return 0
	}

	return 0
}

func (r *Redis) close(L *lua.LState) int {
	if err := r.Close(); err != nil {
		L.RaiseError("redis close error")
		return 0
	}

	return 0
}

func (r *Redis) LHMSet(L *lua.LState) int {
	var err error
	//var res string
	fields := make(map[string]interface{})
	fields[L.CheckString(2)] = L.CheckInt(3)
	_, err = r.client.HMSet(L.CheckString(1), fields).Result()
	if err != nil {
		L.RaiseError("redis hmset error: %v", err)
		return 0
	}

	return 0
}

func (r *Redis) LHMGet(L *lua.LState) int {
	var err error
	var res []interface{}
	res, err = r.client.HMGet(L.CheckString(1), L.CheckString(2)).Result()
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

func (r *Redis) LIncr(L *lua.LState) int {
	var err error

	_, err = r.client.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result()
	if err != nil {
		L.RaiseError("redis hincr error: %v", err)
	}

	return 0
}

func (r *Redis) LExpire(L *lua.LState) int {
	var err error
	_, err = r.client.Expire(L.CheckString(1), time.Duration(L.CheckInt(2))*time.Second).Result()
	if err != nil {
		L.RaiseError("redis set expire error: %v", err)
	}

	return 0
}

func (r *Redis) LDelete(L *lua.LState) int {
	var err error
	_, err = r.client.Del(L.CheckString(1)).Result()
	if err != nil {
		L.RaiseError("redis delete error: %v", err)
	}

	return 0
}

func (r *Redis) LHDelete(L *lua.LState) int {
	var err error
	_, err = r.client.HDel(L.CheckString(1), L.CheckString(2)).Result()
	if err != nil {
		L.RaiseError("redis hdelete error: %v", err)
	}

	return 0
}

// LNewPipe 新建一个pipeline
func (r *Redis) LNewPipe(L *lua.LState) int {
	pipe := Pipe{
		pipe: r.client.Pipeline(),
	}

	L.Push(&lua.LightUserData{Value: &pipe})
	return 1
}

func (r *Redis) Index(L *lua.LState, key string) lua.LValue {
	if key == "start" {
		return lua.NewFunction(r.start)
	}
	if key == "close" {
		return lua.NewFunction(r.close)
	}
	if key == "hmset" {
		return lua.NewFunction(r.LHMSet)
	}
	if key == "hmget" {
		return lua.NewFunction(r.LHMGet)
	}
	if key == "hmdel" {
		return lua.NewFunction(r.LHDelete)
	}
	if key == "incr" {
		return lua.NewFunction(r.LIncr)
	}
	if key == "expire" {
		return lua.NewFunction(r.LExpire)
	}
	if key == "delete" {
		return lua.NewFunction(r.LDelete)
	}
	if key == "pipeline" {
		return lua.NewFunction(r.LNewPipe)
	}

	return lua.LNil
}

func (r *Redis) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		r.C.name = lua.CheckString(L, val)
	case "addr":
		r.C.addr = lua.CheckString(L, val)
	case "password":
		r.C.password = lua.CheckString(L, val)
	case "db":
		r.C.db = lua.CheckInt(L, val)
	case "pool_size":
		r.C.poolSize = lua.CheckInt(L, val)
	case "max_conn_age":
		r.C.maxConnAge = lua.CheckInt(L, val)
	}
}
