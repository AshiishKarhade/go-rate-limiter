package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	client := InitRedisConnection()
	//fmt.Println(client)
	rateLimiter := &RateLimiter{client: client}

	r := mux.NewRouter()

	r.HandleFunc("/register", RegisterUser).Methods("POST")
	r.HandleFunc("/login", LoginUser).Methods("POST")

	r.HandleFunc("/api/v1/token", func(w http.ResponseWriter, r *http.Request) {

	})
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

}

func LoginUser(w http.ResponseWriter, r *http.Request) {

}
