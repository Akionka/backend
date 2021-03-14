package internal

import (
	"fmt"
	"github.com/kate-network/backend/storage"
	"github.com/labstack/echo/v4"
	"time"
)

const updateTypeWebsocket = "websocket"

type MessageService struct {
	Service
}

func (s *MessageService) Pref() string {
	return "/messages"
}

func (s *MessageService) Setup(parent Service, api *echo.Group) {
	s.Service = parent
	api.GET("/update", s.update, s.authenticated)
	api.POST("/send", s.send, s.authenticated)
}

func (s *MessageService) update(ec echo.Context) (err error) {
	c := ec.(*Context)
	updateType := c.Param("type")
	if updateType == "" || updateType == updateTypeWebsocket {
		err = s.updateWebsocket(c)
	}
	if err != nil {
		return err
	}
	return c.nocontent()
}

func (s *MessageService) updateWebsocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

type MessageSendReq struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

type MessageSendResp struct {
	ID int64 `json:"id"`
}

func (s *MessageService) send(ec echo.Context) error {
	c := ec.(*Context)
	t := c.Get("token").(string)
	req := MessageSendReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if req.Text == "" {
		return wrapError(ErrRequiredFields, "")
	}
	token, err := s.ch.UserToken(t)
	if err != nil {
		return err
	}
	// todo: make check field length
	_, err = s.db.Users.ByID(req.ID)
	if err != nil {
		return wrapNotFoundError(fmt.Errorf("user not found"))
	}
	message := &storage.Message{
		ID:          0,
		CreatedAt:   time.Time{},
		SenderID:    token.ID,
		RecipientID: req.ID,
		Message:     req.Text,
	}
	messageID, err := s.db.Messages.Create(message)
	if err != nil {
		return err
	}

	if token.ID != req.ID {
		message.SenderID, message.RecipientID = message.RecipientID, message.SenderID
		_, err = s.db.Messages.Create(message)
		if err != nil {
			return err
		}
	}

	return c.json(MessageSendResp{
		ID: messageID,
	})
}
