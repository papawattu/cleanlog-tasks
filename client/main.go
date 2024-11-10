package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
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
		baseUri = "http://localhost:3000"
	}

	id := CreateTask("Task 1", baseUri)

	GetTask(id, baseUri)
}
