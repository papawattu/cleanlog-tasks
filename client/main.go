package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/papawattu/cleanlog-tasks/types"
)

func CreateTask(description string, baseUri string) int {
	url := fmt.Sprintf("%s/api/task", baseUri)
	body := types.CreateTaskRequest{Description: description}
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return -1
	}

	slog.Info("Creating task with description", "description", description, "url", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error: status code", resp.StatusCode)
		return 0
	}

	r := map[string]int{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0
	}

	fmt.Println("Task created with ID:", r["taskId"])

	return r["taskId"]
}
func GetTask(id int, baseUri string) {
	url := fmt.Sprintf("%s/api/task/%d", baseUri, id)

	var resp *http.Response
	var err error
	var count int = 0

	for {
		resp, err = http.Get(url)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		if resp.StatusCode == http.StatusNotFound {
			if count > 20 {
				log.Fatalln("Task not found")
				return
			}
			count++
			log.Printf("Task not found at %s waiting - times %d\n", url, count)
			time.Sleep(1 * time.Second)
		} else {
			if resp.StatusCode == http.StatusOK {
				log.Printf("Task found at %s\n", url)
				break
			} else {
				log.Fatalf("Error: status code %d\n", resp.StatusCode)
				return
			}
		}
	}
	if resp != nil {

		defer resp.Body.Close()
	}

	r := &types.CreateTaskResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Task:", r)

}

func main() {

	baseUri := os.Getenv("BASE_URI")

	if baseUri == "" {
		baseUri = "http://localhost:3002"
	}

	id := CreateTask("Task 1", baseUri)

	GetTask(id, baseUri)
}
