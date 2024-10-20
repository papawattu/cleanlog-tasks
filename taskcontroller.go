package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/papawattu/cleanlog-tasks/types"
)

type TaskController struct {
	taskService TaskService
	handlePost  func(w http.ResponseWriter, r *http.Request)
	handleGet   func(w http.ResponseWriter, r *http.Request)
}

func NewTaskController(taskService TaskService) *TaskController {

	return &TaskController{taskService: taskService, handlePost: func(w http.ResponseWriter, r *http.Request) {
		var t types.CreateTaskRequest

		json.NewDecoder(r.Body).Decode(&t)

		workID, err := taskService.CreateTask(t.Description)
		if err != nil {
			log.Fatalf("Error starting task controller: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"taskId": workID})
	}, handleGet: func(w http.ResponseWriter, r *http.Request) {
		taskId := r.PathValue("taskid")
		if taskId == "" {
			http.Error(w, "taskId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(taskId)
		if err != nil {
			http.Error(w, "taskId must be an integer %s", http.StatusBadRequest)
			return
		}

		t, err := taskService.GetTask(id)
		if err != nil {
			log.Fatalf("Error getting work log: %v", err)
			http.Error(w, "Error getting work log", http.StatusInternalServerError)
		}

		if t == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {

			r := types.CreateTaskResponse{TaskId: *t.TaskID, Description: t.TaskDescription}
			json.NewEncoder(w).Encode(r)
		}

	}}
}

func (wc *TaskController) Start() error {
	log.Printf("Starting task controller")
	return nil
}

func (wc *TaskController) Stop() error {
	log.Printf("Stopping task controller")
	return nil
}

func (wc *TaskController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request: %v", r)

	switch {
	case r.Method == "POST":
		wc.handlePost(w, r)
	case r.Method == "GET":
		wc.handleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
