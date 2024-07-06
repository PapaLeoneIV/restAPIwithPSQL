package service

import (
	sql "database/sql"
	"encoding/json"
	"net/http"
)

type Service struct {
	db *sql.DB
}

type Album struct {
	Id    string  `json:"id"`
	Title  string  `json:"title"`
	Artist string `json:"artist"`
	Price int     `json:"price"`
}

func NewService(d *sql.DB) *Service {
	return &Service{db: d}
}

func (s *Service) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var album Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request payload"))
		return
	}
	var dbName string
	if err := s.db.QueryRow("SELECT current_database()").Scan(&dbName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	dbQuery := "INSERT INTO " +  dbName + " (id, title, artist, price) AS VALUES $1 $2 $3 $4"
	_, err := s.db.Exec(dbQuery, album.Id, album.Title, album.Artist, album.Price)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to create product"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("product created"))
}

func (s *Service) GetProducts(w http.ResponseWriter, r *http.Request) {



}

func (s *Service) UpdateProduct(w http.ResponseWriter, r *http.Request) {
}

func (s *Service) DeleteProduct(w http.ResponseWriter, r *http.Request) {
}
