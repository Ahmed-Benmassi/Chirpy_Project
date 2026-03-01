package main

import (
	"fmt"
	"net/http"
	
)





func (cfg *apiConfig) middlewareFileserverHits(next http.Handler) http.Handler {     //methods that add other functionality to fileserver that add a hit fo evry request to fs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)                                               //calls the next handeler
	})
}



func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request)  {      //get th enubmer of reauests as a hits and return it a sa response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html>\n    <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n    </body>\n    </html>", cfg.fileserverHits.Load())))
}