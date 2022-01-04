package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Language string   `json:"language"`
	Genres   []string `json:"genres"`
}

// String returns a details of a book in a pretty formatted way
func (book *Book) String() string {

	return fmt.Sprintf("Book ID: %d, Title: %s, Author: %s, Language: %s, Genres: %v", book.ID, book.Title, book.Author, book.Language, book.Genres)
}

// ReadWriteMutex is used to allow multiple go routines to read
type Books struct {
	sync.RWMutex
	storage []Book
}

// global variable, acting as a database
var db = Books{sync.RWMutex{}, []Book{}}

// Returns all the books details
func (b *Books) getAllBooks() []Book {

	b.RLock()
	defer b.RUnlock()

	if len(b.storage) == 0 {
		return []Book{}
	}

	return b.storage

}

// Check if a book struct is empty (in the details)
func (book *Book) isEmpty() bool {

	if book.Title == "" || book.Author == "" || book.Language == "" || book.ID < 0 || book.Genres == nil {

		return true
	}
	return false

}

// Add a book to the books storage

func (b *Books) addBook(book Book) error {

	if book.isEmpty() {
		return fmt.Errorf("the provided book details is empty")
	}

	b.Lock()
	b.storage = append(b.storage, book)
	b.Unlock()

	return nil
}

// Find among the books using id
func (b *Books) findBook(id int) (Book, error) {

	if id < 0 {
		return Book{}, fmt.Errorf("the provided book id is negative")
	}

	b.RLock()
	defer b.RUnlock()
	temp := Book{}
	for _, val := range b.storage {

		if val.ID == id {
			temp = val
			break
		}
	}

	return temp, nil
}

// Home page router handler
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.IsAbs())
	fmt.Fprintf(w, "Welcome to home page.")

}

// Handles /book route
func bookHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:

		//jsonBytes, err := json.Marshal(db.getAllBooks())

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(db.getAllBooks())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	case http.MethodPost:

		var bk Book

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(reqBody, &bk)
		if err != nil {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte(err.Error()))
			return
		}

		err = db.addBook(bk)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Location", fmt.Sprintf("/book/%d", bk.ID))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(bk)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed.")

	}
}

// Handles /book/x router, where x = id of book, e.g. /book/11
func aBookHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Some internal error occured.")
		return
	}

	book, err := db.findBook(int(id))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "No such resource exists.")

	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)

}

// Routes definitions and start a http server
func handleRequests() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/book/", aBookHandler)

	log.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main() {
	// Dummy details
	books := []Book{{
		1,
		"1984",
		"George Orwell",
		"English",
		[]string{"Dystopian", "Fiction"},
	},
		{
			2,
			"Karnali Blues",
			"BuddhiSagar",
			"Nepali",
			[]string{"Novel", "Fiction"},
		},
	}

	db.storage = books

	handleRequests()

}
