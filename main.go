package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func startWebServer(port string, ts TaskService) error {

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

	log.Printf("Starting Task server on port %s\n", port)

	return http.ListenAndServe(port, nil)

}
func main() {

	port := ":3000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	taskRepo := NewTaskRepository()
	taskService := NewTaskService(taskRepo)

	if err := startWebServer(port, taskService); err != nil {
		log.Fatal(err)
	}

}
