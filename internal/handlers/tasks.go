package handlers

import (
	"WebTasks/internal/models"
	"WebTasks/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service services.TaskService
}

func NewHandler(service services.TaskService) *Handler {
	return &Handler{service: service}
}

func RegisterTaskRoutes(router *mux.Router, handler *Handler) {
	router.HandleFunc("/tasks", handler.GetTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", handler.GetTaskByID).Methods(http.MethodGet)
	router.HandleFunc("/tasks", handler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", handler.UpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", handler.DeleteTask).Methods(http.MethodDelete)
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	h.writeJSON(w, http.StatusOK, task)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdTask, err := h.service.Create(ctx, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, http.StatusCreated, createdTask)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.ID = id

	updatedTask, err := h.service.Update(ctx, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.writeJSON(w, http.StatusOK, updatedTask)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) parseID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	return strconv.Atoi(idStr)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}
}
