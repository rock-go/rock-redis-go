package redis

import (
	"github.com/rock-go/rock/lua"
	"time"
	"reflect"
	"github.com/rock-go/rock/xcall"
)

var TRedis = reflect.TypeOf((*Redis)(nil)).String()

func newLuaRedis(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name , TRedis)
	if proc.IsNil() {
		proc.Set(newRedis(cfg))
	} else {
		proc.Value.(*Redis).cfg = cfg
	}
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

	L.Push(L.NewAnyData(pipe))
	return 1
}

func (r *Redis) Index(L *lua.LState, key string) lua.LValue {
	var lv lua.LValue
	lv = r.meta.Get(key)
	if lv != lua.LNil {
		return lv
	}

	switch key {
	case "start":
		lv = L.NewFunction(r.start)
	case "close":
		lv = L.NewFunction(r.close)
	case "hmset":
		lv = L.NewFunction(r.LHMSet)
	case "hmget":
		lv = L.NewFunction(r.LHMGet)
	case "incr":
		lv = L.NewFunction(r.LIncr)
	case "expire":
		lv = L.NewFunction(r.LExpire)
	case "delete":
		lv = L.NewFunction(r.LDelete)
	case "pipeline":
		lv = L.NewFunction(r.LNewPipe)
	default:
		L.RaiseError("%s redis %s not found" , r.Name() , key)
		return lua.LNil
	}

	r.meta.Set(key , lv)
	return lv
}

func (r *Redis) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		r.cfg.name = lua.CheckString(L, val)
	case "addr":
		r.cfg.addr = lua.CheckString(L, val)
	case "password":
		r.cfg.password = lua.CheckString(L, val)
	case "db":
		r.cfg.db = lua.CheckInt(L, val)
	case "pool_size":
		r.cfg.poolSize = lua.CheckInt(L, val)
	case "max_conn_age":
		r.cfg.maxConnAge = lua.CheckInt(L, val)
	}
}
func LuaInjectApi(env xcall.Env) {
	env.Set("redis", lua.NewFunction(newLuaRedis))
}
