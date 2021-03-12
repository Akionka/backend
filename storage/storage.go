package storage

import (
	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func init() {
	// structs is used with squirrel (sq)
	structs.DefaultTagName = "sq"
}

type DB struct {
	*sqlx.DB
	Users UsersStorage
}

func Open(url string) (*DB, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return &DB{
		DB:    db,
		Users: &Users{DB: db},
	}, nil
}
