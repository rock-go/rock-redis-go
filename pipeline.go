package redis

import (
	"github.com/rock-go/rock/lua"
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
	p.meta.Set("close" , lua.NewFunction(p.LClose))
	p.meta.Set("exec"  , lua.NewFunction(p.LExec))

	p.meta.Set("hmset" , p.NewFunction(hmget))
	p.meta.Set("hmdel" , p.NewFunction(hdel))
	p.meta.Set("incr"  , p.NewFunction(incr))
	p.meta.Set("expire", p.NewFunction(expire))
	p.meta.Set("delete", p.NewFunction(del))
}

func (p *Pipe) NewFunction(fn cmder) *lua.LFunction {
	return lua.NewFunction(func(co *lua.LState) int {
		return fn(p.pipe , co)
	})
	
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