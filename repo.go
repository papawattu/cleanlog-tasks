package main

type TaskRepository interface {
	SaveTask(t *Task) error
	GetTask(id int) (*Task, error)
}

type TaskRepositoryImp struct {
	Tasks map[int]*Task
}

func (wri *TaskRepositoryImp) SaveTask(t *Task) error {
	wri.Tasks[*t.TaskID] = t
	return nil
}

func (wri *TaskRepositoryImp) GetTask(id int) (*Task, error) {
	t, ok := wri.Tasks[id]
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
