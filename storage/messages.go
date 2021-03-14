package storage

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	MessagesStorage interface {
		Create(*Message) (int64, error)
	}

	Messages struct {
		*sqlx.DB
	}

	Message struct {
		ID          int64
		CreatedAt   time.Time
		SenderID    int64
		RecipientID int64
		Message     string
	}
)

func (db *Messages) Create(message *Message) (int64, error) {
	const q = "INSERT INTO messages (sender_id, recipient_id, message) VALUES (?, ?, ?)"
	r, err := db.Exec(q, message.SenderID, message.RecipientID, message.Message)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}
