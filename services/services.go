package service

import (
	"context"
	sql "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const contextKey = "id"

type Service struct {
	db *sql.DB
}

type Album struct {
	Id    string  `json:"id"`
	Title  string  `json:"title"`
	Artist string `json:"artist"`
	Price float64     `json:"price"`
}

func NewService(d *sql.DB) *Service {
	fmt.Printf("Service handler created\n")
	return &Service{db: d}
}
func (s *Service) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, album.Id, album.Title, album.Artist, album.Price)
	if err != nil {
		http.Error(w, "Failed to insert album", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(album)
}

func (s *Service) GetProduct(w http.ResponseWriter, r *http.Request) {
	/*get the context where key value/id number are stored*/
	ctx := r.Context()
	id, ok := ctx.Value(contextKey).(string)
	if id == "" || !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusBadRequest)
		return
	}
	/*query the database with the actual id*/
	var album Album
	err := s.db.QueryRowContext(ctx, "SELECT id, title, artist, price FROM albums WHERE id=$1", id).Scan(&album.Id,&album.Title, &album.Artist, &album.Price )
    if err == sql.ErrNoRows{
		http.Error(w, "Invalid or missing user ID in Database", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Invalid or missing user ID in Database", http.StatusBadRequest)
        return
    }
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(album); err != nil {
		log.Printf("Failed to encode albums to JSON: %v\n", err)
		http.Error(w, "Failed to encode albums to JSON", http.StatusInternalServerError)
	}
}

func (s *Service) ListAllProducts(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("ListProduct Handler!\n")
    fmt.Printf("Preparing the query\n")
    var ctx context.Context = r.Context()
    rows, err := s.db.QueryContext(ctx, "SELECT id, title, artist, price FROM albums")
    if err != nil {
		log.Fatal(err)
        return
    }
    defer rows.Close()
	var albums []Album
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price ); err != nil{
			log.Fatal(err)
			return 
		}
		albums = append(albums, album)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		log.Printf("Failed to encode albums to JSON: %v\n", err)
		http.Error(w, "Failed to encode albums to JSON", http.StatusInternalServerError)
	}
}

func (s *Service) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	var album Album
	if err := json.NewDecoder(body).Decode(&album); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v\n", err)
		return
	}
	dbQuery := `UPDATE albums SET title=$1, artist=$2, price=$3 WHERE id=$4`
	if _, err := s.db.Exec(dbQuery, album.Title, album.Artist, album.Price, album.Id); err != nil {
		http.Error(w, "Failed to update album", http.StatusInternalServerError)
		log.Printf("Error updating album in database: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(album)
}

func (s *Service) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := ctx.Value(contextKey).(string)
	if id == "" || !ok {
		http.Error(w, "Invalid or missing user ID in context", http.StatusBadRequest)
		return
	}
	dbQuery := `DELETE FROM albums WHERE id=$1`
	_, err := s.db.Exec(dbQuery, id)
	if err != nil {
		http.Error(w, "Failed to insert album", http.StatusInternalServerError)
		log.Printf("Error inserting album into database: %v\n", err)
		return
	}
		w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album deleted successfully"))
}
