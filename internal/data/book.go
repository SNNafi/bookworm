package data

import (
	"bookworm.snnafi.dev/internal/validator"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidBookTypeFormat = errors.New("invalid book type format")

type BookType int32

func (t BookType) MarshalJSON() ([]byte, error) {
	switch t {
	case 1:
		return []byte(`"Islamic"`), nil
	case 2:
		return []byte(`"Comparative Religion"`), nil
	}
	return []byte(`""`), nil
}

func (t *BookType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"Islamic"`:
		*t = 1
	case `"Comparative Religion"`:
		*t = 2
	default:
		return ErrInvalidBookTypeFormat
	}
	return nil
}

func (t *BookType) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("incompatible type")
	}

	source := string(bytes)
	source = strings.TrimSpace(source)
	v, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		return err
	}
	*t = BookType(v)
	return nil
}

//func (t BookType) Value() (driver.Value, error) {
//	return int32(t), nil
//}

const (
	Islamic BookType = iota + 1
	ComparativeReligion
)

type Book struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Author     string     `json:"author"`
	Publisher  string     `json:"publisher"`
	Image      string     `json:"image"`
	CoverImage string     `json:"cover_image,omitempty"`
	Type       []BookType `json:"type,omitempty"`
	CreatedAt  time.Time  `json:"-"`
}

func (book *Book) ValidateBook(v *validator.Validator) {
	v.Check(book.Name != "", "name", "must be provided")
	v.Check(book.Author != "", "author", "must be provided")
	v.Check(book.Publisher != "", "publisher", "must be provided")
	v.Check(book.Image != "", "image", "must be provided")
	v.Check(book.Type != nil, "type", "must be provided")
	v.Check(len(book.Type) >= 1, "type", "must contain at least 1 type")
	v.Check(len(book.Type) <= 3, "type", "must not contain more than 3 types")
	v.Check(validator.Unique(book.Type), "type", "must not contain duplicate values")
}

type BookRepository struct {
	DB *sql.DB
}

func (repo BookRepository) Insert(book *Book) error {

	query := `INSERT INTO books (name, author, publisher, image, cover_image, types) VALUES
              ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	args := []any{book.Name, book.Author, book.Publisher, book.Image, book.CoverImage, pq.Array(book.Type)}

	return repo.DB.QueryRow(query, args...).Scan(&book.ID, &book.CreatedAt)
}

func (repo BookRepository) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	stmt := `SELECT id, created_at, name, author, publisher, image, cover_image, types
    FROM books WHERE id = $1`

	var book Book

	err := repo.DB.QueryRow(stmt, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Name,
		&book.Author,
		&book.Publisher,
		&book.Image,
		&book.CoverImage,
		pq.Array(&book.Type),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (repo BookRepository) Update(book *Book) error {

	stmt := `UPDATE books SET name = $1, author = $2, publisher = $3, image = $4, cover_image = $5, types = $6 WHERE id = $7`
	args := []any{book.Name, book.Author, book.Publisher, book.Image, book.CoverImage, pq.Array(book.Type), book.ID}
	_, err := repo.DB.Exec(stmt, args...)
	return err
}

func (repo BookRepository) Delete(id int64) error {

	stmt := `DELETE FROM books WHERE id = $1`

	result, err := repo.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type MockBookRepository struct {
	DB *sql.DB
}

func (repo MockBookRepository) Insert(book *Book) error {
	return nil
}

func (repo MockBookRepository) Get(id int64) (*Book, error) {
	return nil, nil
}

func (repo MockBookRepository) Update(book *Book) error {
	return nil
}

func (repo MockBookRepository) Delete(id int64) error {
	return nil
}
