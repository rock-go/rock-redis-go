package redis

import (
	"github.com/rock-go/rock/lua"
	"time"
	"github.com/go-redis/redis"
)

type Pipe struct {
	lua.NoReflect
	pipe redis.Pipeliner
	meta lua.UserKV
}

func newPipeline( cli *redis.Client ) *Pipe {
	p := new(Pipe)
	p.meta = lua.NewUserKV()
	p.pipe = cli.Pipeline()
	p.initMeta()
	return p
}

func (p *Pipe) initMeta() {
	p.meta.Set("hmset" , lua.NewFunction(p.LHMSet))
	p.meta.Set("hmdel" , lua.NewFunction(p.LHDelete))
	p.meta.Set("incr"  , lua.NewFunction(p.LIncr))
	p.meta.Set("expire", lua.NewFunction(p.LExpire))
	p.meta.Set("exec"  , lua.NewFunction(p.LExec))
	p.meta.Set("delete", lua.NewFunction(p.LDelete))
	p.meta.Set("close" , lua.NewFunction(p.LClose))
}

func (p *Pipe) Get(L *lua.LState , key string) lua.LValue {
	return p.meta.Get(key)
}

func (p *Pipe) LExec(L *lua.LState) int {
	var err error
	_, err = p.pipe.Exec()
	if err != nil {
		L.RaiseError("pipeline execute error: %v", err)
	}
	return 0
}


func (p *Pipe) LClose(L *lua.LState) int {
	if err := p.pipe.Close(); err != nil {
		L.RaiseError("pipeline close error: %v", err)
	}
	return 0
}

func (p *Pipe) LIncr(L *lua.LState) int {
	if _, err := p.pipe.HIncrBy(L.CheckString(1), L.CheckString(2), 1).Result(); err != nil {
		L.RaiseError("pipeline incr error: %v", err)
	}
	return 0
}

func (p *Pipe) LHMSet(L *lua.LState) int {
	field := make(map[string]interface{})
	field[L.CheckString(2)] = L.CheckInt(3)
	if _, err := p.pipe.HMSet(L.CheckString(1), field).Result(); err != nil {
		L.RaiseError("pipeline hmset error: %v", err)
	}

	return 0
}

func (p *Pipe) LExpire(L *lua.LState) int {
	var err error
	_, err = p.pipe.Expire(L.CheckString(1), time.Duration(L.CheckInt(2))*time.Second).Result()
	if err != nil {
		L.RaiseError("pipeline set expire error: %v", err)
	}
	return 0
}

func (p *Pipe) LDelete(L *lua.LState) int {
	var err error
	_, err = p.pipe.Del(L.CheckString(1)).Result()
	if err != nil {
		L.RaiseError("pipeline delete key error: %v", err)
	}

	return 0
}

func (p *Pipe) LHDelete(L *lua.LState) int {
	var err error
	_, err = p.pipe.HDel(L.CheckString(1), L.CheckString(2)).Result()
	if err != nil {
		L.RaiseError("pipeline delete hash error: %v", err)
	}

	return 0
}