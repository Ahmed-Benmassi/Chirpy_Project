package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
)


func (cfg *apiConfig) handlerUpdate(w http.ResponseWriter, r *http.Request){
	type req struct{
		Password string `json:"password"`
		Email  string `json:"email"`
	}
	type res struct{
		User
	}

	token,err :=auth.GetBearerToken(r.Header)
	if err !=nil{
		http.Error(w,"Couldn't find JWT",http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Couldn't validate JWT", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := req{}
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters",http.StatusInternalServerError)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w,"Couldn't hash password", http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		http.Error(w,"Couldn't update user", http.StatusInternalServerError)
		return
	}

	response := res{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}