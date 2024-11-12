package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/papawattu/cleanlog-tasks/types"
)

type TaskController struct {
	taskService  TaskService
	handlePost   func(w http.ResponseWriter, r *http.Request)
	handleGet    func(w http.ResponseWriter, r *http.Request)
	handleDelete func(w http.ResponseWriter, r *http.Request)
}

func NewTaskController(taskService TaskService) *TaskController {

	return &TaskController{taskService: taskService, handlePost: func(w http.ResponseWriter, r *http.Request) {
		var t types.CreateTaskRequest

		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		workID, err := taskService.CreateTask(r.Context(), t.Description)
		if err != nil {
			log.Fatalf("Error starting task controller: %v", err)
		}

		w.Header().Set("Location", "/api/task/"+strconv.Itoa(workID))
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(map[string]int{"taskId": workID})
		if err != nil {
			log.Fatalf("Error encoding response: %v", err)
		}
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

		t, err := taskService.GetTask(r.Context(), id)
		if err != nil {
			log.Fatalf("Error getting work log: %v", err)
			http.Error(w, "Error getting work log", http.StatusInternalServerError)
		}

		if t == nil {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {

			r := types.CreateTaskResponse{TaskId: *t.TaskID, Description: t.TaskDescription}
			err := json.NewEncoder(w).Encode(r)
			if err != nil {
				log.Fatalf("Error encoding response: %v", err)
			}
		}
	}, handleDelete: func(w http.ResponseWriter, r *http.Request) {
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

		err = taskService.DeleteTask(r.Context(), id)
		if err != nil {
			log.Fatalf("Error deleting task: %v", err)
			http.Error(w, "Error deleting task", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	},
	}
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
