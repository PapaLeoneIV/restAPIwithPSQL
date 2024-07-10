package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

const contextKey = "id"

type Service struct {
	db *sql.DB
}

type Album struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func NewService(d *sql.DB) *Service {
	fmt.Printf("Service handler created\n")
	return &Service{db: d}
}

func (s *Service) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateProduct Handler called")
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("Error decoding request body: %v\n", err)
		return
	}
	fmt.Printf("Decoded album: %+v\n", album)

	query := `INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, album.Id, album.Title, album.Artist, album.Price)
	if err != nil {
		http.Error(w, "Failed to insert album", http.StatusInternalServerError)
		fmt.Printf("Error inserting album into database: %v\n", err)
		return
	}
	fmt.Println("Album inserted successfully")

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Service) GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetProduct Handler called")
	ctx := r.Context()
	id, ok := ctx.Value(contextKey).(string)
	if id == "" || !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusBadRequest)
		fmt.Println("Invalid or missing user ID in context")
		return
	}
	fmt.Printf("Context ID: %s\n", id)

	var album Album
	err := s.db.QueryRowContext(ctx, "SELECT id, title, artist, price FROM albums WHERE id=$1", id).Scan(&album.Id, &album.Title, &album.Artist, &album.Price)
	if err == sql.ErrNoRows {
		http.Error(w, "Album not found", http.StatusNotFound)
		fmt.Printf("Album not found for ID: %s\n", id)
		return
	} else if err != nil {
		http.Error(w, "Error fetching album", http.StatusInternalServerError)
		fmt.Printf("Error fetching album: %v\n", err)
		return
	}
	fmt.Printf("Fetched album: %+v\n", album)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(album); err != nil {
		fmt.Printf("Failed to encode album to JSON: %v\n", err)
		http.Error(w, "Failed to encode album to JSON", http.StatusInternalServerError)
	}
}

func (s *Service) ListAllProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListAllProducts Handler called")
	ctx := r.Context()

	rows, err := s.db.QueryContext(ctx, "SELECT id, title, artist, price FROM albums")
	if err != nil {
		http.Error(w, "Failed to fetch albums", http.StatusInternalServerError)
		fmt.Printf("Error fetching albums: %v\n", err)
		return
	}
	defer rows.Close()
	fmt.Println("Albums fetched from database")

	var albums []Album
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			http.Error(w, "Error scanning album row", http.StatusInternalServerError)
			fmt.Printf("Error scanning album row: %v\n", err)
			return
		}
		albums = append(albums, album)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating album rows", http.StatusInternalServerError)
		fmt.Printf("Error iterating album rows: %v\n", err)
		return
	}
	fmt.Printf("Fetched albums: %+v\n", albums)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		fmt.Printf("Failed to encode albums to JSON: %v\n", err)
		http.Error(w, "Failed to encode albums to JSON", http.StatusInternalServerError)
	}
}

func (s *Service) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateProduct Handler called")
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Printf("Error decoding request body: %v\n", err)
		return
	}
	fmt.Printf("Decoded album for update: %+v\n", album)

	dbQuery := `UPDATE albums SET title=$1, artist=$2, price=$3 WHERE id=$4`
	effected, err := s.db.Exec(dbQuery, album.Title, album.Artist, album.Price, album.Id)
	if err != nil {
		http.Error(w, "Failed to update album", http.StatusInternalServerError)
		fmt.Printf("Error updating album in database: %v\n", err)
		return
	}
	if rows, _ := effected.RowsAffected(); rows == 0 {
		http.Error(w, "Album not found", http.StatusNotFound)
		fmt.Printf("Album not found for ID: %s\n", album.Id)
		return
	}
	
	fmt.Println("Album updated successfully")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Service) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteProduct Handler called")
	ctx := r.Context()
	id, ok := ctx.Value(contextKey).(string)
	if id == "" || !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusBadRequest)
		fmt.Println("Invalid or missing user ID in context")
		return
	}
	fmt.Printf("Context ID for delete: %s\n", id)

	dbQuery := `DELETE FROM albums WHERE id=$1`
	effected, err := s.db.Exec(dbQuery, id)
	if rows, _ := effected.RowsAffected(); rows == 0 {
		http.Error(w, "Album not found", http.StatusNotFound)
		fmt.Printf("Album not found for ID: %s\n", id)
		return
	}
	if err != nil {
		http.Error(w, "Failed to delete album", http.StatusInternalServerError)
		fmt.Printf("Error deleting album from database: %v\n", err)
		return
	}
	fmt.Println("Album deleted successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album deleted successfully"))
}
