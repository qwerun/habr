package redis

import (
	rds "github.com/redis/go-redis/v9"
)

type RedisExplorer struct {
	RDB *rds.Client
}

func NewRedisExplorer(rdb *rds.Client) *RedisExplorer {
	return &RedisExplorer{RDB: rdb}
}
