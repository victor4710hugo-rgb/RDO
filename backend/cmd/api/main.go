package main

import (
	"log"
	"net/http"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/database"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/router"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatal(err)
	}

	repo := repositories.New(db)
	mux := router.New(repo)

	log.Println("API RDO iniciada em http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}