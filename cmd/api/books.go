package main

import (
	"bookworm.snnafi.dev/internal/data"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create a new book")
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
		Writer:     "রেইনড্রপস মিডিয়া",
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
