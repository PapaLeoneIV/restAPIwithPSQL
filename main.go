package main

import(
	"net/http"
	"fmt"
	"restAPI/db"
	"restAPI/router"
)


func main(){
	
	db := db.NewDB("postgres", "user=postgres dbname=album_db sslmode=disable")

	router := router.NewRouter(db)

	s := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Printf("Listening on port: 8080")
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

}