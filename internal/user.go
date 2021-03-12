package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserService struct {
	Service
}

func (s *UserService) Pref() string {
	return "/users"
}

func (s *UserService) Setup(parent Service, api *echo.Group) {
	s.Service = parent

	api.GET("/me", s.me)
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Scope    int    `json:"scope"`
}

type UserMeResp struct {
	User
}

func (s *UserService) me(c echo.Context) error {
	token, err := s.token(c)
	if err != nil {
		return wrapForbiddenError(err)
	}
	t, err := s.ch.Token(token)
	if err != nil {
		return wrapForbiddenError(err)
	}
	user, err := s.db.Users.ByID(t.ID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, UserMeResp{
		User{
			ID:       t.ID,
			Username: user.Username,
			Scope:    t.Scope,
		},
	})
}
