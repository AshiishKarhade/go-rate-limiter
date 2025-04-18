package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var client *redis.Client

func main() {
	client = InitRedisConnection()
	//fmt.Println(client)
	rateLimiter := &RateLimiter{client: client}

	r := mux.NewRouter()

	r.HandleFunc("/register", RegisterUser).Methods("POST")
	r.HandleFunc("/login", LoginUser).Methods("POST")

	r.HandleFunc("/api/v1/token", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User-ID")
		if userID == "" {
			http.Error(w, "User-ID header is required", http.StatusBadRequest)
			return
		}

		// Initialize rate limit for the user if it doesn't exist
		rateLimiter.InitializeRateLimit(userID)

		// Check if the request is allowed by the rate limiter
		if rateLimiter.AllowRequest(userID) {
			// Proxy to the backend API (dummy response for now)
			w.Write([]byte("Token generated successfully"))
		} else {
			// Reject request if rate limit exceeded
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			w.Header().Set("Retry-After", "60") // Suggest retry after 60 seconds
		}
	}).Methods("GET")
	fmt.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid user data", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword := hashPassword(user.Password)

	err = client.HSet(ctx, "user:"+user.Username, "username", user.Username, "password_hash", hashedPassword, "registration_date", time.Now().Format(time.RFC3339)).Err()
	if err != nil {
		http.Error(w, "Error registering User", http.StatusInternalServerError)
		return
	}
	client.HSet(ctx, "rate_limit:"+user.Username, "tokens", maxTokens, "last_refill_time", time.Now().Format(time.RFC3339))

	w.Write([]byte("User registered successfully"))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve the stored password hash from Redis
	storedPasswordHash, err := client.HGet(ctx, "user:"+user.Username, "password_hash").Result()
	if err == redis.Nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Error retrieving user data", http.StatusInternalServerError)
		return
	}

	// Verify the password against the stored hash
	if !verifyPassword(user.Password, storedPasswordHash) {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Respond with a success message if login is successful
	w.Write([]byte("User logged in successfully"))
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password")
	}
	return string(hashedPassword)
}

func verifyPassword(password, storedHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}
