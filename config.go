package redis

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/rock-go/rock/lua"
)

type Config struct {
	name       string
	addr       string
	password   string
	db         int
	poolSize   int
	maxConnAge int
}

type Redis struct {
	lua.Super
	C Config

	client *redis.Client

	status lua.LightUserDataStatus
	uptime string

	ctx context.Context
}

// Pipe 批量提交命令，当需要从多个维度去分析一条消息时，使用pipeline
type Pipe struct {
	lua.Super
	pipe redis.Pipeliner
}

func (p *Pipe) Type() string {
	return "pipeline"
}
