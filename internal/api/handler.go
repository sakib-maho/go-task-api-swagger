package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/sakib-maho/go-task-api-swagger/internal/model"
	"github.com/sakib-maho/go-task-api-swagger/internal/store"
)

type Handler struct {
	store store.TaskStore
}

func NewHandler(taskStore store.TaskStore) *Handler {
	return &Handler{
		store: taskStore,
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListTasks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": h.store.List(),
	})
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Status = strings.TrimSpace(req.Status)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Status != "" && !isAllowedStatus(req.Status) {
		writeError(w, http.StatusBadRequest, "status must be one of: todo, in_progress, done")
		return
	}

	task, err := h.store.Create(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"data": task})
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	task, err := h.store.GetByID(taskID)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch task")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": task})
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")

	var req model.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Status = strings.TrimSpace(req.Status)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	if !isAllowedStatus(req.Status) {
		writeError(w, http.StatusBadRequest, "status must be one of: todo, in_progress, done")
		return
	}

	task, err := h.store.Update(taskID, req)
	if err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"data": task})
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if err := h.store.Delete(taskID); err != nil {
		if errors.Is(err, store.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isAllowedStatus(status string) bool {
	switch status {
	case "todo", "in_progress", "done":
		return true
	default:
		return false
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"message": message,
		},
	})
}
