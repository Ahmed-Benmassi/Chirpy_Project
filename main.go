package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {        //storing config setting fo interaction with api
	fileserverHits atomic.Int32     //atomic.Int32   provide a thread safe operation on 32bit int so there is no changes that can occure like go routines and so on
	
}



func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()                  //// NewServeMux allocates and returns a new [ServeMux]
	apicfg:=&apiConfig{}                          // create an instance of apiConfig to store the hits and use it in the handlers
	fs:=http.FileServer(http.Dir(filepathRoot))    // FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
	

	mux.Handle("/app/", http.StripPrefix("/app",apicfg.middlewareFileserverHits(fs)))  // handel the resuest and uses the stripprefixto rmove "/app" and add it to middlewar to wrap the file server and add the hit for every request to fs
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics",apicfg.getHits)
	mux.HandleFunc("POST /admin/reset",apicfg.resetHits)
	mux.HandleFunc("POST /api/validate_chirp",apicfg.handelervalidatechirp)
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

