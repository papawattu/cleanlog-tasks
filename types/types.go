package types

type CreateTaskRequest struct {
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	TaskId      int    `json:"taskId"`
	Description string `json:"description"`
}
