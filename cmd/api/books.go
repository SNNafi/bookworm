package main

import (
	"bookworm.snnafi.dev/internal/data"
	"bookworm.snnafi.dev/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name       string          `json:"name"`
		Author     string          `json:"author"`
		Publisher  string          `json:"publisher"`
		Image      string          `json:"image"`
		CoverImage string          `json:"cover_image"`
		Type       []data.BookType `json:"type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Name:       input.Name,
		Author:     input.Author,
		Publisher:  input.Publisher,
		Image:      input.Image,
		CoverImage: input.CoverImage,
		Type:       input.Type,
	}

	v := validator.New()

	if book.ValidateBook(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	book := data.Book{
		ID:         id,
		Name:       "সীরাহ ১ম খণ্ড",
		Author:     "রেইনড্রপস মিডিয়া",
		Publisher:  "রেইনড্রপস মিডিয়া",
		Image:      "seerahlastpart.webp",
		CoverImage: "seerah_cover.webp",
		Type: []data.BookType{
			data.Islamic,
		},
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
