package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {        //storing config setting fo interaction with api
	fileserverHits atomic.Int32     //atomic.Int32   provide a thread safe operation on 32bit int so there is no changes that can occure like go routines and so on
	
}


func (cfg *apiConfig) middlewareFileserverHits(next http.Handler) http.Handler {     //methods that add other functionality to fileserver that add a hit fo evry request to fs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)                                               //calls the next handeler
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request)  {      //get th enubmer of reauests as a hits and return it a sa response
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits: "+ fmt.Sprintf("%d",cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {      //reset the hits to 0
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to  0"))
}



func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()                  //// NewServeMux allocates and returns a new [ServeMux]
	apicfg:=&apiConfig{}                          // create an instance of apiConfig to store the hits and use it in the handlers
	fs:=http.FileServer(http.Dir(filepathRoot))    // FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
	

	mux.Handle("/app/", http.StripPrefix("/app",apicfg.middlewareFileserverHits(fs)))  // handel the resuest and uses the stripprefixto rmove "/app" and add it to middlewar to wrap the file server and add the hit for every request to fs
	mux.HandleFunc("/", handlerReadiness)
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics",apicfg.getHits)
	mux.HandleFunc("/reset",apicfg.resetHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

