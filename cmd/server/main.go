package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/lecture-backend/internal/handler"
	"github.com/go-chi/chi/v5"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/books", handler.GetBooks)
	r.Get("/books/{id}", handler.GetBookByID)
	r.Post("/books", handler.CreateBook)
	r.Put("/books/{id}", handler.UpdateBook)
	r.Delete("/books/{id}", handler.DeleteBook)

	db, err := sql.Open("sqlite3", "data/lecture.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	handler.SetDB(db)

	log.Println("server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
