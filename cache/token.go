package cache

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TokenGroupUser = "user"
)

type Token struct {
	Group string
	ID    int64
	Scope int
}

func (t *Token) Parse(s string) {
	str := strings.Split(s, "|")
	t.Group = str[0]
	t.ID, _ = strconv.ParseInt(str[1], 10, 64)
	scope, _ := strconv.ParseInt(str[2], 10, 64)
	t.Scope = int(scope)
}

func (t *Token) IDs() string {
	return strconv.FormatInt(t.ID, 10)
}

func (t *Token) Encode() string {
	return fmt.Sprintf("%s|%d|%d", t.Group, t.ID, t.Scope)
}

func (c *Cache) SetToken(token string, t Token) error {
	return c.r.Set(TokenGroupUser+token, t.Encode(), redisTokenExpiration).Err()
}

func (c *Cache) Token(token string) (t Token, _ error) {
	s, err := c.get(TokenGroupUser + token)
	if err != nil {
		return t, err
	}

	t.Parse(s)
	return t, nil
}
