package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {         // handlerRefresh is an HTTP handler function that processes incoming requests to refresh an access token using a provided refresh token. It performs several steps: defining the expected response payload, extracting the refresh token from the request header, validating the refresh token against the database, generating a new access token if valid, and responding with a JSON representation of the new access token along with an HTTP 200 status code. It handles errors by responding with appropriate HTTP status codes and messages for missing tokens, invalid tokens, and database errors.
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)                                  // extract the refresh token from the Authorization header using the GetBearerToken function
	if err != nil {
		http.Error(w,"Couldn't find token", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Couldn't get user for refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		http.Error(w,"Couldn't validate token", http.StatusUnauthorized)
		return
	}


	
	resp := response{
		Token: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {              // handlerRevoke is an HTTP handler function that processes incoming requests to revoke a refresh token, effectively logging the user out. It performs several steps: extracting the refresh token from the request header, revoking the token in the database, and responding with an HTTP 204 No Content status code if successful. It handles errors by responding with appropriate HTTP status codes and messages for missing tokens and database errors.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w,"Couldn't find token", http.StatusBadRequest)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w,"Couldn't revoke session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
