package models

type Task struct {
	TaskID          *int
	TaskDescription string
	TaskType        string
}

func NewTask(description string) (Task, error) {

	t := Task{
		TaskID:          nil,
		TaskDescription: description,
	}
	return t, nil
}
