package internal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kate-network/backend/cache"
	"github.com/labstack/echo/v4"
)

type AuthService struct {
	Service
}

func (s *AuthService) Pref() string {
	return "/auth"
}

func (s *AuthService) Setup(parent Service, api *echo.Group) {
	s.Service = parent

	api.POST("/user", s.user)
}

type userAuthReq struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	Scope        int    `json:"scope"`
	ServerCookie bool   `json:"server_cookie"`
}

type UserAuthResp struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
	Scope int    `json:"scope"`
}

func (s *AuthService) user(c echo.Context) error {
	var req userAuthReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Login == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "required fields are not filled in")
	}

	user, err := s.db.Users.ByLogin(req.Login)
	if err != nil {
		return wrapError(ErrUserNotFound, "user is not exists")
	}

	if user.Password != req.Password {
		return wrapError(ErrUserNotFound, "user is not exists")
	}

	scope := 32768 // 2^5
	token := uuid.New().String()
	t := cache.Token{
		Group: cache.TokenGroupUser,
		ID:    user.ID,
		Scope: scope,
	}

	if err := s.ch.SetToken(token, t); err != nil {
		return err
	}

	if !req.ServerCookie {
		c.SetCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})
	}

	return c.JSON(http.StatusOK, UserAuthResp{
		ID:    user.ID,
		Token: token,
		Scope: scope,
	})
}
