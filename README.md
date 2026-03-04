# 🐦 Chirpy — REST API in Go

### Chirpy is a production-style RESTful JSON API built from scratch using Go’s standard library.
### The project demonstrates core backend engineering principles including routing, middleware, authentication, authorization, PostgreSQL integration, and webhook handling.

### This project intentionally avoids heavy frameworks to showcase a deep understanding of HTTP servers and backend fundamentals.

## ✨ Features

### -RESTful JSON API

### -PostgreSQL persistence

### -JWT authentication

### -Role-based authorization

### -Query filtering and sorting

### -Middleware (logging & metrics)

### -Admin-only development endpoints

### -Webhook support

### -Clean layered architecture


## 🛠 Tech Stack
### -Go (net/http)

### -PostgreSQL

### -SQLC

### -JWT (HMAC)

------------------------------------------
# 🚀 Getting Started
##  1. Clone the repository 

git clone https://github.com/Ahmed-Benmassi/Chirpy_Project.git
cd Chirpy_Project

## 2. Install Go (1.20+)

https://go.dev/dl/

### Verify:

go version

## 3. Set Up PostgreSQL

### Create a database:
createdb chirpy_db

## 4. Environment Variables

### Create a .env file:

DB_URL=postgres://postgres:secret@localhost:5432/chirpy_db?sslmode=disable
JWT_SECRET=supersecretkey
PLATFORM=dev
POLKA_KEY=webhook_secret

## 5. Install Dependencies

go mod download

## 6. Run the Server

go build -o out && ./out

### Server runs at:

http://localhost:8080
--------------------------------------------
# 📘 API Documentation

## Base URL:
http://localhost:8080/api/healthz


----------------------------------------


#👤 Users
## Create User

curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'


  --------------------------------------------
  # 🔐 Authentication

##Login

curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

  ## Response:

  {
  "token": "JWT_TOKEN"
}

### Use this token for protected routes:

Authorization: Bearer JWT_TOKEN

---------------------------------------------------
# 🐦 Chirps

## Create Chirp (Authenticated)

curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "body": "Hello Chirpy!"
  }'

## Get All Chirps

curl http://localhost:8080/api/chirps

## Filter by Author

curl "http://localhost:8080/api/chirps?author_id=USER_UUID"


## Sort Chirps

### Ascending (default):

curl "http://localhost:8080/api/chirps?sort=asc"

### Descending:

curl "http://localhost:8080/api/chirps?sort=desc"

## Get Single Chirp

curl http://localhost:8080/api/chirps/CHIRP_ID

## Delete Chirp (Owner Only)

curl -X DELETE http://localhost:8080/api/chirps/CHIRP_ID \
  -H "Authorization: Bearer JWT_TOKEN"

---------------------------------------------------------------
# 🛡 Admin Endpoints (Development Only)

## Requires:

PLATFORM=dev

##  Reset Server Metrics:

curl -X POST http://localhost:8080/admin/reset

-----------------------------------------------------------------------
# 🧠 Architecture Overview

### -Standard library HTTP server

### -Layered design (handlers → service → database)

### -SQLC for type-safe queries

### -JWT-based stateless authentication

### -Context-based request handling

### -Explicit error handling

-----------------------------------------------------------------------
# 📌 Project Goals

## -This project demonstrates:

## -Understanding of HTTP request lifecycle

## -Manual routing without frameworks

## -Secure authentication implementation

## -Database query design

## -Clean API design

## -Production-style error handling

------------------------------------------------------------
#💡 Future Ideas

## Integrate Chirpy with http_from_tcp Project

## Create a Go client that interacts with the Chirpy API to:
   ###    -Make HTTP requests (GET, POST, DELETE, etc.)
   ###    -Integrate with your HTTP_from_TCP project for raw TCP requests
   ###     -Handle JSON responses automatically
   ###    -Authenticate using JWT tokens
   
## Dockerize the Project for Easy Deployment
###       -Run the Chirpy server + PostgreSQL database with a single command.

##Testing & CI/CD
    ###   -Unit tests for handlers and services.

   ###    -GitHub Actions to build and test automatically.
   
   
   
