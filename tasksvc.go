package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/papawattu/cleanlog-tasks/internal/models"
	"github.com/papawattu/cleanlog-tasks/internal/repo"
)

type TaskService interface {
	CreateTask(ctx context.Context, description string) (int, error)

	GetTask(ctx context.Context, id int) (*models.Task, error)
}

type TaskServiceImp struct {
	repo repo.Repository[*models.Task, int]
}

func nextId() int {
	return rand.Intn(1000)
}

func (wsi *TaskServiceImp) CreateTask(ctx context.Context, description string) (int, error) {

	wl, err := models.NewTask(description)
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
		return 0, err
	}

	nextId := nextId()
	wl.TaskID = &nextId

	err = wsi.repo.Save(ctx, &wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return 0, err
	}
	return nextId, nil
}

func (wsi *TaskServiceImp) GetTask(ctx context.Context, id int) (*models.Task, error) {

	wl, err := wsi.repo.Get(ctx, id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}
func NewTaskService(repo repo.Repository[*models.Task, int]) TaskService {

	return &TaskServiceImp{
		repo: repo,
	}
}
