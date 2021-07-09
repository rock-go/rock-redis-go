package redis

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
	"reflect"
)

var TRedis = reflect.TypeOf((*Redis)(nil)).String()

func newLuaRedis(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, TRedis)
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

func (r *Redis) NewLFunction(L *lua.LState, fn cmder) *lua.LFunction {
	return L.NewFunction(func(co *lua.LState) int {
		return fn(r.client, co)
	})
}

// LNewPipe 新建一个pipeline
func (r *Redis) LNewPipe(L *lua.LState) int {
	//pipe := Pipe{
	//	pipe: r.client.Pipeline(),
	//}

	pipe := newPipeline(r.client)
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
	case "pipeline":
		lv = L.NewFunction(r.LNewPipe)

	case "hmset":
		lv = r.NewLFunction(L, hmset)
	case "hmget":
		lv = r.NewLFunction(L, hmget)
	case "incr":
		lv = r.NewLFunction(L, incr)
	case "expire":
		lv = r.NewLFunction(L, expire)
	case "delete":
		lv = r.NewLFunction(L, del)

	default:
		L.RaiseError("%s redis %s not found", r.Name(), key)
		return lua.LNil
	}

	r.meta.Set(key, lv)
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
