package main

import (
	"context"
	"log"
	"net/http"
)

func startWebServer(ts TaskService) error {

	controllers, err := MakeControllers(context.Background(), NewTaskController(ts))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err = controllers.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	http.HandleFunc("/api/task/{taskid}", controllers.HandleRequest)
	http.HandleFunc("/api/task", controllers.HandleRequest)

	http.ListenAndServe(":3000", nil)

	return nil
}
func main() {

	taskRepo := NewTaskRepository()
	taskService := NewTaskService(taskRepo)

	startWebServer(taskService)

}
