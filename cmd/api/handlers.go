package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"provider-api/internal/data"
	"strconv"
	"time"
	"BE-PROVIDER/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     "1.0.0",
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	js = append(js, '\n')
	w.Write(js)

}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "----->\n", r.Method)
	if r.Method == http.MethodGet {
		books, err := app.models.Books.GetAll()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
		fmt.Fprintln(w, "Display a list of books on reading list")
		
		return
	}
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "-----\n")
		// The pieces we need to extract from the request body.
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float64  `json:"rating"`
		}
		// Read Body
		// body, err := ioutil.ReadAll(r.Body)
		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		


		// &input is a pointer to the input struct
		// err = json.Unmarshal(body, &input)
		// if err != nil {
		// 	http.Error(w, "Bad Request", http.StatusBadRequest)
		// 	return
		// }
		fmt.Fprintf(w, " %v\n", input)

	}
}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		fmt.Fprintln(w, "update a book")
		app.updateBook(w, r)
	case http.MethodGet:
		fmt.Fprintln(w, "GET a book")
		app.getBook(w, r)
	case http.MethodPost:
		app.postBook(w, r)
	case http.MethodDelete:
		fmt.Fprintln(w, "Delete a book")
		app.deleteBook(w, r)
	default:
		fmt.Fprintf(w, "came here")
		app.updateBook(w, r)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	book := data.Book{
		ID:        idInt,
		Title:     "The Great Gatsby",
		CreatedAt: time.Now(),
		Published: time.Now(),
		Pages:     218,
		Genres:    []string{"Fiction", "Tragedy"},
		Rating:    4.5,
		Version:   1,
	}

	erro := app.writeJSON(w, http.StatusOK, envelope{"book": book})

	// js, err := json.MarshalIndent(book, "", "\t")
	if erro != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)
	// return
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	// we use pointer to update the existing ones
	var input struct {
		Title     *string   `json:"title"`
		Published *int      `json:"published"`
		Pages     *int      `json:"pages"`
		Genres    *[]string `json:"genres"`
		Rating    *float64  `json:"rating"`
	}
	book := data.Book{
		ID:        idInt,
		Title:     "The Great Gatsby",
		CreatedAt: time.Now(),
		Published: time.Now(),
		Pages:     218,
		Genres:    []string{"Fiction", "Tragedy"},
		Rating:    4.5,
		Version:   1,
	}
	err = app.readJSON(w, r, &input)

	// body, err := ioutil.ReadAll(r.Body)

	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }
	// //read request body and unmarshal to put it into input
	// err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Published != nil {
		book.Published = time.Now()
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if len(*input.Genres) > 0 {
		book.Genres = *input.Genres
	}

	if input.Rating != nil {
		book.Rating = *input.Rating
	}

	fmt.Fprintf(w, "%v\n", book)
}

func (app *application) postBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "post book with ID %d\n", idInt)
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/books/"):]
	idInt, err := strconv.ParseInt(string(id), 10, 64)
	if err != nil {
		http.Error(w, "Invalid book ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "delete book with ID %d\n", idInt)
}
