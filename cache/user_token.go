package cache

import (
	"fmt"
	"strconv"
	"strings"
)

type UserToken struct {
	Group string
	ID    int64
	Scope int
}

func (t *UserToken) Parse(s string) {
	str := strings.Split(s, "|")
	if len(str) < 3 {
		return
	}
	t.Group = str[0]
	t.ID, _ = strconv.ParseInt(str[1], 10, 64)
	scope, _ := strconv.ParseInt(str[2], 10, 64)
	t.Scope = int(scope)
}

func (t *UserToken) IDs() string {
	return strconv.FormatInt(t.ID, 10)
}

func (t *UserToken) Encode() string {
	return fmt.Sprintf("%s|%d|%d", t.Group, t.ID, t.Scope)
}

func (c *Cache) SetUserToken(token string, t UserToken) error {
	return c.setToken(TokenGroupUser, token, t.Encode(), redisUserTokenExpiration)
}

func (c *Cache) UserToken(token string) (t UserToken, _ error) {
	s, err := c.get(TokenGroupUser + token)
	if err != nil {
		return t, err
	}

	t.Parse(s)
	return t, nil
}
