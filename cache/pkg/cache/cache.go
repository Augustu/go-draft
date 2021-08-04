package cache

import (
	"github.com/Augustu/go-draft/cache/pkg/config"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Client *redis.Client
}

func New(c *config.Cache) *Cache {
	cli := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Password: c.Passwd,
	})

	return &Cache{
		Client: cli,
	}
}
