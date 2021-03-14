package cache

import (
	"fmt"
	"time"
)

const (
	TokenGroupUser    = "user"
	TokenGroupMessage = "message"
)

func (c *Cache) setToken(key, token, value string, expiration time.Duration) error {
	return c.r.Set(key+token, value, expiration).Err()
}

func (c *Cache) token(key, token string) (t string, _ error) {
	s, err := c.get(key + token)
	if err != nil {
		return t, err
	}
	if s == "" {
		return t, fmt.Errorf("invalid token")
	}
	return t, nil
}
