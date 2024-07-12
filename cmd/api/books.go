package main

import (
	"bookworm.snnafi.dev/internal/data"
	"bookworm.snnafi.dev/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string
		BookType []data.BookType
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")

	input.BookType = Map(app.readCSV(qs, "type", []string{}), data.NewBookTypeFromString)
	input.BookType = Filter(input.BookType, func(bookType data.BookType) bool {
		return bookType != 0
	})

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)

	input.SortBy = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"id", "-id", "name", "-name", "author", "-author", "publisher", "-publisher"}

	if input.ValidateFilters(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := app.repos.BookRepo.GetAll(input.Name, input.BookType, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

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

	err = app.repos.BookRepo.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.repos.BookRepo.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.repos.BookRepo.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Name       *string         `json:"name"`
		Author     *string         `json:"author"`
		Publisher  *string         `json:"publisher"`
		Image      *string         `json:"image"`
		CoverImage *string         `json:"cover_image"`
		Type       []data.BookType `json:"type"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		book.Name = *input.Name
	}

	if input.Author != nil {
		book.Author = *input.Author
	}

	if input.Publisher != nil {
		book.Publisher = *input.Publisher
	}

	if input.Image != nil {
		book.Image = *input.Image
	}

	if input.CoverImage != nil {
		book.CoverImage = *input.CoverImage
	}

	if input.Type != nil {
		book.Type = input.Type
	}

	v := validator.New()

	if book.ValidateBook(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.repos.BookRepo.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.repos.BookRepo.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
