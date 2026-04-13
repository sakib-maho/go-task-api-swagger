package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/sakib-maho/go-task-api-swagger/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStore interface {
	List() []model.Task
	Create(req model.CreateTaskRequest) (model.Task, error)
	GetByID(id string) (model.Task, error)
	Update(id string, req model.UpdateTaskRequest) (model.Task, error)
	Delete(id string) error
}

type MemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[string]model.Task
}

func NewMemoryTaskStore() *MemoryTaskStore {
	return &MemoryTaskStore{
		tasks: make(map[string]model.Task),
	}
}

func (s *MemoryTaskStore) List() []model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *MemoryTaskStore) Create(req model.CreateTaskRequest) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	status := req.Status
	if status == "" {
		status = "todo"
	}

	task := model.Task{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	return task, nil
}

func (s *MemoryTaskStore) GetByID(id string) (model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (s *MemoryTaskStore) Update(id string, req model.UpdateTaskRequest) (model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return model.Task{}, ErrTaskNotFound
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.UpdatedAt = time.Now().UTC()
	s.tasks[id] = task
	return task, nil
}

func (s *MemoryTaskStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
