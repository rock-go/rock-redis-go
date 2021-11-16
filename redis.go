package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
)

type Redis struct {
	lua.Super
	cfg    *config
	client *redis.Client
	ctx    context.Context
	meta   lua.UserKV
}

func newRedis(cfg *config) *Redis {
	r := &Redis{cfg: cfg}
	r.V(lua.INIT , redisTypeOf)
	r.meta = lua.NewUserKV()
	return r
}

func (r *Redis) Start() error {
	r.client = redis.NewClient(r.cfg.Options())
	r.meta = lua.NewUserKV()
	logger.Infof("%s redis start successfully", r.cfg.name)
	return nil
}

func (r *Redis) Close() error {
	r.V(lua.CLOSE)
	return r.client.Close()
}

func (r *Redis) Name() string {
	return r.cfg.name
}