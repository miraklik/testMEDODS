package main

import (
	"log"
	"net/http"

	"main/internal/handlers"
	"main/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	store, err := storage.NewStorage("postgres://username:password@localhost/dbname")
	if err != nil {
		log.Fatal(err)
	}

	authHandler := handlers.NewAuthHandler(store)

	r := mux.NewRouter()
	r.HandleFunc("/token", authHandler.TokenHandler).Methods("GET")
	r.HandleFunc("/refresh", authHandler.RefreshHandler).Methods("POST")

	log.Println("Auth service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
