package main

import (
	"context"
	"log"
	"net/http"
	"os"

	events "github.com/papawattu/cleanlog-eventstore/events"
	"github.com/papawattu/cleanlog-tasks/internal/repo"
)

const (
	defaultPort      = "3002"
	defaultAmqpURI   = "amqp://guest:guest@localhost:5672/"
	defaultQueueName = "taskqueue"
)

type config struct {
	port      string
	amqpURI   string
	queueName string
}

func getConfig() config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = defaultAmqpURI
	}

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = defaultQueueName
	}

	return config{port, amqpURI, queueName}
}
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

	cfg := getConfig()
	ctx := context.Background()

	taskRepo := repo.NewSimpleTaskRepository()

	var taskService TaskService

	if os.Getenv("EVENT_STORE") == "" || os.Getenv("EVENT_STREAM") == "" {
		taskService = NewTaskService(taskRepo)
	} else {
		eventBroadcaster := events.NewEventBroadcaster(ctx, taskRepo, os.Getenv("EVENT_STORE"), os.Getenv("EVENT_STREAM"), "task", "Task")
		taskService = NewTaskService(eventBroadcaster)
	}

	if err := startWebServer(":"+cfg.port, taskService); err != nil {
		log.Fatal(err)
	}

}
