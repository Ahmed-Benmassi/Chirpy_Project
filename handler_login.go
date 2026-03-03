package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {     // loginHandler handles user login by validating the provided email and password, generating JWT access and refresh tokens upon successful authentication, and returning the user information along with the tokens in the response.
	// 1️⃣ Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		
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

	ctx := r.Context()                                                               // Create a context from the request, which can be used for database operations and other context-aware functions. This allows for better control over request lifecycle and cancellation.

	// 3️⃣ Look up user by email
	dbUser, err := cfg.db.GetUserByEmail(ctx, req.Email)                             // Use the database query method GetUserByEmail to retrieve the user record associated with the provided email address. This will return a user object if found, or an error if the user does not exist or if there was an issue with the database query.
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



	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		http.Error(w, "Couldn't create access JWT", http.StatusInternalServerError)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    dbUser.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})



	if err != nil {
		http.Error(w, "Couldn't save refresh token",  http.StatusInternalServerError)
		return
	}




	// 5️⃣ Map DB user to response struct (omit password)
	resp := struct {                                                                           // Define an anonymous struct type for the response, which includes the user ID, creation and update timestamps, email, access token, and refresh token. This struct will be used to format the JSON response sent back to the client.
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
		RefreshToken: refreshToken,
	}

	// 6️⃣ Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}