package main

import (
	"log"
	"net/http"
)




func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {      //reset the hits to 0
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8\r\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to  0\r\n"))
}



func (cfg *apiConfig) adminResethandler(w http.ResponseWriter, r *http.Request) {   

	if cfg.platform != "dev" {
		http.Error(w,"forbidden",http.StatusForbidden)
		return
	}

	ctx:=r.Context()

	if err:=cfg.db.DeleteAllUsers(ctx);err!=nil{
		log.Printf("failed to delete  user users: %v\n", err)
		http.Error(w,"failed to delete users",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All users deleted successfully\n"))

}    