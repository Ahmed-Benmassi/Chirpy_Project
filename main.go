package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)



type apiConfig struct {        //storing config setting fo interaction with api
	fileserverHits atomic.Int32     //atomic.Int32   provide a thread safe operation on 32bit int so there is no changes that can occure like go routines and so on
	db             *database.Queries
	platform        string
	jwtSecret      string
}


type User struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
}

type Chirp struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string    `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}


func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}


	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)


	platform := os.Getenv("PLATFORM")





	mux := http.NewServeMux()                  //// NewServeMux allocates and returns a new [ServeMux]
	
	apicfg:=apiConfig{                             // create an instance of apiConfig to store the hits and use it in the handlers
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:  platform,
		jwtSecret :     jwtSecret,
	}                                           
	fs:=http.FileServer(http.Dir(filepathRoot))    // FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
	
    
	mux.Handle("/app/", http.StripPrefix("/app",apicfg.middlewareFileserverHits(fs)))  // handel the resuest and uses the stripprefixto rmove "/app" and add it to middlewar to wrap the file server and add the hit for every request to fs
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics",apicfg.getHits)
	mux.HandleFunc("POST /admin/reset",apicfg.resetHits)
	mux.HandleFunc("GET /api/chirps", apicfg.listChirpsHandler)   // new GET endpoint
	mux.HandleFunc("GET /admin/reset", apicfg.adminResethandler)
	mux.HandleFunc("POST /api/users", apicfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apicfg.createChirpHandler)
    mux.HandleFunc("GET /api/chirps/{chirpID}",apicfg.getsinglechirphandeler)
	mux.HandleFunc("POST /api/login",apicfg.loginHandler)
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

