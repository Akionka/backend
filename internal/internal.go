package internal

import (
	"fmt"
	"github.com/kate-network/backend/cache"
	"github.com/kate-network/backend/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	ErrRequiredFields = iota
	ErrUserNotFound
	ErrUserExist
	ErrTokenEmpty
)

type Service struct {
	e  *echo.Echo
	db *storage.DB
	ch *cache.Cache
}

// CustomService allows you to create your sub-services, through which all requests will pass.
// There is a distinct feature is that your service must HAVE an rps limit.
// You can use the standard rps level, which is specified in the defaultRps variable.
// Inside the service, you can implement the API as you like.
type CustomService interface {
	// Service is one of the key things that are necessary for the specific operation of your service.
	// Properties of the main service type are inherited. This allows you to avoid duplicating code.
	// Service

	// Pref returns the unique prefix of the service, which will indicate under what conditions it should be accessed.
	// Example: you have a block for working with documents, so the prefix will be /documents.
	Pref() string

	// Setup allows you to configure the necessary APIs
	// so that the framework understands where to send what request and how to process it.
	// Do not forget that each service has an isolated group, and it should not have access to others.
	Setup(Service, *echo.Group)
}

func NewServer(db *storage.DB, ch *cache.Cache) *Service {
	return &Service{
		e:  echo.New(),
		db: db,
		ch: ch,
	}
}

func (s *Service) Init() {
	apiGroup := s.e.Group("/api")
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	s.e.HTTPErrorHandler = s.handlerError

	authService := &AuthService{}
	seedService := &SeedService{}
	userService := &UserService{}

	s.add(apiGroup, authService)
	s.add(apiGroup, seedService)
	s.add(apiGroup, userService)
}

// add connects custom services for working with the framework and the system itself
func (s Service) add(parentGroup *echo.Group, customService CustomService) {
	group := parentGroup.Group(customService.Pref())
	customService.Setup(s, group)
}

func (s *Service) Listen(address string) error {
	return s.e.Start(address)
}

func wrapError(code int, err interface{}) error {
	return echo.NewHTTPError(http.StatusConflict, map[string]interface{}{
		"code":  code,
		"error": err,
	})
}

func wrapForbiddenError(err ...error) error {
	errMsg := "unknown forbidden error"
	if err != nil {
		errMsg = err[0].Error()
	}
	return echo.NewHTTPError(http.StatusForbidden, errMsg)
}

// tokens gets token
func (s *Service) token(c echo.Context) (string, error) {
	token := c.Request().Header.Get("Token")
	if token == "" {
		return "", fmt.Errorf("token header is empty")
	}
	return token, nil
}

type errorResp struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
}

func (s *Service) handlerError(err error, c echo.Context) {
	var resp errorResp
	if he, ok := err.(*echo.HTTPError); ok {
		resp.Code = he.Code
		resp.Error = he.Message
	} else {
		resp.Code = http.StatusConflict
		resp.Error = err.Error()
	}
	logrus.Error(err)
	err = c.JSON(resp.Code, resp)
	if err != nil {
		logrus.Errorln(err)
	}
}
