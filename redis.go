package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
	"time"
)

func (r *Redis) Start() error {
	r.client = redis.NewClient(&redis.Options{
		DB:         r.C.db,
		Addr:       r.C.addr,
		Password:   r.C.password,
		PoolSize:   r.C.poolSize,
		MaxConnAge: time.Duration(r.C.maxConnAge) * time.Second,
	})

	r.status = lua.RUNNING
	r.uptime = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (r *Redis) Close() error {
	r.status = lua.CLOSE
	return r.client.Close()
}

func (r *Redis) State() lua.LightUserDataStatus {
	return r.status
}

func (r *Redis) Name() string {
	return r.C.name
}

func (r *Redis) Type() string {
	return "redis client"
}

func (r *Redis) Status() string {
	return fmt.Sprintf("name:%s , status:%s , uptime:%s",
		r.Name(), r.status.String(), r.uptime)
}
