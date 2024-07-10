package main

import(
	"net/http"
	"fmt"
	"restAPI/db"
	"restAPI/router"
	"restAPI/env"
)


func main(){


	envManager := env.SetupEnv("env/.env")
	strin := fmt.Sprintf("%s, %s", envManager.DbSource, envManager.DbDriver)

	fmt.Printf("DB Source: %s\n", strin)
	db := db.NewDB(envManager.DbDriver, envManager.DbSource)

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