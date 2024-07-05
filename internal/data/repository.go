package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Repositories struct {
	BookRepo interface {
		Insert(book *Book) error
		Get(id int64) (*Book, error)
		Update(book *Book) error
		Delete(id int64) error
	}
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		BookRepo: BookRepository{DB: db},
	}
}

func NewMockRepositories(db *sql.DB) Repositories {
	return Repositories{
		BookRepo: MockBookRepository{DB: db},
	}
}
