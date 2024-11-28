package main

import (
	"context"
	"log"
	"net/http"
	"os"

	common "github.com/papawattu/cleanlog-common"
	events "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-tasks/internal/models"
)

type config struct {
	port      string
	amqpURI   string
	queueName string
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

	var taskService TaskService

	if os.Getenv("EVENT_STORE") == "" || os.Getenv("EVENT_STREAM") == "" {
		taskRepo := common.NewInMemoryRepository[*models.Task]()
		taskService = NewTaskService(taskRepo)
	} else {
		ht := events.NewHttpTransport(os.Getenv("EVENT_STORE"), os.Getenv("EVENT_STREAM"), 10)
		taskRepo := common.NewMemcacheRepository[*models.Task]("localhost:11211", "task", nil)
		eventBroadcaster := common.NewEventService(taskRepo, ht, "Task")

		taskService = NewTaskService(eventBroadcaster)

		eventBroadcaster.StartEventRunner(context.Background())
	}

	if err := startWebServer(":3000", taskService); err != nil {
		log.Fatal(err)
	}

}
