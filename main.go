package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/englandrecoil/go-avito-shop/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db     *database.Queries
	secret string
	conn   *sql.DB
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("`.env` file must be created")
	}
	const port = "8080"

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:     dbQueries,
		secret: secret,
		conn:   dbConn,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/reg", apiCfg.handlerRegister)
	mux.HandleFunc("POST /api/auth", apiCfg.handlerAuth)
	mux.HandleFunc("POST /api/sendCoin", apiCfg.handlerSendCoins)

	mux.HandleFunc("GET /api/buy/{item}", apiCfg.handlerBuyItem)
	mux.HandleFunc("GET /api/info", apiCfg.handlerInfo)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
