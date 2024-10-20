package main

type TaskRepository interface {
	SaveTask(t *Task) error
	GetTask(id int) (*Task, error)
}

type TaskRepositoryImp struct {
	Tasks map[int]*Task
}

func (tri *TaskRepositoryImp) SaveTask(t *Task) error {
	tri.Tasks[*t.TaskID] = t
	return nil
}

func (tri *TaskRepositoryImp) GetTask(id int) (*Task, error) {
	t, ok := tri.Tasks[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}
func NewTaskRepository() TaskRepository {
	return &TaskRepositoryImp{
		Tasks: make(map[int]*Task),
	}
}
