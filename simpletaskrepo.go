package main

type TaskSimpleRepositoryImp struct {
	Tasks map[int]*Task
}

func (tri *TaskSimpleRepositoryImp) SaveTask(t *Task) error {
	tri.Tasks[*t.TaskID] = t
	return nil
}

func (tri *TaskSimpleRepositoryImp) GetTask(id int) (*Task, error) {
	t, ok := tri.Tasks[id]
	if !ok {
		return nil, nil
	}
	return t, nil
}
func NewSimpleTaskRepository() TaskRepository {
	return &TaskSimpleRepositoryImp{
		Tasks: make(map[int]*Task),
	}
}
