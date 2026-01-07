package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type createBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type updateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

/*
GET /books
*/
func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, title, author, created_at, updated_at
		FROM books
		ORDER BY id ASC
	`)
	if err != nil {
		http.Error(w, "failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		var b Book
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			http.Error(w, "failed to scan book", http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

/*
GET /books/{id}
*/
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var b Book
	err = db.QueryRow(`
		SELECT id, title, author, created_at, updated_at
		FROM books
		WHERE id = ?
	`, id).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

/*
POST /books
*/
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var req createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)

	result, err := db.Exec(`
		INSERT INTO books (title, author, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`, req.Title, req.Author, now, now)
	if err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	book := Book{
		ID:        int(id),
		Title:     req.Title,
		Author:    req.Author,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

/*
PUT /books/{id}
*/
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)

	result, err := db.Exec(`
		UPDATE books
		SET title = ?, author = ?, updated_at = ?
		WHERE id = ?
	`, req.Title, req.Author, now, id)
	if err != nil {
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
DELETE /books/{id}
*/
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM books WHERE id = ?
	`, id)
	if err != nil {
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
