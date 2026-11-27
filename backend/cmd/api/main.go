package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/victor4710hugo-rgb/RDO/backend/internal/database"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/repositories"
	"github.com/victor4710hugo-rgb/RDO/backend/internal/router"
)

// loadDotEnv lê um arquivo .env simples (CHAVE=valor) e define as variáveis
// de ambiente que ainda não estiverem definidas. Se o arquivo não existir,
// não faz nada.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			key := strings.TrimSpace(k)
			if os.Getenv(key) == "" {
				os.Setenv(key, strings.Trim(strings.TrimSpace(v), `"`))
			}
		}
	}
}

func main() {
	// Precisa vir antes de router.New, que lê JWT_SECRET via AuthService.
	loadDotEnv(".env")

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

	mux, err := router.New(repo)
	if err != nil {
		log.Fatalf("erro ao inicializar router: %v", err)
	}

	log.Println("API RDO iniciada em http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}