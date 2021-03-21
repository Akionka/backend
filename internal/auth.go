package internal

import (
	"fmt"
	"github.com/kate-network/backend/storage"
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

	api.POST("/login", s.login)
	api.POST("/signup", s.signup)
	api.POST("/message", s.messages, s.authenticated)
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

func (s *AuthService) login(ec echo.Context) error {
	c := ec.(*Context)
	var req userAuthReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Login == "" || req.Password == "" {
		return wrapError(ErrRequiredFields, "required fields are not filled in")
	}

	user, err := s.db.Users.ByLogin(req.Login)
	if err != nil {
		return wrapNotFoundError(fmt.Errorf("user does not exist"))
	}

	if user.Password != req.Password {
		return wrapNotFoundError(fmt.Errorf("user does not exist"))
	}

	scope := 32768 // 2^5
	token := uuid.New().String()
	t := cache.UserToken{
		Group: cache.TokenGroupUser,
		ID:    user.ID,
		Scope: scope,
	}

	if err := s.ch.SetUserToken(token, t); err != nil {
		return err
	}

	if !req.ServerCookie {
		c.SetCookie(&http.Cookie{
			Name:     "token",
			Value:    token,
			SameSite: http.SameSiteNoneMode,
		})
	}

	return c.json(UserAuthResp{
		ID:    user.ID,
		Token: token,
		Scope: scope,
	})
}

type AuthMessagesResp struct {
	Token string `json:"token"`
}

func (s *AuthService) messages(ec echo.Context) error {
	c := ec.(*Context)
	userToken := c.Get("token").(string)
	token := uuid.New().String()
	if err := s.ch.SetMessageToken(token, userToken); err != nil {
		return err
	}
	return c.json(AuthMessagesResp{
		Token: token,
	})
}

type AuthSignupReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *AuthService) signup(ec echo.Context) error {
	// todo: make validate fields
	// todo: check field length
	c := ec.(*Context)
	var req AuthSignupReq
	if err := c.Bind(&req); err != nil {
		return wrapError(ErrRequiredFields, "login or password is empty")
	}
	if req.Login == "" || req.Password == "" {
		return wrapError(ErrRequiredFields, "login or password is empty")
	}
	if err := s.db.Users.Create(&storage.User{
		Login:    req.Login,
		Password: req.Password,
	}); err != nil {
		return wrapError(ErrUserExist, "user exists")
	}

	return c.nocontent()
}
