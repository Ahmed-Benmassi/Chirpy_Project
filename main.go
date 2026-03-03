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


type User struct {                                   // User struct represents a user in the system. It contains fields for the user's ID, creation and update timestamps, and email address. The ID is a UUID, which is a universally unique identifier that can be used to uniquely identify a user across different systems. The CreatedAt and UpdatedAt fields are of type time.Time, which is a built-in Go type for representing date and time. The Email field is a string that stores the user's email address.
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Email     string    `json:"email"`
}

type Chirp struct {                                     // Chirp struct represents a chirp in the system. It contains fields for the chirp's ID, creation and update timestamps, body text, and the ID of the user who created the chirp. 
                                                        //The ID is a UUID, which is a universally unique identifier that can be used to uniquely identify a chirp across different systems. The CreatedAt and UpdatedAt fields are of type time.Time, which is a built-in Go type for representing date and time. The Body field is a string that stores the text content of the chirp. The UserID field is a UUID that references the ID of the user who created the chirp.
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string    `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}


func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()                                          // Load environment variables from .env file, if it exists. This is useful for local development and testing.
	dbURL := os.Getenv("DB_URL")                             // Get the database URL from the environment variable DB_URL. This should be set to a valid PostgreSQL connection string. If it is not set, the program will log a fatal error and exit.
	if dbURL == "" { 
		log.Fatal("DB_URL must be set")
	}   
	jwtSecret := os.Getenv("JWT_SECRET")                                         // Get the platform from the environment variable PLATFORM. This can be used to differentiate between different deployment environments (e.g., development, staging, production). If it is not set, it will default to an empty string.
	if jwtSecret == "" {                                                   // Get the JWT secret from the environment variable JWT_SECRET. This should be set to a secure random string that is used to sign and verify JSON Web Tokens (JWTs). If it is not set, the program will log a fatal error and exit.
		log.Fatal("JWT_SECRET environment variable is not set")
	}


	dbConn, err := sql.Open("postgres", dbURL)                         // sql.Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information. In this case, it is opening a PostgreSQL database using the connection string provided in dbURL. If there is an error opening the database, it will log a fatal error and exit.
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)                                 // database.New is a function that initializes a new instance of the Queries struct, which is defined in the internal/database package. This struct contains methods for interacting with the database, such as creating users, creating chirps, and querying for chirps. The dbConn variable is passed to this function to establish a connection to the database.


	platform := os.Getenv("PLATFORM")                                   // Get the platform from the environment variable PLATFORM. This can be used to differentiate between different deployment environments (e.g., development, staging, production). If it is not set, it will default to an empty string.





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
	mux.HandleFunc("POST /api/refresh", apicfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apicfg.handlerRevoke)
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)             // log.Printf is used to print the message to the console with a timestamp and other information, it is used here to indicate that the server is running and serving files from the specified directory on the specified port.
	log.Fatal(srv.ListenAndServe())                                                    // ListenAndServe listens on the TCP network address srv.Addr and then calls Serve to handle requests on incoming connections. It returns an error if any.
}

