package main

import (
	"encoding/json"
	"net/http"
    "time"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2️⃣ Validate inputs
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// 3️⃣ Look up user by email
	dbUser, err := cfg.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// User not found → 401 Unauthorized
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	// 4️⃣ Compare password with stored hash
	match, err := auth.CheckPasswordHash(req.Password, dbUser.HashedPassword)
	if err != nil || !match {
		// Password mismatch → 401 Unauthorized
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}



	expirationTime := time.Hour
	if req.ExpiresInSeconds > 0 && req.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(req.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		http.Error(w, "Couldn't create access JWT",http.StatusInternalServerError)
		return
	}



	// 5️⃣ Map DB user to response struct (omit password)
	resp := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Email:     dbUser.Email,
		Token:     accessToken,
	}

	// 6️⃣ Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}