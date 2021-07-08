package redis

import (
	"github.com/rock-go/rock/lua"
	"time"
)

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

func (p *Pipe) Index(L *lua.LState, key string) lua.LValue {
	if key == "hmset" {
		return lua.NewFunction(p.LHMSet)
	}
	if key == "hmdel" {
		return lua.NewFunction(p.LHDelete)
	}
	if key == "incr" {
		return lua.NewFunction(p.LIncr)
	}
	if key == "expire" {
		return lua.NewFunction(p.LExpire)
	}
	if key == "exec" {
		return lua.NewFunction(p.LExec)
	}
	if key == "delete" {
		return lua.NewFunction(p.LDelete)
	}
	if key == "close" {
		return lua.NewFunction(p.LClose)
	}

	return lua.LNil
}
