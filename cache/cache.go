package cache

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	redisKeyToken      = "token-"
	redisKeyTokenCount = "tokencount-"
)

const (
	redisTokenExpiration      = 48 * time.Hour
	redisTokenCountExpiration = 0
)

type Cache struct {
	r *redis.Client
}

func New(redisAddr string) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		r: rdb,
	}, nil
}

func (c *Cache) exist(key string) (bool, error) {
	r, err := c.r.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return r > 0, nil
}

func (c *Cache) get(key string) (string, error) {
	if ok, err := c.exist(key); err != nil || !ok {
		return "", nil
	}
	cmd := c.r.Get(key)
	return cmd.Result()
}
