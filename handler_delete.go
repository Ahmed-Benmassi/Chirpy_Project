package main

import (
	"net/http"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {  
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		http.Error(w,"Invalid chirp ID", http.StatusBadRequest)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Couldn't find JWT",http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		http.Error(w,"Couldn't validate JWT", http.StatusUnauthorized)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		http.Error(w,"Couldn't get chirp", http.StatusNotFound)
		return
	}
	if dbChirp.UserID != userID {
		http.Error(w,"You can't delete this chirp", http.StatusForbidden)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		http.Error(w, "Couldn't delete chirp", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}