package storage

import (
	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(*User) error
		ByLogin(string) (User, error)
		ByID(int64) (User, error)
	}

	Users struct {
		*sqlx.DB
	}

	User struct {
		ID       int64  `sq:"id"`
		Login    string `sq:"login"`
		Password string `sq:"password"`
		Username string `sq:"username"`
	}
)

func (db *Users) Create(user *User) error {
	const q = "INSERT INTO users (login, password, username) VALUES (?, ?, ?)"
	_, err := db.Exec(q, user.Login, user.Password, user.Username)
	return err
}

func (db *Users) ByLogin(login string) (u User, _ error) {
	const q = "SELECT * FROM users WHERE login = ?"
	return u, db.Get(&u, q, login)
}

func (db *Users) ByID(id int64) (u User, _ error) {
	const q = "SELECT * FROM users WHERE id = ?"
	return u, db.Get(&u, q, id)
}
