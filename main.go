package main

import (
	"context"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	common "github.com/papawattu/cleanlog-common"
	events "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-tasks/internal/models"
)

type Config struct {
	Port        string `envconfig:"PORT" default:"3000"`
	EventStore  string `envconfig:"EVENT_STORE"`
	EventStream string `envconfig:"EVENT_STREAM"`
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

	var cfg Config

	err := envconfig.Process("task", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var taskService TaskService

	if cfg.EventStore == "" || cfg.EventStream == "" {
		taskRepo := common.NewInMemoryRepository[*models.Task]()
		taskService = NewTaskService(taskRepo)
	} else {
		ht := events.NewHttpTransport(cfg.EventStore, cfg.EventStream, 10)
		taskRepo := common.NewMemcacheRepository[*models.Task]("localhost:11211", "task", nil)
		eventBroadcaster := common.NewEventService(taskRepo, ht, "Task")

		taskService = NewTaskService(eventBroadcaster)

		eventBroadcaster.StartEventRunner(context.Background())
	}

	if err := startWebServer(":"+cfg.Port, taskService); err != nil {
		log.Fatal(err)
	}

}
