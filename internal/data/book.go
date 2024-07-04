package data

import (
	"bookworm.snnafi.dev/internal/validator"
	"errors"
	"time"
)

var ErrInvalidBookTypeFormat = errors.New("invalid book type format")

type BookType int

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
	v.Check(book.Publisher != "", "author", "must be provided")
	v.Check(book.Image != "", "image", "must be provided")
	v.Check(book.Type != nil, "type", "must be provided")
	v.Check(len(book.Type) >= 1, "type", "must contain at least 1 type")
	v.Check(len(book.Type) <= 3, "type", "must not contain more than 3 types")
	v.Check(validator.Unique(book.Type), "type", "must not contain duplicate values")
}
