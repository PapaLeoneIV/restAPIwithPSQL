package main

import(
	"net/http"
	"fmt"
	"restAPI/db"
	"restAPI/router"
)


func main(){
	
	db := db.NewDB("postgres", "host=localhost user=postgres password=ric dbname=album_db sslmode=disable")

	router := router.NewRouter(db)

	s := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Printf("Listening on port: 8080\n")
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

}