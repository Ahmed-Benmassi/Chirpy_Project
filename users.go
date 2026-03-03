package main

import (
	"encoding/json"
	"log"
	"net/http"

	
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){   // handel the request to create a new user and return the created user in json format
	var req struct{
		Email string `json:"email"`
	}

	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil{
		http.Error(w,"invalid request body",http.StatusBadRequest)
	    return
	}

	ctx:=r.Context()

	dbUser,err:=cfg.db.CreateUser(ctx,req.Email)
	if err!=nil{
		log.Printf("failed to create user with email %s: %v\n", req.Email, err)
		http.Error(w,"failed to create user",http.StatusInternalServerError)
		return
	}
	

	// Map database.User to main package User struct
    user := User{
        ID:        dbUser.ID,
        CreatedAt: dbUser.CreatedAt,
        UpdatedAt: dbUser.UpdatedAt,
        Email:     dbUser.Email,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)



}