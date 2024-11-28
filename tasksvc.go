package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	repo "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-tasks/internal/models"
)

type TaskService interface {
	CreateTask(ctx context.Context, description string) (int, error)

	GetTask(ctx context.Context, id int) (*models.Task, error)

	DeleteTask(ctx context.Context, id int) error
}

type TaskServiceImp struct {
	repo repo.Repository[*models.Task, string]
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

	found := false

	for !found {
		t, _ := wsi.repo.Get(ctx, strconv.Itoa(nextId))
		if t != nil {
			found = true
		}
	}

	return nextId, nil
}

func (wsi *TaskServiceImp) GetTask(ctx context.Context, id int) (*models.Task, error) {

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}

func (wsi *TaskServiceImp) DeleteTask(ctx context.Context, id int) error {

	t, err := wsi.GetTask(ctx, id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}
	if t == nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil
	}
	err = wsi.repo.Delete(ctx, t)
	if err != nil {
		log.Fatalf("Error deleting work log: %v", err)
		return err
	}

	return nil
}
func NewTaskService(repo repo.Repository[*models.Task, string]) TaskService {

	return &TaskServiceImp{
		repo: repo,
	}
}
