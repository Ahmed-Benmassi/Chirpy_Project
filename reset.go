package main


import (
	"net/http"
	
)




func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {      //reset the hits to 0
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8\r\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to  0\r\n"))
}