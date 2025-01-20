package handlers

import (
	"WebTasks/internal/models"
	"WebTasks/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUserRoutes(router *mux.Router, service services.UserService) {
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		GetUsers(w, r, service)
	}).Methods("GET")

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		CreateUser(w, r, service)
	}).Methods("POST")
}

func GetUsers(w http.ResponseWriter, r *http.Request, service services.UserService) {
	users, err := service.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request, service services.UserService) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdUser, err := service.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}
