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

<pre> git clone https://github.com/Ahmed-Benmassi/Chirpy_Project.git   
    cd Chirpy_Project </pre>

## 2. Install Go (1.20+)

<pre> https://go.dev/dl/</pre>

### Verify:

<pre>go version </pre>

## 3. Set Up PostgreSQL

### Create a database:
<pre> createdb chirpy_db </pre>

## 4. Environment Variables

### Create a .env file:

<pre> DB_URL=postgres://postgres:secret@localhost:5432/chirpy_db?sslmode=disable

    JWT_SECRET=supersecretkey

    PLATFORM=dev

    POLKA_KEY=webhook_secret </pre>

## 5. Install Dependencies

<pre> go mod download </pre>

## 6. Run the Server

<pre> go build -o out && ./out </pre>

### Server runs at:

<pre> http://localhost:8080</pre>

# 📘 API Documentation

## Base URL:

<pre> http://localhost:8080/api/healthz</pre> 


----------------------------------------


## 👤 Users
### Create User

<pre> curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
</pre>

  --------------------------------------------
  ## 🔐Authentication

### Login

<pre> curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
</pre>
  ## Response:

  <pre>  {
  "token": "JWT_TOKEN"
} </pre>

### Use this token for protected routes:

<pre> Authorization: Bearer JWT_TOKEN </pre>

---------------------------------------------------
## 🐦 Chirps

### Create Chirp (Authenticated)

<pre> curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "body": "Hello Chirpy!"
  }' </pre>

## Get All Chirps

<pre> curl http://localhost:8080/api/chirps </pre>

## Filter by Author

<pre> curl "http://localhost:8080/api/chirps?author_id=USER_UUID"</pre>


## Sort Chirps

### Ascending (default):

<pre> curl "http://localhost:8080/api/chirps?sort=asc"</pre>

### Descending:

<pre> curl "http://localhost:8080/api/chirps?sort=desc" </pre>

## Get Single Chirp

<pre> curl http://localhost:8080/api/chirps/CHIRP_ID</pre>

## Delete Chirp (Owner Only)

<pre> curl -X DELETE http://localhost:8080/api/chirps/CHIRP_ID \
  -H "Authorization: Bearer JWT_TOKEN"</pre>

---------------------------------------------------------------
## 🛡 Admin Endpoints (Development Only)

### Requires:

<pre> PLATFORM=dev</pre>

###  Reset Server Metrics:

<pre> curl -X POST http://localhost:8080/admin/reset</pre>

-----------------------------------------------------------------------
## 🧠 Architecture Overview

### -Standard library HTTP server

### -Layered design (handlers → service → database)

### -SQLC for type-safe queries

### -JWT-based stateless authentication

### -Context-based request handling

### -Explicit error handling

-----------------------------------------------------------------------
## 📌 Project Goals

### -This project demonstrates:

### -Understanding of HTTP request lifecycle

### -Manual routing without frameworks

### -Secure authentication implementation

### -Database query design

### -Clean API design

### -Production-style error handling

------------------------------------------------------------
# 💡 Future Ideas

## Integrate Chirpy with http_from_tcp Project:

## Create a Go client that interacts with the Chirpy API to:
   ####    -Make HTTP requests (GET, POST, DELETE, etc.)
   ####   -Integrate with your HTTP_from_TCP project for raw TCP requests
   ####     -Handle JSON responses automatically
   ####    -Authenticate using JWT tokens

   
## Dockerize the Project for Easy Deployment:
####       -Run the Chirpy server + PostgreSQL database with a single command.

## Testing & CI/CD:
####   -Unit tests for handlers and services.
####    -GitHub Actions to build and test automatically.
   
   
   
