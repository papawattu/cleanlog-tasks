package repo

import (
	"context"
	"log/slog"

	repo "github.com/papawattu/cleanlog-eventstore/repository"
	"github.com/papawattu/cleanlog-tasks/internal/models"
)

type TaskSimpleRepositoryImp struct {
	Tasks map[int]*models.Task
}

func (tri *TaskSimpleRepositoryImp) Save(ctx context.Context, t *models.Task) error {
	tri.Tasks[*t.TaskID] = t
	return nil
}

func (tri *TaskSimpleRepositoryImp) Get(ctx context.Context, id int) (*models.Task, error) {
	t, ok := tri.Tasks[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}
func (tri *TaskSimpleRepositoryImp) GetAll(ctx context.Context) ([]*models.Task, error) {
	tasks := []*models.Task{}
	for _, t := range tri.Tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (tri *TaskSimpleRepositoryImp) Delete(ctx context.Context, t *models.Task) error {
	delete(tri.Tasks, *t.TaskID)
	slog.Info("Task deleted", "Task ID", *t.TaskID)
	return nil
}

func (tri *TaskSimpleRepositoryImp) GetId(ctx context.Context, t *models.Task) (int, error) {
	return *t.TaskID, nil
}

func (tri *TaskSimpleRepositoryImp) Exists(ctx context.Context, id int) (bool, error) {
	_, ok := tri.Tasks[id]
	return ok, nil
}

func NewSimpleTaskRepository() repo.Repository[*models.Task, int] {
	return &TaskSimpleRepositoryImp{
		Tasks: make(map[int]*models.Task),
	}
}
