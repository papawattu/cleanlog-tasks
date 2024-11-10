package main

type TaskRepository interface {
	SaveTask(t *Task) error
	GetTask(id int) (*Task, error)
}
