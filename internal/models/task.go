package models

import (
	"strconv"

	common "github.com/papawattu/cleanlog-common"
)

type Task struct {
	common.BaseEntity[int]
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

func (wl *Task) GetID() string {
	return strconv.Itoa(*wl.TaskID)
}
