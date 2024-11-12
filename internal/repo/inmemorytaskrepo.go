package repo

import (
	"context"

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
	return nil
}

func NewSimpleTaskRepository() Repository[*models.Task, int] {
	return &TaskSimpleRepositoryImp{
		Tasks: make(map[int]*models.Task),
	}
}
