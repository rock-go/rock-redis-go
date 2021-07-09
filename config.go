package redis

import (
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
	"time"
)

type config struct {
	name       string
	addr       string
	password   string
	db         int
	poolSize   int
	maxConnAge int
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}
	tab.ForEach(func(key lua.LValue, val lua.LValue) {
		switch key.String() {
		case "name":
			cfg.name = utils.CheckProcName(val , L)

		case "addr":
			cfg.addr = utils.CheckSockets(val , L)

		case "password":
			cfg.password = utils.LValueToStr(val , "")

		case "db":
			cfg.db = utils.LValueToInt(val , 0)

		case "pool_size":
			cfg.poolSize = utils.LValueToInt(val , 10)

		case "max_conn_age":
			cfg.maxConnAge = utils.LValueToInt(val , 10)

		default:
			L.RaiseError("not found %s key" , key.String())
		}
	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v" , e)
		return nil
	}

	return cfg
}

func (cfg *config) Options() *redis.Options {
	return &redis.Options{
		DB:         cfg.db,
		Addr:       cfg.addr,
		Password:   cfg.password,
		PoolSize:   cfg.poolSize,
		MaxConnAge: time.Duration(cfg.maxConnAge) * time.Second,
	}
}

func (cfg *config) verify() error {
	return nil
}

// Pipe 批量提交命令，当需要从多个维度去分析一条消息时，使用pipeline
//type Pipe struct {
//	pipe redis.Pipeliner
//}

