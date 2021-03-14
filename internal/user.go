package internal

import (
	"fmt"
	"github.com/kate-network/backend/storage"
	"github.com/labstack/echo/v4"
	"strconv"
)

type UserService struct {
	Service
}

func (s *UserService) Pref() string {
	return "/users"
}

func (s *UserService) Setup(parent Service, api *echo.Group) {
	s.Service = parent

	api.GET("/me", s.me, s.authenticated)
	api.GET("/find/:param", s.find, s.authenticated)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Scope    *int   `json:"scope,omitempty"`
}

func newUser(user storage.User) User {
	return User{
		ID:       user.ID,
		Username: user.Username,
		Scope:    nil,
	}
}

type UserMeResp struct {
	User
}

func (s *UserService) me(ec echo.Context) error {
	c := ec.(*Context)
	token := c.Get("token").(string)
	t, err := s.ch.UserToken(token)
	if err != nil {
		return wrapForbiddenError(err)
	}
	user, err := s.db.Users.ByID(t.ID)
	if err != nil {
		return wrapForbiddenError(err)
	}
	return c.json(UserMeResp{
		User{
			ID:       t.ID,
			Username: user.Username,
			Scope:    &t.Scope,
		},
	})
}

type UsersFindResp struct {
	User
}

func (s *UserService) find(ec echo.Context) (err error) {
	c := ec.(*Context)
	var user storage.User
	param := c.Param("param")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		user, err = s.db.Users.ByUsername(param)
	} else {
		user, err = s.db.Users.ByID(id)
	}
	if err != nil {
		return wrapNotFoundError(fmt.Errorf("user not found"))
	}
	u := newUser(user)
	return c.json(u)
}
