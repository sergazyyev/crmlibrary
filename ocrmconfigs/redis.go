package ocrmconfigs

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

func (r *RedisConfig) GetRedisOptions() *redis.Options {
	address := fmt.Sprintf("%s:%d", r.Host, r.Port)
	return &redis.Options{
		Addr:     address,
		Password: r.Password,
		DB:       r.Db,
	}
}
