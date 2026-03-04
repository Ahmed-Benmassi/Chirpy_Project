package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {                // createChirpHandler is an HTTP handler function that processes incoming requests to create a new chirp. It performs several steps: defining the expected request payload, decoding the JSON body, validating the chirp content, inserting the chirp into the database using SQLC, mapping the database chirp to an API response struct, and finally responding with a JSON representation of the created chirp along with an HTTP 201 status code.
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
	chirp := Chirp{                                                                        // map the database chirp to the API response struct
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


func(cfg *apiConfig) getsinglechirphandeler(w http.ResponseWriter, r *http.Request){                         // getsinglechirphandeler is an HTTP handler function that retrieves a single chirp by its ID from the database, maps it to the API response struct, and responds with a JSON representation of the chirp along with an HTTP 200 status code. It handles errors by responding with appropriate HTTP status codes and messages for invalid chirp IDs, not found chirps, and unexpected database errors.
	 
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

	chirp:=Chirp{                                                                                 // map the database chirp to the API response struct
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


func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	ctx:=r.Context()
	if err != nil {
		http.Error(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	var dbChirps []database.Chirp

	if authorID != uuid.Nil {
		dbChirps, err = cfg.db.GetChirpsByAuthor(ctx, authorID)
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		http.Error(w, "Couldn't retrieve chirps", http.StatusInternalServerError)
		return
	}
 
	sortDirection := "asc"                                                        
	sortDirectionParam := r.URL.Query().Get("sort")                            
	if sortDirectionParam == "desc" {
		sortDirection = "desc"
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}


	sort.Slice(chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirps)
}
