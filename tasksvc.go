package main

import (
	"log"
	"math/rand"
)

type TaskService interface {
	CreateTask(description string) (int, error)

	GetTask(id int) (*Task, error)
}

type TaskServiceImp struct {
	repo TaskRepository
}

func nextId() int {
	return rand.Intn(1000)
}

func (wsi *TaskServiceImp) CreateTask(description string) (int, error) {

	wl, err := NewTask(description)
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
		return 0, err
	}

	nextId := nextId()
	wl.TaskID = &nextId

	err = wsi.repo.SaveTask(&wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return 0, err
	}
	return nextId, nil
}

func (wsi *TaskServiceImp) GetTask(id int) (*Task, error) {

	wl, err := wsi.repo.GetTask(id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}
func NewTaskService(repo TaskRepository) TaskService {

	return &TaskServiceImp{
		repo: repo,
	}
}
