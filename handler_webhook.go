package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
    "github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {      // handlerWebhook handles incoming webhook requests from Polka, validating the API key, decoding the event data, and upgrading the user's account to Chirpy Red if the event indicates a user upgrade. 
	                                                                                //It checks for the presence of a valid API key in the request header, decodes the JSON payload to extract the event type and user ID, and if the event is "user.upgraded"
                                                                                    // it calls the database method to upgrade the user's account. The handler responds with appropriate HTTP status codes based on the success or failure of these operations.
    type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, "Couldn't find api key", http.StatusUnauthorized)
		return
	}
	if apiKey != cfg.polkaKey {
		http.Error(w,  "API key is invalid", http.StatusUnauthorized,)
		return
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w,"Couldn't find user", http.StatusNotFound)
			return
		}
		http.Error(w, "Couldn't update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}