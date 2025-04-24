package models

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type PostTaskDTO struct {
	Name string `json:"name"`
}

type Task struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	FilePath string    `json:"file_path"`
}

func NewTask(name, status, filePath string) *Task {
	return &Task{
		Id:       uuid.New(),
		Name:     name,
		Status:   status,
		FilePath: filePath,
	}
}

type TaskMap struct {
	sync.RWMutex
	tasks map[uuid.UUID]*Task
}

func NewTaskMap() *TaskMap {
	return &TaskMap{
		tasks: make(map[uuid.UUID]*Task),
	}
}

func (tm *TaskMap) PostTask(task *Task) {
	tm.Lock()
	defer tm.Unlock()
	tm.tasks[task.Id] = task
}

func (tm *TaskMap) GetTaskById(id uuid.UUID) (Task, error) {
	tm.RLock()
	defer tm.RUnlock()
	task, ok := tm.tasks[id]
	if !ok {
		errorMsg := fmt.Sprintf("Task with id %s not found", id.String())
		return Task{}, errors.New(errorMsg)
	}

	return *task, nil
}

func (tm *TaskMap) UpdateTaskStatus(id uuid.UUID, filePath string) error {
	tm.Lock()
	defer tm.Unlock()
	task, ok := tm.tasks[id]
	if !ok {
		errorMsg := fmt.Sprintf("Task with id %s not found", id.String())
		return errors.New(errorMsg)
	}

	task.Status = "completed"
	task.FilePath = filePath
	return nil
}
