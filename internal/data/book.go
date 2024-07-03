package data

import (
	"time"
)

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

const (
	Islamic BookType = iota + 1
	ComparativeReligion
)

type Book struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Writer     string     `json:"writer"`
	Publisher  string     `json:"publisher"`
	Image      string     `json:"image"`
	CoverImage string     `json:"cover_image,omitempty"`
	Type       []BookType `json:"type,omitempty"`
	CreatedAt  time.Time  `json:"-"`
}
