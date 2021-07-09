package redis

import (
	"fmt"
	"context"
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
	"time"
)

type Redis struct {
	lua.Super
	cfg *config
	client *redis.Client
	ctx context.Context
	meta lua.UserKV
}

func newRedis(cfg *config) *Redis {
	r := &Redis{cfg:cfg}
	r.S = lua.INIT
	r.T = TRedis
	return r
}

func (r *Redis) Start() error {
	r.client = redis.NewClient(r.cfg.Options())
	r.S = lua.RUNNING
	r.U = time.Now()
	return nil
}

func (r *Redis) Close() error {
	r.S = lua.CLOSE
	return r.client.Close()
}

func (r *Redis) Name() string {
	return r.cfg.name
}

func (r *Redis) Status() string {
	return fmt.Sprintf("name:%s , status:%s , uptime:%s",
		r.Name(), r.S.String(), r.U)
}
