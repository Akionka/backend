package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kate-network/backend/storage"
	"github.com/labstack/echo/v4"
)

type SeedService struct {
	Service
}

func (s *SeedService) Pref() string {
	return "/seed"
}

func (s *SeedService) Setup(parent Service, api *echo.Group) {
	s.Service = parent

	api.GET("/user", s.user)

}

func (s *SeedService) user(c echo.Context) error {
	if strings.Split(c.Request().RemoteAddr, ":")[0] != "127.0.0.1" {
		return echo.NewHTTPError(http.StatusForbidden, "")
	}

	if err := s.db.Users.Create(&storage.User{
		Login:    "admin",
		Password: "admin",
		Username: "admin",
	}); err != nil {
		fmt.Println(err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
