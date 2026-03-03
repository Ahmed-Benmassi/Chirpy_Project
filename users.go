package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ahmed-Benmassi/chirpy_Project/internal/auth"
	"github.com/Ahmed-Benmassi/chirpy_Project/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){   // handel the request to create a new user and return the created user in json format
	var req struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil{                        // decode the request body into the req struct
		http.Error(w,"invalid request body",http.StatusBadRequest)                
	    return
	}

	// 3️⃣ Validate inputs
	if req.Email == "" || req.Password == "" {                                    // Check if email and password are provided
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	ctx:=r.Context()

	hashedPassword, err := auth.HashPassword(req.Password)                          // Hash the password using the auth package
	if err != nil {
		log.Printf("failed to hash password: %v\n", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}


	dbUser,err:=cfg.db.CreateUser(ctx,database.CreateUserParams{                     // Create a new user in the database using the database package
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	
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