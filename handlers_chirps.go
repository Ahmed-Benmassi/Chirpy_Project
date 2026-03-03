package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Define request payload
	type ChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}


	var req ChirpRequest

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, "Couldn't find JWT", http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		http.Error(w, "Couldn't validate JWT", http.StatusUnauthorized)
		return
	}




	// 2️⃣ Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3️⃣ Validate chirp
	if len(req.Body) == 0 {
		http.Error(w, "Chirp body cannot be empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 4️⃣ Insert chirp using SQLC CreateChirp
	dbChirp, err := cfg.db.CreateChirp(ctx, database.CreateChirpParams{
		Body:   req.Body,
		UserID: userID,
	})
	if err != nil {
		http.Error(w, "Failed to create chirp", http.StatusInternalServerError)
		return
	}

	// 5️⃣ Map to API Chirp struct
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// 6️⃣ Respond with HTTP 201 and JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirp)
}


func (cfg *apiConfig) listChirpsHandler(w http.ResponseWriter, r *http.Request) {


	ctx := r.Context()
	dbChirp,err:=cfg.db.ListALLChirps(ctx)
	if err!=nil{
		http.Error(w, "Failed to list chirps", http.StatusInternalServerError)
		return 
	}

	chirps:=[]Chirp{}
	for _,c:=range dbChirp{
		chirps=append(chirps,Chirp{
			ID:        c.ID,
            CreatedAt: c.CreatedAt,
            UpdatedAt: c.UpdatedAt,
            Body:      c.Body,
            UserID:    c.UserID,
		})
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}

func(cfg *apiConfig) getsinglechirphandeler(w http.ResponseWriter, r *http.Request){
	
	ctx:=r.Context()
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		http.Error(w, "Invalid chirp ID",http.StatusBadRequest)
		return
	}


	dbChirp,err:=cfg.db.GetChirp(ctx,chirpID)
	if err!=nil{
		if err == sql.ErrNoRows {
            // chirp does not exist → 404
            http.Error(w, "Chirp not found", http.StatusNotFound)
            return
        }
        // unexpected DB error → 500
        http.Error(w, "Could not get the chirp", http.StatusInternalServerError)
        return
	}

	chirp:=Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirp)



}