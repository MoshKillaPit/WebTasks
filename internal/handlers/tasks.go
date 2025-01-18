package handlers

import (
	"WebTasks/internal/models"
	"WebTasks/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func RegisterTaskRoutes(router *mux.Router, service *services.TaskServiceImpl) {
	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		GetTasks(w, r, service)
	}).Methods("GET")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetTaskByID(w, r, service)
	}).Methods("GET")

	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		CreateTask(w, r, service)
	}).Methods("POST")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		UpdateTask(w, r, service)
	}).Methods("PUT")

	router.HandleFunc("/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteTask(w, r, service)
	}).Methods("DELETE")
}

func GetTasks(w http.ResponseWriter, r *http.Request, service services.TaskService) {
	tasks, err := service.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request, service services.TaskService) {
	var task models.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	createdTask, err := service.Create(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(createdTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request, service services.TaskService) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var task models.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = id
	updatedTask, err := service.Update(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(updatedTask)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request, service services.TaskService) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := service.GetByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request, service services.TaskService) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = service.Delete(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
