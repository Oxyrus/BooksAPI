package main

import (
	"net/http"

	goji "goji.io"
	"goji.io/pat"

	mgo "gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/books"), allBooks(session))
	mux.HandleFunc(pat.Post("/books"), addBook(session))
	mux.HandleFunc(pat.Get("/books/:isbn"), bookByISBN(session))
	mux.HandleFunc(pat.Put("/books/:isbn"), updateBook(session))
	mux.HandleFunc(pat.Delete("/books/:isbn"), deleteBook(session))

	http.ListenAndServe(":8080", mux)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("store").C("books")

	index := mgo.Index{
		Key:        []string{"isbn"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)

	if err != nil {
		panic(err)
	}
}
